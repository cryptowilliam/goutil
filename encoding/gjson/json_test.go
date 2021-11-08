package gjson

import (
	"github.com/cryptowilliam/goutil/basic/gtest"
	"github.com/cryptowilliam/goutil/container/gstring"
	"github.com/tidwall/gjson"
	"testing"
)

const (
	demo_json1 = `{"T":"2018-04-10T16:14:08.364623+08:00","Content":123}`
	demo_json2 = `{"name":{"first":"Janet","last":"Prichard"},"age":47}`
	demoJson3  = `
{
	"DSN.modifyMe": "hello dsn",
	"KlineLocation": "kline",
	"LivingLogLocation": "living-log",
	"TimeZone": "Asia/Shanghai",
    "ServerName": null,
	"StrategyOpts": [
		{
			"StrategyConfig": {
				"IsBackTest": false,
				"Name": "StrategyDemo",
				"Portfolio": {
					"DirectSet": [
						"ETH/USD.spot.coinbase"
					]
				},
				"Period": "1hour",
				"KlineCloseTimeout": "5 minutes",
				"TradeCmdExecTimeout": "5 minutes",
				"Fee": "0",
				"MaxSlippage": "0",
				"MaxDrawDown": "0",
				"IndExpr": {
					"IndExpr1": "MACD(50)",
					"IndExpr2": "MACD(100)"
				}
			},
			"AccInfos": [
				{
					"Email": "demo@example.com",
					"Platform": "coinbase",
					"ApiKey.modifyMe": "api key plain",
					"SecretKey.modifyMe": "secret key plain",
					"Proxy": "proxy address"
				}
			]
		}
	]
}`
)

func TestGet(t *testing.T) {
	jv := Get(demo_json1, "Content")
	if !jv.Exists() {
		t.Error("Get error, not exists")
		return
	}
	if jv.Type != gjson.Number {
		t.Error("Get error, type error")
		return
	}
}

func TestIterate(t *testing.T) {
	iterFn := func(key string, val interface{}) (newVal interface{}, modified bool, err error) {
		if gstring.EndWith(key, ".modifyMe") {
			return "new modified string", true, nil
		}
		return nil, false, nil
	}

	jsonStr := demoJson3
	err := Iterate(&jsonStr, true, iterFn)
	gtest.Assert(t, err)
}
