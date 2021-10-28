package cloud

/**
Note:

--- Spot instance charge ---
Spot instances are charged for actual seconds used,
even if they are released early during the protection period(default 1 hour on aliyun).
*/

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gdecimal"
	"github.com/cryptowilliam/goutil/container/gstring"
	"github.com/cryptowilliam/goutil/container/gvolume"
	"github.com/cryptowilliam/goutil/encoding/gjson"
	"github.com/cryptowilliam/goutil/i18n/gfiat"
	"github.com/cryptowilliam/goutil/sys/gio"
	"github.com/cryptowilliam/goutil/sys/gtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type (
	AliyunClient struct {
		accessKey string
		secretKey string
		lan       bool
		ecsList   map[string]*ecs.Client
		ecsListMu sync.RWMutex
		ossList   map[string]*oss.Client
		ossListMu sync.RWMutex
		buckets   map[string]*oss.Bucket
		bucketsMu sync.RWMutex
	}

	Endpoint struct {
		WAN string
		LAN string
	}

	RegionInfo struct {
		Country   string
		City      string
		Endpoints map[Production]Endpoint
	}
)

func newAliyun(accessKey string, secretKey string, LAN bool) (*AliyunClient, error) {
	return &AliyunClient{
		accessKey: accessKey,
		secretKey: secretKey,
		lan:       LAN,
		ecsList:   map[string]*ecs.Client{},
		ossList:   map[string]*oss.Client{},
		buckets:   map[string]*oss.Bucket{},
	}, nil
}

func (ac *AliyunClient) GetBalance() (*Balance, error) {
	cli, err := bssopenapi.NewClientWithAccessKey("cn-hongkong", ac.accessKey, ac.secretKey)
	if err != nil {
		return nil, err
	}
	request := bssopenapi.CreateQueryAccountBalanceRequest()
	resp, err := cli.QueryAccountBalance(request)
	if err != nil {
		return nil, err
	}
	asset, err := gfiat.ParseFiat(resp.Data.Currency)
	if err != nil {
		return nil, err
	}
	available, err := gdecimal.NewFromString(resp.Data.AvailableAmount)
	if err != nil {
		return nil, err
	}

	return &Balance{
		Currency:  asset,
		Available: available,
	}, err
}

func (ac *AliyunClient) ListRegions() ([]string, error) {
	cli, err := ac.getEcsClient("cn-hongkong")
	if err != nil {
		return nil, err
	}
	resp, err := cli.DescribeRegions(ecs.CreateDescribeRegionsRequest())
	if err != nil {
		return nil, err
	}
	var res []string
	for _, v := range resp.Regions.Region {
		res = append(res, v.RegionId)
	}
	return res, nil
}

func (ac *AliyunClient) ListZones(regionId string) ([]string, error) {
	cli, err := ac.getEcsClient(regionId)
	if err != nil {
		return nil, err
	}
	req := ecs.CreateDescribeZonesRequest()
	req.RegionId = regionId
	resp, err := cli.DescribeZones(req)
	if err != nil {
		return nil, err
	}
	var res []string
	for _, v := range resp.Zones.Zone {
		res = append(res, v.ZoneId)
	}
	return res, nil
}

func (ac *AliyunClient) EcListOnDemandSpecs(regionId string) (*VpsSpecList, error) {
	cli, err := ac.getEcsClient(regionId)
	if err != nil {
		return nil, err
	}
	resMap := map[string]VpsSpec{}

	// Get hardware.
	itReq := ecs.CreateDescribeInstanceTypesRequest()
	itReq.RegionId = regionId
	ditResp, err := cli.DescribeInstanceTypes(itReq)
	if err != nil {
		return nil, err
	}
	for _, instanceType := range ditResp.InstanceTypes.InstanceType {
		memoryVolume, err := gvolume.FromByteSize(instanceType.MemorySize * float64(gvolume.GB.Bytes()))
		if err != nil {
			return nil, err
		}
		entry := VpsSpec{
			RegionId: regionId,
			Id:            instanceType.InstanceTypeId,
			LogicalCpuNum: instanceType.CpuCoreCount,
			MemoryVolume:  memoryVolume,
		}
		entry.IsCredit = gstring.StartWith(entry.Id, "ecs.t")
		resMap[instanceType.InstanceTypeId] = entry
	}

	// Get available zones.
	arrReq := ecs.CreateDescribeAvailableResourceRequest()
	arrReq.DestinationResource = "InstanceType"
	arrReq.RegionId = regionId
	darResp, err := cli.DescribeAvailableResource(arrReq)
	if err != nil {
		return nil, err
	}
	for _, zoneInfo := range darResp.AvailableZones.AvailableZone { // Iterate all zones.
		for _, resInfo := range zoneInfo.AvailableResources.AvailableResource[0].SupportedResources.SupportedResource { // Iterate all instance specs at specified zone.
			if resInfo.Status != "Available" {
				continue
			}
			originVal := resMap[resInfo.Value] // resInfo.Value is instance spec like "ecs.t6-c2m1.large".
			originVal.AvailableZoneIds = append(originVal.AvailableZoneIds, zoneInfo.ZoneId)
			resMap[resInfo.Value] = originVal
		}
	}

	// Get on demand price.
	pReq := ecs.CreateDescribePriceRequest()
	pReq.RegionId = regionId
	pReq.ResourceType = "instance"
	for _, spec := range resMap {
		if len(spec.AvailableZoneIds) == 0 {
			continue
		}
		pReq.InstanceType = spec.Id // like "ecs.g6.large"
		priceErr := error(nil)
		// In document, if set "ResourceType" as "instance", user has no necessary to
		// set "SystemDiskCategory" in request, its default value is "cloud_efficiency".
		// But in fact, if user didn't set "SystemDiskCategory" in some instance specs,
		// api will return error, because default disk type is not available for that instance spec.
		// So, I use below string array to try to find instance price.
		for _, diskCategory := range []string{"cloud_efficiency", "cloud_ssd", "ephemeral_ssd", "cloud_essd"} {
			pReq.SystemDiskCategory = diskCategory
			time.Sleep(time.Millisecond * 10) // avoid request throttling
			dpr := ecs.CreateDescribePriceResponse()
			dpr, priceErr = cli.DescribePrice(pReq)
			if priceErr == nil {
				spec.Currency, err = gfiat.ParseFiat(dpr.PriceInfo.Price.Currency)
				if err != nil {
					return nil, err
				}
				for _, zoneId := range spec.AvailableZoneIds {
					spec.OnDemandPrices[zoneId] = gdecimal.NewFromFloat64(dpr.PriceInfo.Price.OriginalPrice)
				}
				resMap[spec.Id] = spec
				break
			}
		}
		if priceErr != nil {
			return nil, priceErr
		}
	}

	res := &VpsSpecList{}
	for _, v := range resMap {
		if len(v.AvailableZoneIds) == 0 || v.LogicalCpuNum == 0 {
			continue
		}
		res.Specs = append(res.Specs, v)
	}
	return res, nil
}

func (ac *AliyunClient) EcListSpotSpecs(regionId string) (*VpsSpecList, error) {
	cli, err := ac.getEcsClient(regionId)
	if err != nil {
		return nil, err
	}
	resMap := map[string]VpsSpec{}

	// Get hardware.
	itReq := ecs.CreateDescribeInstanceTypesRequest()
	itReq.RegionId = regionId
	ditResp, err := cli.DescribeInstanceTypes(itReq)
	if err != nil {
		return nil, err
	}
	for _, instanceType := range ditResp.InstanceTypes.InstanceType {
		memoryVolume, err := gvolume.FromByteSize(instanceType.MemorySize * float64(gvolume.GB.Bytes()))
		if err != nil {
			return nil, err
		}
		entry := VpsSpec{
			RegionId: regionId,
			Id:            instanceType.InstanceTypeId,
			LogicalCpuNum: instanceType.CpuCoreCount,
			MemoryVolume:  memoryVolume,
		}
		entry.IsCredit = gstring.StartWith(entry.Id, "ecs.t")
		resMap[instanceType.InstanceTypeId] = entry
	}

	// Get available zones.
	arrReq := ecs.CreateDescribeAvailableResourceRequest()
	arrReq.DestinationResource = "InstanceType"
	arrReq.RegionId = regionId
	darResp, err := cli.DescribeAvailableResource(arrReq)
	if err != nil {
		return nil, err
	}
	for _, zoneInfo := range darResp.AvailableZones.AvailableZone { // Iterate all zones.
		for _, resInfo := range zoneInfo.AvailableResources.AvailableResource[0].SupportedResources.SupportedResource { // Iterate all instance specs at specified zone.
			if resInfo.Status != "Available" {
				continue
			}
			originVal := resMap[resInfo.Value] // resInfo.Value is instance spec like "ecs.t6-c2m1.large".
			originVal.AvailableZoneIds = append(originVal.AvailableZoneIds, zoneInfo.ZoneId)
			resMap[resInfo.Value] = originVal
		}
	}

	// Get spot price.
	sphReq := ecs.CreateDescribeSpotPriceHistoryRequest()
	sphReq.RegionId = regionId
	sphReq.NetworkType = "vpc"
	for _, vpsSpec := range resMap {
		sphReq.InstanceType = vpsSpec.Id // Instance spec like "ecs.t6-c2m1.large".

		time.Sleep(time.Millisecond * 10) // avoid request throttling
		sphResp, err := cli.DescribeSpotPriceHistory(sphReq)
		if err != nil {
			return nil, err
		}
		if len(sphResp.SpotPrices.SpotPriceType) == 0 {
			continue
		}
		currency, err := gfiat.ParseFiat(sphResp.Currency)
		if err != nil {
			return nil, err
		}
		for _, v := range sphResp.SpotPrices.SpotPriceType {
			originVal := resMap[vpsSpec.Id]
			originVal.Currency = currency
			if originVal.OnDemandPrices == nil {
				originVal.OnDemandPrices = map[string]gdecimal.Decimal{}
			}
			if originVal.SpotPricePerHour == nil {
				originVal.SpotPricePerHour = map[string]gdecimal.Decimal{}
			}
			originVal.OnDemandPrices[v.ZoneId] = gdecimal.NewFromFloat64(v.OriginPrice)
			originVal.SpotPricePerHour[v.ZoneId] = gdecimal.NewFromFloat64(v.SpotPrice)
			resMap[vpsSpec.Id] = originVal
		}
	}

	res := &VpsSpecList{}
	for _, v := range resMap {
		if len(v.AvailableZoneIds) == 0 || v.LogicalCpuNum == 0 {
			continue
		}
		res.Specs = append(res.Specs, v)
	}
	return res, nil
}

func (ac *AliyunClient) EcListImages(regionId string) ([]SysImage, error) {
	cli, err := ac.getEcsClient(regionId)
	if err != nil {
		return nil, err
	}
	dir, err := cli.DescribeImages(ecs.CreateDescribeImagesRequest())
	if err != nil {
		return nil, err
	}
	var res []SysImage
	for _, v := range dir.Images.Image {
		entry := SysImage{
			Id:        v.ImageId,
			Name:      v.OSNameEn,
			OS:        v.OSType,
			Arch:      v.Architecture,
			Distro:    v.Platform,
			Available: v.Status == "Available",
		}
		res = append(res, entry)
	}
	return res, nil
}

func (ac *AliyunClient) EcListVps(regionId string) ([]VpsInfo, error) {
	c, err := ecs.NewClientWithAccessKey(regionId, ac.accessKey, ac.secretKey)
	if err != nil {
		return nil, err
	}
	dir, err := c.DescribeInstances(ecs.CreateDescribeInstancesRequest())
	if err != nil {
		return nil, err
	}
	var res []VpsInfo
	for _, v := range dir.Instances.Instance {
		ict, err := ac.chargeTypeAliyunToStd(v.InstanceChargeType)
		if err != nil {
			return nil, err
		}
		nct, err := ac.chargeTypeAliyunToStd(v.InstanceChargeType)
		if err != nil {
			return nil, err
		}
		v.StartTime = strings.Replace(v.StartTime, "Z", ":00Z", 1)
		startTime := gtime.ZeroTime
		if v.StartTime != "" {
			startTime, err = gtime.ParseTimeStringStrict(v.StartTime)
			if err != nil {
				return nil, err
			}
		}
		v.CreationTime = strings.Replace(v.CreationTime, "Z", ":00Z", 1)
		creationTime := gtime.ZeroTime
		if v.CreationTime != "" {
			creationTime, err = gtime.ParseTimeStringStrict(v.CreationTime)
			if err != nil {
				return nil, err
			}
		}
		v.AutoReleaseTime = strings.Replace(v.AutoReleaseTime, "Z", ":00Z", 1)
		autoReleaseTime := gtime.ZeroTime
		if v.AutoReleaseTime != "" {
			autoReleaseTime, err = gtime.ParseTimeStringStrict(v.AutoReleaseTime)
			if err != nil {
				return nil, err
			}
		}
		memoryVolume, err := gvolume.FromByteSize(float64(v.Memory) * float64(gvolume.MB.Bytes()))
		if err != nil {
			return nil, err
		}
		var privateIPs []string
		for _, nic := range v.NetworkInterfaces.NetworkInterface {
			for _, pis := range nic.PrivateIpSets.PrivateIpSet {
				privateIPs = append(privateIPs, pis.PrivateIpAddress)
			}
		}
		privateIPs = gstring.RemoveDuplicate(privateIPs)
		entry := VpsInfo{
			Id:                 v.InstanceId,
			Name:               v.InstanceName,
			Specs:              v.InstanceType,
			InstanceChargeType: ict,
			InternetChargeType: nct,
			SpotPriceLimit:     v.SpotPriceLimit,
			SpotStrategy:       v.SpotStrategy,
			SpotStartTime:      startTime,
			NetworkType:        v.InstanceNetworkType, // "classic", "vpc"
			PrivateIPs:         privateIPs,
			PublicIPs:          v.PublicIpAddress.IpAddress,
			RegionId:           v.RegionId,
			ZoneId:             v.ZoneId,
			SecurityGroupIds:   v.SecurityGroupIds.SecurityGroupId,
			ImageId:            v.ImageId,
			Status:             v.Status,
			KeyPairName:        v.KeyPairName,
			PhysicalCpuNum:     v.CpuOptions.CoreCount,
			LogicalCpuNum:      v.Cpu,
			GpuNum:             v.GPUAmount,
			CreationTime:       creationTime,
			AutoReleaseTime:    autoReleaseTime,
			MemorySize:         memoryVolume,
			SysImageId:         v.ImageId,
			SysImageName:       v.OSNameEn,
			SysImageOS:         v.OSType,
		}
		res = append(res, entry)
	}
	return res, nil
}

func (ac *AliyunClient) EcListSecurityGroups(regionId string) ([]SecurityGroup, error) {
	cli, err := ac.getEcsClient(regionId)
	if err != nil {
		return nil, err
	}
	niReq := ecs.CreateDescribeSecurityGroupsRequest()
	niReq.RegionId = regionId
	niResp, err := cli.DescribeSecurityGroups(niReq)
	if err != nil {
		return nil, err
	}

	var res []SecurityGroup
	for _, v := range niResp.SecurityGroups.SecurityGroup {
		item := SecurityGroup{
			Id:               v.SecurityGroupId,
			Name:             v.SecurityGroupName,
		}
		gaReq := ecs.CreateDescribeSecurityGroupAttributeRequest()
		gaReq.SecurityGroupId = v.SecurityGroupId
		gaResp, err := cli.DescribeSecurityGroupAttribute(gaReq)
		if err != nil {
			return nil, err
		}
		for _, aliPermission := range gaResp.Permissions.Permission {
			stdPermission := SecurityPermission{
				Description: aliPermission.Description,
				Direction:   aliPermission.Direction,
				Protocol:    aliPermission.IpProtocol,
			}
			if aliPermission.SourcePortRange != "" {
				srcPortRangeSS := strings.Split(aliPermission.SourcePortRange, "/")
				if len(srcPortRangeSS) != 2 {
					return nil, gerrors.New("invalid SourcePortRange(%s)", aliPermission.SourcePortRange)
				}
				srcPortBegin, err := strconv.ParseInt(srcPortRangeSS[0], 10, 64)
				if err != nil {
					return nil, err
				}
				srcPortEnd, err := strconv.ParseInt(srcPortRangeSS[1], 10, 64)
				if err != nil {
					return nil, err
				}
				stdPermission.SrcPortRange = [2]int{int(srcPortBegin), int(srcPortEnd)}
				stdPermission.SrcCidrIP = aliPermission.SourceCidrIp
			}
			if aliPermission.PortRange != "" {
				dstPortRangeSS := strings.Split(aliPermission.PortRange, "/")
				if len(dstPortRangeSS) != 2 {
					return nil, gerrors.New("invalid PortRange(%s)", aliPermission.PortRange)
				}
				dstPortBegin, err := strconv.ParseInt(dstPortRangeSS[0], 10, 64)
				if err != nil {
					return nil, err
				}
				dstPortEnd, err := strconv.ParseInt(dstPortRangeSS[1], 10, 64)
				if err != nil {
					return nil, err
				}
				stdPermission.DstPortRange = [2]int{int(dstPortBegin), int(dstPortEnd)}
				stdPermission.DstCidrIP = aliPermission.SourceCidrIp
			}
			item.Permissions = append(item.Permissions, stdPermission)
		}
		res = append(res, item)
	}
	return res, nil
}

func (ac *AliyunClient) EcCreateSecurityGroup(regionId string, sg SecurityGroup) (string, error) {
	cli, err := ac.getEcsClient(regionId)
	if err != nil {
		return "", err
	}

	// If existed, delete it first.
	niReq := ecs.CreateDescribeSecurityGroupsRequest()
	niReq.RegionId = regionId
	var allRegions [][2]string // [][2]{Name, Id}
	niResp, err := cli.DescribeSecurityGroups(niReq)
	if err != nil {
		return "", err
	}
	for _, v := range niResp.SecurityGroups.SecurityGroup {
		allRegions = append(allRegions, [2]string{v.SecurityGroupName, v.SecurityGroupId})
	}
	for _, v := range allRegions {
		if v[0] == sg.Name {
			if err := ac.EcDeleteSecurityGroup(regionId, v[1]); err != nil {
				return "", err
			}
		}
	}

	// Create new security group.
	csReq := ecs.CreateCreateSecurityGroupRequest()
	csReq.RegionId = regionId
	csReq.SecurityGroupName = sg.Name
	csResp, err := cli.CreateSecurityGroup(csReq)
	if err != nil {
		return "", err
	}
	sgId := csResp.SecurityGroupId

	// Create permmisions for security group.
	for _, permission := range sg.Permissions {
		if strings.ToLower(permission.Direction) == "in" {
			asReq := ecs.CreateAuthorizeSecurityGroupRequest()
			asReq.RegionId = regionId
			asReq.SecurityGroupId = sg.Id
			asReq.IpProtocol = permission.Protocol
			asReq.SourceCidrIp = permission.SrcCidrIP
			asReq.SourcePortRange = fmt.Sprintf("%d/%d", permission.SrcPortRange[0], permission.SrcPortRange[1])
			asReq.DestCidrIp = permission.DstCidrIP
			asReq.PortRange = fmt.Sprintf("%d/%d", permission.DstPortRange[0], permission.DstPortRange[1])
			if _, err := cli.AuthorizeSecurityGroup(asReq); err != nil {
				ac.EcDeleteSecurityGroup(regionId, sgId)
				return "", err
			}
		} else if strings.ToLower(permission.Direction) == "out" {
			asReq := ecs.CreateAuthorizeSecurityGroupEgressRequest()
			asReq.RegionId = regionId
			asReq.SecurityGroupId = sg.Id
			asReq.IpProtocol = permission.Protocol
			asReq.SourceCidrIp = permission.SrcCidrIP
			asReq.SourcePortRange = fmt.Sprintf("%d/%d", permission.SrcPortRange[0], permission.SrcPortRange[1])
			asReq.DestCidrIp = permission.DstCidrIP
			asReq.PortRange = fmt.Sprintf("%d/%d", permission.DstPortRange[0], permission.DstPortRange[1])
			if _, err := cli.AuthorizeSecurityGroupEgress(asReq); err != nil {
				ac.EcDeleteSecurityGroup(regionId, sgId)
				return "", err
			}
		} else {
			ac.EcDeleteSecurityGroup(regionId, sgId)
			return "", gerrors.New("unsupported direction(%s)", permission.Direction)
		}
	}

	return sgId, nil
}

func (ac *AliyunClient) EcDeleteSecurityGroup(regionId, securityGroupId string) error {
	cli, err := ac.getEcsClient(regionId)
	if err != nil {
		return err
	}
	dsReq := ecs.CreateDeleteSecurityGroupRequest()
	dsReq.RegionId = regionId
	dsReq.SecurityGroupId = securityGroupId
	_, err = cli.DeleteSecurityGroup(dsReq)
	return err
}

func (ac *AliyunClient) EcListSwitches(regionId, zoneId string) ([]string, error) {
	cli, err := ac.getEcsClient(regionId)
	if err != nil {
		return nil, err
	}
	niReq := ecs.CreateDescribeNetworkInterfacesRequest()
	niReq.RegionId = regionId
	niResp, err := cli.DescribeNetworkInterfaces(niReq)
	if err != nil {
		return nil, err
	}

	var res []string
	for _, v := range niResp.NetworkInterfaceSets.NetworkInterfaceSet {
		if v.ZoneId != zoneId {
			continue
		}
		res = append(res, v.VSwitchId)
	}
	res = gstring.RemoveDuplicate(res)
	return res, nil
}

func (ac *AliyunClient) EcCreateVps(regionId string, tmpl VpsCreationTmpl) (string, error) {
	c, err := ac.getEcsClient(regionId)
	if err != nil {
		return "", err
	}
	vpsCharge, err := ac.chargeTypeStdToAliyun(tmpl.VpsCharge)
	if err != nil {
		return "", err
	}
	ciReq := ecs.CreateCreateInstanceRequest()
	ciReq.RegionId = regionId
	ciReq.KeyPairName = tmpl.KeyPair
	ciReq.Password = tmpl.Password
	ciReq.ZoneId = tmpl.ZoneId
	ciReq.InstanceType = tmpl.Specs
	ciReq.InstanceName = tmpl.Name
	ciReq.ImageId = tmpl.ImageId
	ciReq.SecurityGroupId = tmpl.SecurityGroupId
	// Only "vpc" network type supported now, "classic" network type is fading away.
	// At "vpc" network type, switch Id required always.
	ciReq.VSwitchId = tmpl.SwitchId
	ciReq.SystemDiskSize = requests.NewInteger(tmpl.SystemDiskGB)
	ciReq.InstanceChargeType = vpsCharge
	switch tmpl.InternetCharge {
	case PayByBandwidth:
		ciReq.InternetChargeType = "PayByBandwidth"
	case PayByTraffic:
		ciReq.InternetChargeType = "PayByTraffic"
	default:
		return "", gerrors.New("unsupported InternetCharge(%s)", tmpl.InternetCharge)
	}
	ciReq.InternetMaxBandwidthIn = requests.NewInteger(tmpl.BandWidthMbIn)
	ciReq.InternetMaxBandwidthOut = requests.NewInteger(tmpl.BandWidthMbOut)
	ciReq.SpotStrategy = "NoSpot"
	if tmpl.UnlimitedPerformance {
		ciReq.CreditSpecification = "Unlimited"
	} else {
		ciReq.CreditSpecification = "Standard"
	}
	if tmpl.IsSpot {
		if tmpl.SpotStrategy == SpotWithMaxPrice {
			ciReq.SpotPriceLimit = requests.NewFloat(tmpl.SpotMaxPrice.Float64())
		}
		if tmpl.OpenSpotDuration {
			ciReq.SpotDuration = requests.NewInteger(1) // At aliyun, spot duration is an integer in hour. Support 1 for now.
		}
		ciReq.SpotInterruptionBehavior = "Terminate" // "Terminate" supported for now.
		if tmpl.SpotStrategy == SpotWithMaxPrice {
			ciReq.SpotStrategy = "SpotWithPriceLimit"
		} else if tmpl.SpotStrategy == SpotAsPriceGo {
			ciReq.SpotStrategy = "SpotAsPriceGo"
		} else {
			return "", gerrors.New("unsupported spot strategy(%s)", tmpl.SpotStrategy)
		}
	}
	response, err := c.CreateInstance(ciReq)
	if err != nil {
		return "", err
	}
	return response.InstanceId, nil
}

func (ac *AliyunClient) EcStartVps(regionId, vpsId string) error {
	cli, err := ac.getEcsClient(regionId)
	if err != nil {
		return err
	}
	req := ecs.CreateStartInstanceRequest()
	req.RegionId = regionId
	req.InstanceId = vpsId
	_, err = cli.StartInstance(req)
	return err
}

func (ac *AliyunClient) EcDeleteVps(regionId string, vpsIds []string, force bool) error {
	ec, err := ac.getEcsClient(regionId)
	if err != nil {
		return err
	}
	diReq := ecs.CreateDeleteInstancesRequest()
	var idList []string
	for _, id := range vpsIds {
		idList = append(idList, id)
	}
	diReq.InstanceId = &idList
	diReq.Force = requests.NewBoolean(force)
	_, err = ec.DeleteInstances(diReq)
	return err
}

func (ac *AliyunClient) ObsIsBucketExist(regionId, bucketName string) (bool, error) {
	ossClient, err := ac.getOssClient(regionId)
	if err != nil {
		return false, err
	}
	return ossClient.IsBucketExist(bucketName)
}

func (ac *AliyunClient) ObsNewBucket(regionId, bucketName string) error {
	ossClient, err := ac.getOssClient(regionId)
	if err != nil {
		return err
	}
	return ossClient.CreateBucket(bucketName, oss.ACL(oss.ACLPrivate))
}

func (ac *AliyunClient) ObsDeleteBucket(regionId, bucketName string, deleteIfNotEmpty *bool) error {
	ossClient, err := ac.getOssClient(regionId)
	if err != nil {
		return err
	}

	if deleteIfNotEmpty != nil && *deleteIfNotEmpty {
		allKeys, err := ac.ObsListObjectKeys(regionId, bucketName)
		if err != nil {
			return err
		}
		for _, key := range allKeys {
			if err := ac.ObsDeleteObject(regionId, bucketName, key); err != nil {
				return err
			}
		}
	}

	return ossClient.DeleteBucket(bucketName)
}

func (ac *AliyunClient) ObsListBuckets(regionId string) ([]string, error) {
	ossClient, err := ac.getOssClient(regionId)
	if err != nil {
		return nil, err
	}
	properties, err := ossClient.ListBuckets()
	if err != nil {
		return nil, err
	}

	var res []string
	for _, v := range properties.Buckets {
		res = append(res, v.Name)
	}
	return res, nil
}

func (ac *AliyunClient) ObsListObjectKeys(regionId string, bucketName string) ([]string, error) {
	bucket, err := ac.getBucket(regionId, bucketName)
	if err != nil {
		return nil, err
	}

	maxKeys := 1000 // Max 1000 objects returned for each list request.
	pageToken := ""
	var res []string
	for {
		options := []oss.Option{oss.MaxKeys(maxKeys)}
		if pageToken != "" {
			options = append(options, oss.ContinuationToken(pageToken))
		}
		objs, err := bucket.ListObjectsV2(options...)
		if err != nil {
			return nil, err
		}
		for _, v := range objs.Objects {
			res = append(res, v.Key)
		}
		// Note: don't use 'pageToken = objs.ContinuationToken'
		pageToken = objs.NextContinuationToken
		if pageToken == "" {
			break
		}
	}

	return res, nil
}

func (ac *AliyunClient) ObsScanObjectKeys(regionId, bucketName, pageToken string) ([]string, string, error) {
	bucket, err := ac.getBucket(regionId, bucketName)
	if err != nil {
		return nil, "", err
	}

	maxKeys := 1000 // Max 1000 objects returned for each list request.
	options := []oss.Option{oss.MaxKeys(maxKeys)}
	if pageToken != "" {
		options = append(options, oss.ContinuationToken(pageToken))
	}
	objs, err := bucket.ListObjectsV2(options...)
	if err != nil {
		return nil, "", err
	}
	var res []string
	for _, v := range objs.Objects {
		res = append(res, v.Key)
	}
	// Note: don't use 'objs.ContinuationToken'
	nextPageToken := objs.NextContinuationToken
	return res, nextPageToken, nil
}

func (ac *AliyunClient) ObsUpsertObject(regionId, bucketName, blobId string, blob []byte) error {
	bucket, err := ac.getBucket(regionId, bucketName)
	if err != nil {
		return err
	}
	br := gio.BytesToReadCloser(blob)
	defer func() { _ = br.Close() }()
	return bucket.PutObject(blobId, br)
}

func (ac *AliyunClient) ObsGetObject(regionId, bucketName, blobId string) ([]byte, error) {
	bucket, err := ac.getBucket(regionId, bucketName)
	if err != nil {
		return nil, err
	}
	rc, err := bucket.GetObject(blobId)
	if err != nil {
		exist, err := bucket.IsObjectExist(blobId)
		if err != nil {
			return nil, err
		}
		if !exist {
			return nil, gerrors.ErrNotExist
		}
		return nil, err
	}
	defer func() { _ = rc.Close() }()
	return gio.ReadCloserToBytes(rc)
}

func (ac *AliyunClient) ObsGetObjectSize(regionId, bucketName, blobId string) (*gvolume.Volume, error) {
	bucket, err := ac.getBucket(regionId, bucketName)
	if err != nil {
		return nil, err
	}
	head, err := bucket.GetObjectMeta(blobId)
	if err != nil {
		return nil, err
	}
	sizeStr, ok := head["Content-Length"]
	if !ok || len(sizeStr) == 0 {
		return nil, gerrors.New("Content-Length not found in response of GetObjectMeta")
	}
	sizeBytes, err := strconv.ParseInt(sizeStr[0], 10, 64)
	if err != nil {
		return nil, err
	}
	vol, err := gvolume.FromByteSize(float64(sizeBytes))
	if err != nil {
		return nil, err
	}
	return &vol, nil
}

func (ac *AliyunClient) ObsDeleteObject(regionId, bucketName, blobId string) error {
	bucket, err := ac.getBucket(regionId, bucketName)
	if err != nil {
		return err
	}
	return bucket.DeleteObject(blobId)
}

func (ac *AliyunClient) ObsRenameObject(regionId, bucketName, oldObjectKey, newObjectKey string) error {
	bucket, err := ac.getBucket(regionId, bucketName)
	if err != nil {
		return err
	}
	_, err = bucket.CopyObject(oldObjectKey, newObjectKey)
	if err != nil {
		return err
	}
	if err = bucket.DeleteObject(oldObjectKey); err != nil {
		return err
	}
	return nil
}

func (ac *AliyunClient) SmsSendTmpl(regionId string, mobiles []string, sign, tmpl string, params map[string]string) error {
	client, err := dysmsapi.NewClientWithAccessKey(regionId, ac.accessKey, ac.secretKey)
	if err != nil {
		return err
	}
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = strings.Join(mobiles, ",")
	request.SignName = sign
	request.TemplateCode = tmpl
	request.TemplateParam, err = gjson.MarshalString(params, false)
	if err != nil {
		return err
	}
	_, err = client.SendSms(request)
	return err
}

func (ac *AliyunClient) SmsSendMsg(fromMobile, toMobile, message string) (err error) {
	return gerrors.ErrNotImplemented
}

func (ac *AliyunClient) Close() error {
	return nil
}

func (ac *AliyunClient) getEcsClient(regionId string) (*ecs.Client, error) {
	var res *ecs.Client = nil
	ok := false
	ac.ecsListMu.RLock()
	res, ok = ac.ecsList[regionId]
	ac.ecsListMu.RUnlock()
	if ok {
		return res, nil
	}
	res, err := ecs.NewClientWithAccessKey(regionId, ac.accessKey, ac.secretKey)
	if err != nil {
		return nil, err
	}
	ac.ecsListMu.Lock()
	ac.ecsList[regionId] = res
	ac.ecsListMu.Unlock()
	return res, nil
}

func (ac *AliyunClient) getOssClient(regionId string) (*oss.Client, error) {
	var res *oss.Client = nil
	ok := false
	ac.ossListMu.RLock()
	res, ok = ac.ossList[regionId]
	ac.ossListMu.RUnlock()
	if ok {
		return res, nil
	}
	endpoint, err := ac.getEndpoint(regionId, OBS, ac.lan)
	if err != nil {
		return nil, err
	}
	res, err = oss.New(endpoint, ac.accessKey, ac.secretKey)
	if err != nil {
		return nil, err
	}
	ac.ossListMu.Lock()
	ac.ossList[regionId] = res
	ac.ossListMu.Unlock()
	return res, nil
}

func (ac *AliyunClient) getBucket(regionId string, bucketName string) (*oss.Bucket, error) {
	ossClient, err := ac.getOssClient(regionId)
	if err != nil {
		return nil, err
	}

	var bucket *oss.Bucket = nil
	ok := false
	ac.bucketsMu.RLock()
	bucket, ok = ac.buckets[bucketName]
	ac.bucketsMu.RUnlock()
	if ok {
		return bucket, nil
	}

	bucket, err = ossClient.Bucket(bucketName)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

/*func (ac *AliyunClient) listRegionsByJSON() ([]string, error) {
	eps, err := ac.getEndpoints()
	if err != nil {
		return nil, err
	}
	var res []string
	for k := range eps {
		res = append(res, k)
	}
	return res, nil
}*/

func (ac *AliyunClient) getEndpoints() (map[string]RegionInfo, error) {
	res := map[string]RegionInfo{}
	if err := json.Unmarshal([]byte(aliyunJsonString), &res); err != nil {
		return nil, err
	}
	return res, nil
}

// Doc: https://help.aliyun.com/document_detail/31837.html?spm=a2c4g.11186623.2.18.109d6685wSGpx7#concept-zt4-cvy-5db
func (ac *AliyunClient) getEndpoint(regionId string, production Production, LAN bool) (string, error) {
	epl, err := ac.getEndpoints()
	if err != nil {
		return "", err
	}
	regionInfo, ok := epl[regionId]
	if !ok {
		return "", gerrors.New("region %s not found", regionId)
	}
	ep, ok := regionInfo.Endpoints[production]
	if !ok {
		return "", gerrors.New("region %s production %s not found", regionId, production)
	}
	if LAN {
		return ep.LAN, nil
	}
	return ep.WAN, nil
}

func (ac *AliyunClient) chargeTypeAliyunToStd(chargeType string) (ChargeType, error) {
	switch strings.ToLower(chargeType) {
	case "prepaid", "postpaid":
		return ChargeType(strings.ToLower(chargeType)), nil
	default:
		return "", gerrors.New("unsupported charge type %s", chargeType)
	}
}

func (ac *AliyunClient) chargeTypeStdToAliyun(chargeType ChargeType) (string, error) {
	switch chargeType {
	case PrePaid:
		return "PrePaid", nil
	case PostPaid:
		return "PostPaid", nil
	default:
		return "", gerrors.New("unsupported charge type %s", chargeType)
	}
}

var (
	aliyunJsonString = `{
	"cn-hangzhou": {
		"Country": "CHN",
		"City": "Hangzhou",
		"Endpoints": {
			"obs": {
				"WAN": "oss-cn-hangzhou.aliyuncs.com",
				"LAN": "oss-cn-hangzhou-internal.aliyuncs.com"
			},
			"fc":{
				"WAN":"cn-hangzhou.fc.aliyuncs.com",
				"LAN":"cn-hangzhou-internal.fc.aliyuncs.com"
			}
		}
	},
	"cn-shanghai": {
		"Country": "CHN",
		"City": "Shanghai",
		"Endpoints": {
			"obs": {
				"WAN": "oss-cn-shanghai.aliyuncs.com",
				"LAN": "oss-cn-shanghai-internal.aliyuncs.com"
			},
			"fc":{
				"WAN":"cn-shanghai.fc.aliyuncs.com",
				"LAN":"cn-shanghai-internal.fc.aliyuncs.com"
			}
		}
	},
	"cn-qingdao": {
		"Country": "CHN",
		"City": "Qingdao",
		"Endpoints": {
			"obs": {
				"WAN": "oss-cn-qingdao.aliyuncs.com",
				"LAN": "oss-cn-qingdao-internal.aliyuncs.com"
			},
			"fc":{
				"WAN":"cn-qingdao.fc.aliyuncs.com",
				"LAN":"cn-qingdao-internal.fc.aliyuncs.com"
			}
		}
	},
	"cn-beijing": {
		"Country": "CHN",
		"City": "Beijing",
		"Endpoints": {
			"obs": {
				"WAN": "oss-cn-beijing.aliyuncs.com",
				"LAN": "oss-cn-beijing-internal.aliyuncs.com"
			},
			"fc":{
				"WAN":"cn-beijing.fc.aliyuncs.com",
				"LAN":"cn-beijing-internal.fc.aliyuncs.com"
			}
		}
	},
	"cn-zhangjiakou": {
		"Country": "CHN",
		"City": "Zhangjiakou",
		"Endpoints": {
			"obs": {
				"WAN": "oss-cn-zhangjiakou.aliyuncs.com",
				"LAN": "oss-cn-zhangjiakou-internal.aliyuncs.com"
			},
			"fc":{
				"WAN":"cn-zhangjiakou.fc.aliyuncs.com",
				"LAN":"cn-zhangjiakou-internal.fc.aliyuncs.com"
			}
		}
	},
	"cn-huhehaote": {
		"Country": "CHN",
		"City": "Huhehaote",
		"Endpoints": {
			"obs": {
				"WAN": "oss-cn-huhehaote.aliyuncs.com",
				"LAN": "oss-cn-huhehaote-internal.aliyuncs.com"
			},
			"fc":{
				"WAN":"cn-huhehaote.fc.aliyuncs.com",
				"LAN":"cn-huhehaote-internal.fc.aliyuncs.com"
			}
		}
	},
	"cn-wulanchabu": {
		"Country": "CHN",
		"City": "Wulanchabu",
		"Endpoints": {
			"obs": {
				"WAN": "oss-cn-wulanchabu.aliyuncs.com",
				"LAN": "oss-cn-wulanchabu-internal.aliyuncs.com"
			}
		}
	},
	"cn-shenzhen": {
		"Country": "CHN",
		"City": "Shenzhen",
		"Endpoints": {
			"obs": {
				"WAN": "oss-cn-shenzhen.aliyuncs.com",
				"LAN": "oss-cn-shenzhen-internal.aliyuncs.com"
			},
			"fc":{
				"WAN":"cn-shenzhen.fc.aliyuncs.com",
				"LAN":"cn-shenzhen-internal.fc.aliyuncs.com"
			}
		}
	},
	"cn-heyuan": {
		"Country": "CHN",
		"City": "Heyuan",
		"Endpoints": {
			"obs": {
				"WAN": "oss-cn-heyuan.aliyuncs.com",
				"LAN": "oss-cn-heyuan-internal.aliyuncs.com"
			}
		}
	},
	"cn-guangzhou": {
		"Country": "CHN",
		"City": "Guangzhou",
		"Endpoints": {
			"obs": {
				"WAN": "oss-cn-guangzhou.aliyuncs.com",
				"LAN": "oss-cn-guangzhou-internal.aliyuncs.com"
			}
		}
	},
	"cn-chengdu": {
		"Country": "CHN",
		"City": "Chengdu",
		"Endpoints": {
			"obs": {
				"WAN": "oss-cn-chengdu.aliyuncs.com",
				"LAN": "oss-cn-chengdu-internal.aliyuncs.com"
			},
			"fc":{
				"WAN":"cn-chengdu.fc.aliyuncs.com",
				"LAN":"cn-chengdu-internal.fc.aliyuncs.com"
			}
		}
	},
	"cn-hongkong": {
		"Country": "CHN",
		"City": "Hongkong",
		"Endpoints": {
			"obs": {
				"WAN": "oss-cn-hongkong.aliyuncs.com",
				"LAN": "oss-cn-hongkong-internal.aliyuncs.com"
			},
			"fc":{
				"WAN":"cn-hongkong.fc.aliyuncs.com",
				"LAN":"cn-hongkong-internal.fc.aliyuncs.com"
			}
		}
	},
	"us-west-1": {
		"Country": "USA",
		"City": "San Mateo",
		"Endpoints": {
			"obs": {
				"WAN": "oss-us-west-1.aliyuncs.com",
				"LAN": "oss-us-west-1-internal.aliyuncs.com"
			},
			"fc":{
				"WAN":"us-west-1.fc.aliyuncs.com",
				"LAN":"us-west-1-internal.fc.aliyuncs.com"
			}
		}
	},
	"us-east-1": {
		"Country": "USA",
		"City": "Ashburn",
		"Endpoints": {
			"obs": {
				"WAN": "oss-us-east-1.aliyuncs.com",
				"LAN": "oss-us-east-1-internal.aliyuncs.com"
			},
			"fc":{
				"WAN":"us-east-1.fc.aliyuncs.com",
				"LAN":"us-east-1-internal.fc.aliyuncs.com"
			}
		}
	},
	"ap-southeast-1": {
		"Country": "SG",
		"City": "Singapore",
		"Endpoints": {
			"obs": {
				"WAN": "oss-ap-southeast-1.aliyuncs.com",
				"LAN": "oss-ap-southeast-1-internal.aliyuncs.com"
			},
			"fc":{
				"WAN":"ap-southeast-1.fc.aliyuncs.com",
				"LAN":"ap-southeast-1-internal.fc.aliyuncs.com"
			}
		}
	},
	"ap-southeast-2": {
		"Country": "AU",
		"City": "Sydney",
		"Endpoints": {
			"obs": {
				"WAN": "oss-ap-southeast-2.aliyuncs.com",
				"LAN": "oss-ap-southeast-2-internal.aliyuncs.com"
			},
			"fc":{
				"WAN":"ap-southeast-2.fc.aliyuncs.com",
				"LAN":"ap-southeast-2-internal.fc.aliyuncs.com"
			}
		}
	},
	"ap-southeast-3": {
		"Country": "MY",
		"City": "Kuala Lumpur",
		"Endpoints": {
			"obs": {
				"WAN": "oss-ap-southeast-3.aliyuncs.com",
				"LAN": "oss-ap-southeast-3-internal.aliyuncs.com"
			},
			"fc":{
				"WAN":"ap-southeast-3.fc.aliyuncs.com",
				"LAN":"ap-southeast-3-internal.fc.aliyuncs.com"
			}
		}
	},
	"ap-southeast-5": {
		"Country": "IDN",
		"City": "Jakarta",
		"Endpoints": {
			"obs": {
				"WAN": "oss-ap-southeast-5.aliyuncs.com",
				"LAN": "oss-ap-southeast-5-internal.aliyuncs.com"
			},
			"fc":{
				"WAN":"ap-southeast-5.fc.aliyuncs.com",
				"LAN":"ap-southeast-5-internal.fc.aliyuncs.com"
			}
		}
	},
	"ap-northeast-1": {
		"Country": "JP",
		"City": "Tokyo",
		"Endpoints": {
			"obs": {
				"WAN": "oss-ap-northeast-1.aliyuncs.com",
				"LAN": "oss-ap-northeast-1-internal.aliyuncs.com"
			},
			"fc":{
				"WAN":"ap-northeast-1.fc.aliyuncs.com",
				"LAN":"ap-northeast-1-internal.fc.aliyuncs.com"
			}
		}
	},
	"ap-south-1": {
		"Country": "IRI",
		"City": "Mumbai",
		"Endpoints": {
			"obs": {
				"WAN": "oss-ap-south-1.aliyuncs.com",
				"LAN": "oss-ap-south-1-internal.aliyuncs.com"
			},
			"fc":{
				"WAN":"ap-south-1.fc.aliyuncs.com",
				"LAN":"ap-south-1-internal.fc.aliyuncs.com"
			}
		}
	},
	"eu-central-1": {
		"Country": "GER",
		"City": "Frankfurt",
		"Endpoints": {
			"obs": {
				"WAN": "oss-eu-central-1.aliyuncs.com",
				"LAN": "oss-eu-central-1-internal.aliyuncs.com"
			},
			"fc":{
				"WAN":"eu-central-1.fc.aliyuncs.com",
				"LAN":"eu-central-1-internal.fc.aliyuncs.com"
			}
		}
	},
	"eu-west-1": {
		"Country": "UK",
		"City": "London",
		"Endpoints": {
			"obs": {
				"WAN": "oss-eu-west-1.aliyuncs.com",
				"LAN": "oss-eu-west-1-internal.aliyuncs.com"
			},
			"fc":{
				"WAN":"eu-west-1.fc.aliyuncs.com",
				"LAN":"eu-west-1-internal.fc.aliyuncs.com"
			}
		}
	},
	"me-east-1": {
		"Country": "ARE",
		"City": "Dubai",
		"Endpoints": {
			"obs": {
				"WAN": "oss-me-east-1.aliyuncs.com",
				"LAN": "oss-me-east-1-internal.aliyuncs.com"
			}
		}
	}
}`
)
