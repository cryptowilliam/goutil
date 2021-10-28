package cloud

import (
	"github.com/cryptowilliam/goutil/container/gvolume"
	"time"
)

type (
	Retry struct {
		cloud         Cloud
		retry         int
		retryInterval time.Duration
	}

	RetryCloud interface {
		Cloud
		SetRetry(retry uint8)
		SetRetryInterval(dur time.Duration)
	}
)

func NewRetry(platform Platform, accessKey string, secretKey string, LAN bool) (RetryCloud, error) {
	c, err := New(platform, accessKey, secretKey, LAN)
	if err != nil {
		return nil, err
	}
	res := Retry{
		cloud:         c,
		retry:         1,
		retryInterval: time.Duration(0),
	}
	return &res, nil
}

func (c *Retry) SetRetry(retry uint8) {
	c.retry = int(retry)
}

func (c *Retry) SetRetryInterval(dur time.Duration) {
	c.retryInterval = dur
}

func (c *Retry) GetBalance() (balance *Balance, err error) {
	for i := 0; i < c.retry; i++ {
		balance, err = c.cloud.GetBalance()
		if err == nil {
			break
		}
		time.Sleep(c.retryInterval)
	}
	return balance, err
}

func (c *Retry) ListRegions() (names []string, err error) {
	for i := 0; i < c.retry; i++ {
		names, err = c.cloud.ListRegions()
		if err == nil {
			break
		}
		time.Sleep(c.retryInterval)
	}
	return names, err
}

func (c *Retry) ListZones(regionId string) (names []string, err error) {
	for i := 0; i < c.retry; i++ {
		names, err = c.cloud.ListZones(regionId)
		if err == nil {
			break
		}
		time.Sleep(c.retryInterval)
	}
	return names, err
}

func (c *Retry) EcListSpotSpecs(regionId string) (res *VpsSpecList, err error) {
	for i := 0; i < c.retry; i++ {
		res, err = c.cloud.EcListSpotSpecs(regionId)
		if err == nil {
			break
		}
		time.Sleep(c.retryInterval)
	}
	return res, err
}

func (c *Retry) EcListImages(regionId string) (res []SysImage, err error) {
	for i := 0; i < c.retry; i++ {
		res, err = c.cloud.EcListImages(regionId)
		if err == nil {
			break
		}
		time.Sleep(c.retryInterval)
	}
	return res, err
}

func (c *Retry) EcListVps(regionId string) (res []VpsInfo, err error) {
	for i := 0; i < c.retry; i++ {
		res, err = c.cloud.EcListVps(regionId)
		if err == nil {
			break
		}
		time.Sleep(c.retryInterval)
	}
	return res, err
}

func (c *Retry) EcListSecurityGroups(regionId string) (res []SecurityGroup, err error) {
	for i := 0; i < c.retry; i++ {
		res, err = c.cloud.EcListSecurityGroups(regionId)
		if err == nil {
			break
		}
		time.Sleep(c.retryInterval)
	}
	return res, err
}

func (c *Retry) EcCreateSecurityGroup(regionId string, sg SecurityGroup) (securityGroupId string, err error) {
	for i := 0; i < c.retry; i++ {
		securityGroupId, err = c.cloud.EcCreateSecurityGroup(regionId, sg)
		if err == nil {
			break
		}
		time.Sleep(c.retryInterval)
	}
	return securityGroupId, err
}

func (c *Retry) EcDeleteSecurityGroup(regionId, securityGroupId string) (err error) {
	for i := 0; i < c.retry; i++ {
		err = c.cloud.EcDeleteSecurityGroup(regionId, securityGroupId)
		if err == nil {
			break
		}
		time.Sleep(c.retryInterval)
	}
	return err
}

func (c *Retry) EcListSwitches(regionId, zoneId string) (res []string, err error) {
	for i := 0; i < c.retry; i++ {
		res, err = c.cloud.EcListSwitches(regionId, zoneId)
		if err == nil {
			break
		}
		time.Sleep(c.retryInterval)
	}
	return res, err
}

func (c *Retry) EcCreateVps(regionId string, tmpl VpsCreationTmpl) (res string, err error) {
	for i := 0; i < c.retry; i++ {
		res, err = c.cloud.EcCreateVps(regionId, tmpl)
		if err == nil {
			break
		}
		time.Sleep(c.retryInterval)
	}
	return res, err
}

func (c *Retry) EcStartVps(regionId, vpsId string) (err error) {
	for i := 0; i < c.retry; i++ {
		err = c.cloud.EcStartVps(regionId, vpsId)
		if err == nil {
			break
		}
		time.Sleep(c.retryInterval)
	}
	return err
}

func (c *Retry) EcDeleteVps(regionId string, vpsIds []string, force bool) (err error) {
	for i := 0; i < c.retry; i++ {
		err = c.cloud.EcDeleteVps(regionId, vpsIds, force)
		if err == nil {
			break
		}
		time.Sleep(c.retryInterval)
	}
	return err
}


func (c *Retry) ObsIsBucketExist(regionId, bucketName string) (exist bool, err error) {
	for i := 0; i < c.retry; i++ {
		exist, err = c.cloud.ObsIsBucketExist(regionId, bucketName)
		if err == nil {
			break
		}
		time.Sleep(c.retryInterval)
	}
	return exist, err
}

func (c *Retry) ObsNewBucket(regionId, bucketName string) (err error) {
	for i := 0; i < c.retry; i++ {
		err = c.cloud.ObsNewBucket(regionId, bucketName)
		if err == nil {
			break
		}
		time.Sleep(c.retryInterval)
	}
	return err
}

func (c *Retry) ObsDeleteBucket(regionId, bucketName string, deleteIfNotEmpty *bool) (err error) {
	for i := 0; i < c.retry; i++ {
		err = c.cloud.ObsDeleteBucket(regionId, bucketName, deleteIfNotEmpty)
		if err == nil {
			break
		}
		time.Sleep(c.retryInterval)
	}
	return err
}

func (c *Retry) ObsListBuckets(regionId string) (names []string, err error) {
	for i := 0; i < c.retry; i++ {
		names, err = c.cloud.ObsListBuckets(regionId)
		if err == nil {
			break
		}
		time.Sleep(c.retryInterval)
	}
	return names, err
}

func (c *Retry) ObsListObjectKeys(regionId, bucketName string) (keys []string, err error) {
	for i := 0; i < c.retry; i++ {
		keys, err = c.cloud.ObsListObjectKeys(regionId, bucketName)
		if err == nil {
			break
		}
		time.Sleep(c.retryInterval)
	}
	return keys, err
}

func (c *Retry) ObsScanObjectKeys(regionId, bucketName, pageToken string) (keys []string, nextPageToken string, err error) {
	for i := 0; i < c.retry; i++ {
		keys, nextPageToken, err = c.cloud.ObsScanObjectKeys(regionId, bucketName, pageToken)
		if err == nil {
			break
		}
		time.Sleep(c.retryInterval)
	}
	return keys, nextPageToken, err
}

func (c *Retry) ObsGetObjectSize(regionId, bucketName, objectKey string) (vol *gvolume.Volume, err error) {
	for i := 0; i < c.retry; i++ {
		vol, err = c.cloud.ObsGetObjectSize(regionId, bucketName, objectKey)
		if err == nil {
			break
		}
		time.Sleep(c.retryInterval)
	}
	return vol, err
}

func (c *Retry) ObsGetObject(regionId, bucketName, blobId string) (val []byte, err error) {
	for i := 0; i < c.retry; i++ {
		val, err = c.cloud.ObsGetObject(regionId, bucketName, blobId)
		if err == nil {
			break
		}
		time.Sleep(c.retryInterval)
	}
	return val, err
}

func (c *Retry) ObsUpsertObject(regionId, bucketName, blobId string, blob []byte) (err error) {
	for i := 0; i < c.retry; i++ {
		err = c.cloud.ObsUpsertObject(regionId, bucketName, blobId, blob)
		if err == nil {
			break
		}
		time.Sleep(c.retryInterval)
	}
	return err
}

func (c *Retry) ObsRenameObject(regionId, bucketName, oldKey, newKey string) (err error) {
	for i := 0; i < c.retry; i++ {
		err = c.cloud.ObsRenameObject(regionId, bucketName, oldKey, newKey)
		if err == nil {
			break
		}
		time.Sleep(c.retryInterval)
	}
	return err
}

func (c *Retry) ObsDeleteObject(regionId, bucketName, blobId string) (err error) {
	for i := 0; i < c.retry; i++ {
		err = c.cloud.ObsDeleteObject(regionId, bucketName, blobId)
		if err == nil {
			break
		}
		time.Sleep(c.retryInterval)
	}
	return err
}

func (c *Retry) SmsSendTmpl(regionId string, mobiles []string, sign, tmpl string, params map[string]string) (err error) {
	for i := 0; i < c.retry; i++ {
		err = c.cloud.SmsSendTmpl(regionId, mobiles, sign, tmpl, params)
		if err == nil {
			break
		}
		time.Sleep(c.retryInterval)
	}
	return err
}

func (c *Retry) SmsSendMsg(fromMobile, toMobile, message string) (err error) {
	for i := 0; i < c.retry; i++ {
		err = c.cloud.SmsSendMsg(fromMobile, toMobile, message)
		if err == nil {
			break
		}
		time.Sleep(c.retryInterval)
	}
	return err
}

func (c *Retry) Close() (err error) {
	return c.Close()
}
