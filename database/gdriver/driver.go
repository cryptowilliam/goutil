package gdriver

import "github.com/cryptowilliam/goutil/basic/gerrors"

var (
	ErrUnsupportedDriver = gerrors.New("Unsupported Database Driver")
)

type Driver string

// Don't support old databases such like MSSQL, Oracle, Access...
// Old databases will break compatibility.
const (
	DragonDB    = Driver("dragondb")
	MongoDB     = Driver("mongodb")
	Redis       = Driver("redis")
	SQLite      = Driver("sqlite")
	MySQL       = Driver("mysql")
	PgSQL       = Driver("pgsql")
	TiDB        = Driver("tidb")
	CockroachDB = Driver("cockroachdb")
	Mssql       = Driver("mssql")
	Oracle      = Driver("oracle")
)

func (dvr Driver) String() string {
	return string(dvr)
}
