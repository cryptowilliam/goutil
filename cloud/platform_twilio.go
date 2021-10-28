package cloud

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gdecimal"
	"github.com/cryptowilliam/goutil/container/gvolume"
	"github.com/cryptowilliam/goutil/encoding/gjson"
	"github.com/cryptowilliam/goutil/i18n/gfiat"
	"github.com/sfreiberg/gotwilio"
	"net/http"
	"net/url"
)

type (
	TwilioClient struct {
		accountSid string
		authToken  string
		tw         *gotwilio.Twilio
	}
)

// New twilio client.
func newTwilio(accountSid, authToken string) (*TwilioClient, error) {
	return &TwilioClient{accountSid: accountSid,
		authToken: authToken,
		tw:        gotwilio.NewTwilioClient(accountSid, authToken),
	}, nil
}

// Reference: https://github.com/e154/smart-home/blob/08bfcca81c321b893cad87842a4ea40713f62e09/system/twilio/twilio.go#L102
// Get balance.
func (t *TwilioClient) GetBalance() (*Balance, error) {
	// Balance ...
	type twBalance struct {
		Currency   string `json:"currency"`
		Balance    string `json:"balance"`
		AccountSid string `json:"account_sid"`
	}

	// Build request.
	uri, err := url.Parse(fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Balance.json", t.accountSid))
	if err != nil {
		return nil, gerrors.New(err.Error())
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", uri.String(), nil)
	if err != nil {
		return nil, gerrors.New(err.Error())
	}
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", t.accountSid, t.authToken))))

	// Send request and decode response.
	resp, err := client.Do(req)
	if err != nil {
		return nil, gerrors.New(err.Error())
	}
	defer resp.Body.Close()
	tb := &twBalance{}
	if err = json.NewDecoder(resp.Body).Decode(tb); err != nil {
		return nil, gerrors.New(err.Error())
	}

	// Build 'Balance'.
	res := &Balance{}
	res.Currency, err = gfiat.ParseFiat(tb.Currency)
	if err != nil {
		return nil, err
	}
	res.Available, err = gdecimal.NewFromString(tb.Balance)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (t *TwilioClient) ListRegions() ([]string, error) {
	return nil, gerrors.ErrNotSupport
}

func (t *TwilioClient) ListZones(regionId string) ([]string, error) {
	return nil, gerrors.ErrNotSupport
}

func (t *TwilioClient) EcListSpotSpecs(regionId string) (*VpsSpecList, error) {
	return nil, gerrors.ErrNotSupport
}

func (t *TwilioClient) EcListImages(regionId string) ([]SysImage, error) {
	return nil, gerrors.ErrNotSupport
}

func (t *TwilioClient) EcListVps(regionId string) ([]VpsInfo, error) {
	return nil, gerrors.ErrNotSupport
}

func (t *TwilioClient) EcListSecurityGroups(regionId string) ([]SecurityGroup, error) {
	return nil, gerrors.ErrNotSupport
}

func (t *TwilioClient) EcCreateSecurityGroup(regionId string, sg SecurityGroup) (string, error) {
	return "", gerrors.ErrNotSupport
}

func (t *TwilioClient) EcDeleteSecurityGroup(regionId, securityGroupId string) error {
	return gerrors.ErrNotSupport
}

func (t *TwilioClient) EcListSwitches(regionId, zoneId string) ([]string, error) {
	return nil, gerrors.ErrNotSupport
}

func (t *TwilioClient) EcCreateVps(regionId string, tmpl VpsCreationTmpl) (string, error) {
	return "", gerrors.ErrNotSupport
}

func (t *TwilioClient) EcStartVps(regionId, vpsId string) error {
	return gerrors.ErrNotSupport
}

func (t *TwilioClient) EcDeleteVps(regionId string, vpsIds []string, force bool) error {
	return gerrors.ErrNotSupport
}

func (t *TwilioClient) ObsIsBucketExist(regionId, bucketName string) (bool, error) {
	return false, gerrors.ErrNotSupport
}

func (t *TwilioClient) ObsNewBucket(regionId, bucketName string) error {
	return gerrors.ErrNotSupport
}

func (t *TwilioClient) ObsDeleteBucket(regionId, bucketName string, deleteIfNotEmpty *bool) error {
	return gerrors.ErrNotSupport
}

func (t *TwilioClient) ObsListBuckets(regionId string) ([]string, error) {
	return nil, gerrors.ErrNotSupport
}

func (t *TwilioClient) ObsListObjectKeys(regionId string, bucketName string) ([]string, error) {
	return nil, gerrors.ErrNotSupport
}

func (t *TwilioClient) ObsScanObjectKeys(regionId, bucketName, pageToken string) ([]string, string, error) {
	return nil, "", gerrors.ErrNotSupport
}

func (t *TwilioClient) ObsGetObjectSize(regionId, bucketName, objectKey string) (*gvolume.Volume, error) {
	return nil, gerrors.ErrNotSupport
}

func (t *TwilioClient) ObsGetObject(regionId, bucketName, objectKey string) ([]byte, error) {
	return nil, gerrors.ErrNotSupport
}

func (t *TwilioClient) ObsUpsertObject(regionId, bucketName, objectKey string, objectVal []byte) error {
	return gerrors.ErrNotSupport
}

func (t *TwilioClient) ObsRenameObject(regionId, bucketName, oldObjectKey, newObjectKey string) error {
	return gerrors.ErrNotSupport
}

func (t *TwilioClient) ObsDeleteObject(regionId, bucketName, objectKey string) error {
	return gerrors.ErrNotSupport
}

func (t *TwilioClient) SmsSendTmpl(regionId string, mobiles []string, sign, tmpl string, params map[string]string) error {
	return gerrors.ErrNotSupport
}

// Send SMS.
func (t *TwilioClient) SmsSendMsg(fromMobile, toMobile, message string) error {
	_, ex, err := t.tw.SendSMS(fromMobile, toMobile, message, "", "")
	if ex != nil && ex.Code != 200 {
		return gerrors.New(gjson.MarshalStringDefault(ex, false))
	}
	return err
}

func (t *TwilioClient) FcListFunctions(regionId, service string) ([]string, error) {
	return nil, gerrors.ErrNotSupport
}

func (t *TwilioClient) Close() error {
	return nil
}