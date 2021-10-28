package gsqldb

import (
	"github.com/cryptowilliam/goutil/basic/gtest"
	"testing"
)

type sheet struct {
	K string `xorm:"varchar(256) pk not null 'k'"`
	V string `xorm:"JSON 'v'"`
}

func dialTestDB() (*SqlDB, error) {
	return Dial("mysql://root:msq%!888@tcp(192.168.9.20:3306)/whale?charset=utf8")
}

func TestSqlDB_Tables(t *testing.T) {
	s, err := dialTestDB()
	gtest.Assert(t, err)
	defer s.Close()

	tables, err := s.Tables()
	gtest.Assert(t, err)

	t.Log(tables)
}

func TestSqlDB_SelectOne(t *testing.T) {
	s, err := dialTestDB()
	gtest.Assert(t, err)
	defer s.Close()

	out := sheet{K: "本季度2"}
	ok, err := s.SelectOne(&out)
	gtest.Assert(t, err)
	t.Log(ok)
	t.Log(out)
}

func TestSqlDB_SelectAll(t *testing.T) {
	s, err := dialTestDB()
	gtest.Assert(t, err)
	defer s.Close()

	out := make([]sheet, 0)
	err = s.SelectAll(&out)
	gtest.Assert(t, err)
	t.Log(out)
}

func TestSqlDB_UpsertOne(t *testing.T) {
	s, err := dialTestDB()
	gtest.Assert(t, err)
	defer s.Close()

	newRecord := sheet{
		K: "本季度22",
		V: `{"name":"tom", "age":2222}`,
	}
	n, err := s.UpsertOne(newRecord, &sheet{K: newRecord.K})
	gtest.Assert(t, err)
	t.Log(n)
}

func TestSqlDB_Exist(t *testing.T) {
	s, err := dialTestDB()
	gtest.Assert(t, err)
	defer s.Close()

	exist, err := s.Exist(&sheet{K: "本季度"})
	gtest.Assert(t, err)
	t.Log(exist)
}

func TestSqlDB_Remove(t *testing.T) {
	s, err := dialTestDB()
	gtest.Assert(t, err)
	defer s.Close()

	n, err := s.Delete(&sheet{K: "本季度2"})
	gtest.Assert(t, err)
	t.Log(n)
}
