package cloud

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gdecimal"
	"github.com/cryptowilliam/goutil/container/gvolume"
)

/**
AWS S3，新建从未访问过的Key是强一致性的，新建已访问过的Key、Delete、Modify、ListBucket等其他操作都是最终一致性eventual consistency的
https://segmentfault.com/a/1190000022079593

各大平台特性对比
https://www.cnblogs.com/fastone/p/11766161.html

除了AWS S3，其他平台的对象存储都是强一致性的。
*/

const (
	// ACLPrivate definition : private read and write
	ACLPrivate ACLType = "private"
	// ACLPublicRead definition : public read and private write
	ACLPublicRead ACLType = "public-read"
	// ACLPublicReadWrite definition : public read and public write
	ACLPublicReadWrite ACLType = "public-read-write"
)

type (
	Platform string

	Production string

	PlatformInfo struct {
		SupportPrepaid         bool
		SupportPostpaid        bool
		IsGetsStrongConsist    bool
		IsNonGetsStrongConsist bool
	}

	Balance struct {
		Currency  string
		Available gdecimal.Decimal
	}

	SysImage struct {
		Id        string
		Name      string
		OS        string // "linux", "windows"...
		Distro    string // "debian", "ubuntu"...
		Arch      string // "x86_64"...
		Available bool
	}

	// TODO：RegionId 要不要优化掉
	// AWS Lambda的demo
	// https://github.com/razeone/serverless-blog-api/blob/master/hello/main.go
	Cloud interface {
		ListRegions() ([]string, error)
		ListZones(regionId string) ([]string, error)
		GetBalance() (*Balance, error)
		EcListSpotSpecs(regionId string) (*VpsSpecList, error)
		EcListImages(regionId string) ([]SysImage, error)
		EcListVps(regionId string) ([]VpsInfo, error)
		EcListSecurityGroups(regionId string) ([]SecurityGroup, error)
		EcCreateSecurityGroup(regionId string, sg SecurityGroup) (string, error)
		EcDeleteSecurityGroup(regionId, securityGroupId string) error
		EcListSwitches(regionId, zoneId string) ([]string, error)
		EcCreateVps(regionId string, tmpl VpsCreationTmpl) (string, error)
		EcStartVps(regionId, vpsId string) error
		EcDeleteVps(regionId string, vpsIds []string, force bool) error
		ObsIsBucketExist(regionId, bucketName string) (bool, error)
		ObsNewBucket(regionId, bucketName string) error
		ObsDeleteBucket(regionId, bucketName string, deleteIfNotEmpty *bool) error
		ObsListBuckets(regionId string) ([]string, error)
		ObsListObjectKeys(regionId, bucketName string) ([]string, error)
		ObsScanObjectKeys(regionId, bucketName, pageToken string) ([]string, string, error)
		ObsGetObjectSize(regionId, bucketName, objectKey string) (*gvolume.Volume, error)
		ObsGetObject(regionId, bucketName, objectKey string) ([]byte, error)
		ObsUpsertObject(regionId, bucketName, objectKey string, objectVal []byte) error
		ObsRenameObject(regionId, bucketName, oldObjectKey, newObjectKey string) error
		ObsDeleteObject(regionId, bucketName, objectKey string) error
		SmsSendTmpl(regionId string, mobiles []string, sign, tmpl string, params map[string]string) error
		SmsSendMsg(fromMobile, toMobile, message string) error
		Close() error
	}

	Auth struct {
		AccessKey string
		SecretKey string
	}

	// ACLType bucket/object ACL
	ACLType string
)

var (
	VPS = Production("vps")
	OBS = Production("obs")
	FC  = Production("fc")

	allPlatformInfos = map[Platform]PlatformInfo{}

	Aliyun = enrollPlatform("aliyun", PlatformInfo{SupportPrepaid: true, SupportPostpaid: false})
	AWS    = enrollPlatform("aws", PlatformInfo{SupportPrepaid: false, SupportPostpaid: true})
	GCP    = enrollPlatform("gcp", PlatformInfo{SupportPrepaid: false, SupportPostpaid: true})
	Twilio = enrollPlatform("twilio", PlatformInfo{SupportPrepaid: true, SupportPostpaid: false})
)

func enrollPlatform(name string, info PlatformInfo) Platform {
	allPlatformInfos[Platform(name)] = info
	return Platform(name)
}

func (p *Platform) Info() PlatformInfo {
	info, ok := allPlatformInfos[*p]
	if !ok {
		return PlatformInfo{}
	}
	return info
}

func (p Platform) String() string {
	return string(p)
}

func New(platform Platform, accessKey string, secretKey string, LAN bool) (Cloud, error) {
	switch platform {
	case Aliyun:
		return newAliyun(accessKey, secretKey, LAN)
	case Twilio:
		return newTwilio(accessKey, secretKey)
	default:
		return nil, gerrors.New("unsupported platform %s", platform)
	}
}

func GetAllVps(cloud Cloud) ([]VpsInfo, error) {
	allRegions, err := cloud.ListRegions()
	if err != nil {
		return nil, err
	}
	var result []VpsInfo
	for _, v := range allRegions {
		items, err := cloud.EcListVps(v)
		if err != nil {
			return nil, err
		}
		result = append(result, items...)
	}
	return result, nil
}

type (
	CheapestSpotVpsScanner struct {
		cli    Cloud
		tmpRes map[string]*VpsSpecExList
	}
)

func NewCheapestSpotVpsScanner(cli Cloud) *CheapestSpotVpsScanner {
	return &CheapestSpotVpsScanner{
		cli:    cli,
		tmpRes: map[string]*VpsSpecExList{},
	}
}

func (s *CheapestSpotVpsScanner) Scan() error {
	allRegions, err := s.cli.ListRegions()
	if err != nil {
		return err
	}
	for _, region := range allRegions {
		if _, exist := s.tmpRes[region]; exist {
			continue
		}
		items, err := s.cli.EcListSpotSpecs(region)
		if err != nil {
			return err
		}
		s.tmpRes[region] = items.ToSpecExList()
	}
	return nil
}

func (s *CheapestSpotVpsScanner) GetCheapestNonCreditSpotVpsSpec() *VpsSpecExList {
	res := &VpsSpecExList{}
	for _, v := range s.tmpRes {
		res = res.Append(v)
	}
	return res.RemoveCreditInstance().Sort()
}
