package cloud

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gtest"
	"github.com/cryptowilliam/goutil/container/gdecimal"
	"github.com/cryptowilliam/goutil/encoding/gjson"
	"github.com/cryptowilliam/goutil/sys/gsysinfo"
	"testing"
)

func TestAliyunClient_ListRegions(t *testing.T) {
	cli, err := newAliyun(gsysinfo.GetEnv("ALI_ACCESS"), gsysinfo.GetEnv("ALI_SECRET"), true)
	gtest.Assert(t, err)
	fmt.Println(cli.ListRegions())
}

func TestAliyunClient_GetBalance(t *testing.T) {
	cli, err := newAliyun(gsysinfo.GetEnv("ALI_ACCESS"), gsysinfo.GetEnv("ALI_SECRET"), true)
	gtest.Assert(t, err)
	fmt.Println(cli.GetBalance())
}

func TestAliyunClient_EcListSpotSpecs(t *testing.T) {
	gsysinfo.SetEnv("ALI_ACCESS", "")
	gsysinfo.SetEnv("ALI_SECRET", "")
	cli, err := newAliyun(gsysinfo.GetEnv("ALI_ACCESS"), gsysinfo.GetEnv("ALI_SECRET"), true)
	gtest.Assert(t, err)
	/*specs, err := cli.EcListSpotSpecs("cn-shenzhen")
	gtest.Assert(t, err)
	sortedSpecs := specs.ToSpecExList().RemoveCreditInstance().Sort()
	fmt.Println(gjson.MarshalStringDefault(sortedSpecs, true))
	*/
	res := NewCheapestSpotVpsScanner(cli)
	for {
		if err := res.Scan(); err != nil {
			fmt.Println(err)
		} else {
			break
		}
	}
	fmt.Println(gjson.MarshalStringDefault(res.GetCheapestNonCreditSpotVpsSpec(), true))
}

func TestAliyunClient_EcListImages(t *testing.T) {
	cli, err := newAliyun(gsysinfo.GetEnv("ALI_ACCESS"), gsysinfo.GetEnv("ALI_SECRET"), true)
	gtest.Assert(t, err)
	imgs, err := cli.EcListImages("cn-hongkong")
	gtest.Assert(t, err)
	fmt.Println(gjson.MarshalStringDefault(imgs, true))
}

func TestAliyunClient_EcListVps(t *testing.T) {
	cli, err := newAliyun(gsysinfo.GetEnv("ALI_ACCESS"), gsysinfo.GetEnv("ALI_SECRET"), true)
	gtest.Assert(t, err)
	ins, err := cli.EcListVps("cn-hongkong")
	gtest.Assert(t, err)
	fmt.Println(gjson.MarshalStringDefault(ins, true))
}

func TestAliyunClient_EcCreateVps(t *testing.T) {
	regionId := "cn-hongkong"
	zoneId := "cn-hongkong-c"
	cli, err := newAliyun(gsysinfo.GetEnv("ALI_ACCESS"), gsysinfo.GetEnv("ALI_SECRET"), true)
	gtest.Assert(t, err)

	switchIds, err := cli.EcListSwitches(regionId, zoneId)
	gtest.Assert(t, err)

	sgs, err := cli.EcListSecurityGroups(regionId)
	gtest.Assert(t, err)

	tmpl := VpsCreationTmpl{
		ZoneId:          zoneId,
		SwitchId:        switchIds[0],
		IsSpot:          true,
		Specs:           "ecs.t5-lc2m1.nano",
		Name:            "test-spot-instance-name",
		ImageId:         "ubuntu_20_04_x64_20G_alibase_20201120.vhd",
		SecurityGroupId: sgs[0].Id,
		KeyPair:         "",
		Password:        "jcnde8r74BVGF",
		SystemDiskGB:    35,
		VpsCharge:       PostPaid,
		InternetCharge:  PayByTraffic,
		BandWidthMbIn:   21,
		BandWidthMbOut:  22,
		SpotStrategy:    SpotAsPriceGo,
		SpotMaxPrice:    gdecimal.NewFromFloat64(0.98),
	}
	ins, err := cli.EcCreateVps(regionId, tmpl)
	gtest.Assert(t, err)
	fmt.Println(gjson.MarshalStringDefault(ins, true))
}

func TestAliyunClient_EcDeleteVps(t *testing.T) {
	cli, err := newAliyun(gsysinfo.GetEnv("ALI_ACCESS"), gsysinfo.GetEnv("ALI_SECRET"), true)
	gtest.Assert(t, err)
	err = cli.EcDeleteVps("cn-hangzhou", []string{"i-j6cfq6ofuvebfuzy8qi1"}, true)
	gtest.Assert(t, err)
}

func TestAliyunClient_OssNewBucket(t *testing.T) {
	cli, err := newAliyun(gsysinfo.GetEnv("ALI_ACCESS"), gsysinfo.GetEnv("ALI_SECRET"), true)
	gtest.Assert(t, err)
	fmt.Println(cli.ObsListBuckets("cn-beijing"))
}

func TestAliyunClient_ObsGetObjectSize(t *testing.T) {
	cli, err := newAliyun(gsysinfo.GetEnv("ALI_ACCESS"), gsysinfo.GetEnv("ALI_SECRET"), false)
	gtest.Assert(t, err)
	fmt.Println(cli.ObsGetObjectSize("cn-hongkong", "infdbchunk", "HTZrA5opN07tUu02XVZnzZ07QPlntIuS"))
}

func TestAliyunClient_SmsSendTmpl(t *testing.T) {
	cli, err := newAliyun(gsysinfo.GetEnv("ALI_ACCESS"), gsysinfo.GetEnv("ALI_SECRET"), false)
	gtest.Assert(t, err)
	fmt.Println(cli.SmsSendTmpl("cn-hongkong", []string{""}, "", "", map[string]string{"": ""}))
}
