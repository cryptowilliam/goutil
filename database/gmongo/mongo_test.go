package gmongo

import (
	"fmt"
	"github.com/cryptowilliam/goutil/sys/gtime"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"time"
	"testing"
)

/*
func TestNewConn(t *testing.T) {
	conn, err := Dial("mongodb://192.168.9.11:27717")
	if err != nil {
		t.Error(err)
		return
	}

	now := time.Now()
	kl := finance.K{}
	kl.T = now
	if err != nil {
		t.Error(err)
		return
	}
	kl.H = 500
	kl.L = 1.00000123
	err = conn.Database("coins").Collection("kline").Insert(kl)
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(time.Second)
	kl.L = 1.000001234567
	err = conn.Database("coins").Collection("kline").Upsert(bson.M{"T":kl.T}, bson.M{"$set": kl})
	if err != nil {
		t.Error(err)
		return
	}

	kl.T = kl.T.Add(time.Second)
	err = conn.Database("coins").Collection("kline").Upsert(bson.M{"T":kl.T}, bson.M{"$set": kl})
	if err != nil {
		t.Error(err)
		return
	}

	maxTime, _, err := conn.Database("coins").Collection("kline").MaxTime("T")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(maxTime)
}
*/
func TestColl_Find(t *testing.T) {
	conn, err := Dial("mongodb://127.0.0.1:27017")
	if err != nil {
		t.Error(err)
		return
	}

	sz, err := conn.Database("Details2").Collection("123").Count(nil)
	fmt.Println(sz, err)

	cur, err := conn.Database("Details").Collection("123").Find(nil)
	fmt.Println(cur.Next(), err)
}

func TestColl_UpsertEntireDoc_UpsertFields(t *testing.T) {
	conn, err := Dial("mongodb://127.0.0.1:27017")
	if err != nil {
		t.Error(err)
		return
	}

	type TestItemA struct {
		Name  string `bson:"_id"`
		Age   int    `bson:"Age"`
		Class string `bson:"Class"`
	}

	type TestItemB struct {
		Name string `bson:"_id"`
		Age  int    `bson:"Age"`
		Step string `bson:"Step"`
	}

	type TestItemC struct {
		Age  int    `bson:"Age"`
		Step string `bson:"Step"`
	}

	itemA := TestItemA{Name: "WangJiaChen", Age: 18, Class: "#1"}
	itemRetA := TestItemA{}
	if err = conn.Database("test").Collection("test").UpsertEntireDoc(itemA.Name, itemA); err != nil {
		t.Error(err)
		return
	}
	ok, err := conn.Database("test").Collection("test").FindId(itemA.Name, &itemRetA)
	if err != nil {
		t.Error(err)
		return
	}
	if !ok {
		t.Errorf("must FindId")
		return
	}
	fmt.Println(itemRetA)
	if itemRetA != itemA {
		t.Error("UpsertEntireDoc error")
		return
	}

	itemB := TestItemB{Name: "WangJiaChen", Age: 19, Step: "#2"}
	itemRetB := TestItemB{}
	if err = conn.Database("test").Collection("test").UpsertEntireDoc(itemB.Name, itemB); err != nil {
		t.Error(err)
		return
	}
	ok, err = conn.Database("test").Collection("test").FindId(itemB.Name, &itemRetB)
	if err != nil {
		t.Error(err)
		return
	}
	if !ok {
		t.Errorf("must FindId")
		return
	}
	fmt.Println(itemRetB)
	if itemRetB != itemB {
		t.Error("UpsertEntireDoc error")
		return
	}

	itemC := TestItemC{Age: 20, Step: "#3"}
	itemRetC := TestItemB{}
	if err = conn.Database("test").Collection("test").UpsertFields(itemRetB.Name, itemC); err != nil {
		t.Error(err)
		return
	}
	ok, err = conn.Database("test").Collection("test").FindId(itemB.Name, &itemRetC)
	if err != nil {
		t.Error(err)
		return
	}
	if !ok {
		t.Errorf("must FindId")
		return
	}
	fmt.Println(itemRetC)
	if itemRetC.Name != "WangJiaChen" || itemRetC.Age != 20 || itemRetC.Step != "#3" {
		t.Error("UpsertEntireDoc error")
		return
	}
}

type tmpKline struct {
	Time      time.Time `json:"T" bson:"_id"`
	Open      float64   `json:"O" bson:"O"` // open price in USD
	Close     float64   `json:"C" bson:"C"`
	High      float64   `json:"H" bson:"H"`
	Low       float64   `json:"L" bson:"L"`
	Volume    float64   `json:"V" bson:"V"`                                     // volume in USD
	MarketCap float64   `json:"MarketCap,omitempty" bson:"MarketCap,omitempty"` // exchangeName cap in USD CoinMarketCap接口中用到
}

func TestColl_FindCmp(t *testing.T) {
	conn, err := Dial("mongodb://127.0.0.1:27017")
	if err != nil {
		t.Error(err)
		return
	}

	date, err := gtime.NewDate(2018, 11, 1)
	if err != nil {
		t.Error(err)
		return
	}

	cur, err := conn.Database("TradeHistory-YF").Collection("s.nasdaq.aapl").FindCmp("_id", CmpGTE, bsonx.Time(date.ToTimeUTC()))
	item := tmpKline{}
	total := 0
	for cur.Next() {
		if err := cur.Decode(&item); err != nil {
			t.Error(err)
			return
		} else {
			t.Log(item)
			total++
		}
	}
	t.Log(total)
}

func TestColl_RemoveCmp(t *testing.T) {
	conn, err := Dial("mongodb://127.0.0.1:27017")
	if err != nil {
		t.Error(err)
		return
	}

	date, err := gtime.NewDate(2019, 1, 1)
	if err != nil {
		t.Error(err)
		return
	}

	ds, err := conn.Database("TradeHistory-CMC").Collection("c.0chain").RemoveCmp("_id", CmpGTE, bsonx.Time(date.ToTimeUTC()))
	t.Log(ds, err)
	t.Log(conn.Database("TradeHistory-CMC").Collection("c.0chain").MaxTime("_id"))
}
