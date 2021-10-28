package cloud

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gdecimal"
	"github.com/cryptowilliam/goutil/container/gvolume"
	"sort"
	"time"
)

type (
	// How to charge.
	ChargeType string

	// Internet charge type.
	InternetChargeType string

	// Spot strategy.
	SpotStrategy string

	// Instance specification.
	VpsSpec struct {
		RegionId string
		Id               string
		IsCredit bool 	// Is credit instance(突发性能实例) or normal instance.
		Currency         string
		LogicalCpuNum    int
		MemoryVolume     gvolume.Volume
		AvailableZoneIds []string
		OnDemandPrices   map[string]gdecimal.Decimal
		SpotPricePerHour map[string]gdecimal.Decimal
	}

	VpsSpecEx struct {
		RegionId 		 string
		ZoneId           string
		Id               string
		IsCredit bool 	// Is credit instance(突发性能实例) or normal instance.
		Currency         string
		LogicalCpuNum    int
		MemoryVolume     gvolume.Volume
		Network          string
		OnDemandPrices   gdecimal.Decimal
		SpotPricePerHour gdecimal.Decimal
	}

	VpsSpecList struct {
		Specs []VpsSpec
	}

	VpsSpecExList struct {
		SpecExs []VpsSpecEx
	}

	// Instance is a nested struct in ecs response.
	VpsInfo struct {
		Id                 string
		Name               string
		Specs              string // "ecs.g5.large"
		InstanceChargeType ChargeType
		InternetChargeType ChargeType
		SpotPriceLimit     float64
		SpotStrategy       string
		SpotStartTime      time.Time
		NetworkType        string // "classic", "vpc"
		PrivateIPs         []string
		PublicIPs          []string
		RegionId           string
		ZoneId             string
		SecurityGroupIds   []string
		ImageId            string
		Status             string
		KeyPairName        string
		PhysicalCpuNum     int
		LogicalCpuNum      int
		GpuNum             int
		CreationTime       time.Time
		AutoReleaseTime    time.Time
		MemorySize         gvolume.Volume
		SysImageId         string
		SysImageName       string
		SysImageOS         string
	}

	// Create VPS template.
	VpsCreationTmpl struct {
		// Zone Id like "cn-hangzhou-c"
		ZoneId string
		// Start a spot instance.
		IsSpot bool
		// Open spot instance protection duration, only available on aliyun for now.
		// If open "OpenSpotDuration", price will be about a little (10% on aliyun) higher than closing it.
		OpenSpotDuration bool
		// "SwitchId" required in "vpc" network type, and now only "vpc" network type supported on aliyun,
		// "classic" network type is fading away, so "SwitchId" is required.
		SwitchId string
		// Instance specs like "ecs.g5.large".
		Specs string
		// Instance name like "my-instance-temp-name".
		Name string
		// OS image Id like "ubuntu_20_04_x64_20G_alibase_20201120.vhd".
		ImageId         string
		SecurityGroupId string
		KeyPair         string
		Password        string
		SystemDiskGB    int
		VpsCharge       ChargeType
		InternetCharge  InternetChargeType
		BandWidthMbIn   int
		BandWidthMbOut  int
		SpotStrategy    SpotStrategy
		SpotMaxPrice    gdecimal.Decimal
		// true: unlimited performance mode, false: limited performance mode.
		// Available param option on aliyun.
		UnlimitedPerformance bool
	}

	SecurityPermission struct {
		Description string
		Direction string // "in", "out"
		Protocol string // "tcp","udp"...
		SrcPortRange [2]int
		SrcCidrIP string
		DstPortRange [2]int
		DstCidrIP string
	}

	SecurityGroup struct {
		Id string
		Name string
		Permissions []SecurityPermission
	}
)

var (
	PrePaid  = ChargeType("prepaid")  // Yearly package or monthly package.
	PostPaid = ChargeType("postpaid") // Pay on demand.

	PayByBandwidth = InternetChargeType("PayByBandwidth")
	PayByTraffic   = InternetChargeType("PayByTraffic")

	SpotWithMaxPrice = SpotStrategy("SpotWithMaxPrice")
	SpotAsPriceGo    = SpotStrategy("SpotAsPriceGo")
)

func (t VpsCreationTmpl) Verify(platform Platform) error {
	switch platform {
	case Aliyun:
		if t.IsSpot {
			if t.VpsCharge != PostPaid {
				return gerrors.New("invalid spot charge: VpsCharge %s", t.VpsCharge)
			}
		}
	}
	return nil
}

func (vsl *VpsSpecList) ToSpecExList() *VpsSpecExList {
	res := &VpsSpecExList{}
	for _, spec := range vsl.Specs {
		for zoneId, spotPrice := range spec.SpotPricePerHour {
			entry := VpsSpecEx{
				RegionId:		  spec.RegionId,
				ZoneId:           zoneId,
				Id:               spec.Id,
				IsCredit: 		  spec.IsCredit,
				Currency:         spec.Currency,
				LogicalCpuNum:    spec.LogicalCpuNum,
				MemoryVolume:     spec.MemoryVolume,
				Network:          "vpc",
				OnDemandPrices:   spec.OnDemandPrices[zoneId],
				SpotPricePerHour: spotPrice,
			}
			res.SpecExs = append(res.SpecExs, entry)
		}
	}
	return res
}

func (vl VpsSpecExList) Len() int {
	return len(vl.SpecExs)
}

func (vl VpsSpecExList) Less(i, j int) bool {
	iCpuNum := gdecimal.NewFromInt(vl.SpecExs[i].LogicalCpuNum)
	jCpuNum := gdecimal.NewFromInt(vl.SpecExs[j].LogicalCpuNum)
	if iCpuNum.IsZero() || jCpuNum.IsZero() {
		return false
	}
	iSpotPrice := vl.SpecExs[i].SpotPricePerHour
	jSpotPrice := vl.SpecExs[j].SpotPricePerHour
	return iSpotPrice.Div(iCpuNum).LessThan(jSpotPrice.Div(jCpuNum))
}

func (vl VpsSpecExList) Swap(i, j int) {
	vl.SpecExs[i], vl.SpecExs[j] = vl.SpecExs[j], vl.SpecExs[i]
}

func (vl *VpsSpecExList) RemoveCreditInstance() *VpsSpecExList {
	res := &VpsSpecExList{}
	for _, v := range vl.SpecExs {
		if v.IsCredit {
			continue
		}
		res.SpecExs = append(res.SpecExs, v)
	}
	return res
}

func (vl *VpsSpecExList) Sort() *VpsSpecExList {
	sort.Sort(vl)
	return vl
}

func (vl *VpsSpecExList) Append(newList *VpsSpecExList) *VpsSpecExList {
	for _, v := range newList.SpecExs {
		vl.SpecExs = append(vl.SpecExs, v)
	}
	return vl
}
