package gdsn

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gstring"
	. "github.com/cryptowilliam/goutil/database/gdriver"
	"github.com/cryptowilliam/goutil/net/gnet"
	"github.com/jackc/pgx"
	"net"
	"strconv"
	"strings"
	"upper.io/db.v3/mongo"
	"upper.io/db.v3/mysql"
	"upper.io/db.v3/sqlite"
)

type DSN struct {
	Driver   Driver
	User     string
	Password string
	Host     string // Domain or IP address without port (e.g. localhost) or path to unix domain socket directory (e.g. /private/tmp)
	Port     uint16
	Path     string // Sub path after port, but before Options. It is database sometimes.
	Database string
	Options  map[string]string
}

// Compare to gaddr.ParseURL, gdsn.Parse use professional database package to parse data source name string.
// Some databases have DSNs in specific formats, for example, MongoDB DSN supports multiple ports.
func Parse(s string) (*DSN, error) {
	res := DSN{}
	res.Options = make(map[string]string)

	if gstring.StartWith(strings.ToLower(s), "dragondb://") {
		res.Driver = DragonDB
	} else if gstring.StartWith(strings.ToLower(s), "mongodb://") {
		res.Driver = MongoDB
	} else if gstring.StartWith(strings.ToLower(s), "redis://") {
		res.Driver = Redis
	} else if gstring.StartWith(strings.ToLower(s), "sqlite://") {
		res.Driver = SQLite
	} else if gstring.StartWith(strings.ToLower(s), "mysql://") {
		res.Driver = MySQL
	} else if gstring.StartWith(strings.ToLower(s), "postgres://") {
		res.Driver = PgSQL
	} else if gstring.StartWith(strings.ToLower(s), "tidb://") {
		res.Driver = TiDB
	} else if gstring.StartWith(strings.ToLower(s), "cockroach://") {
		res.Driver = CockroachDB
	} else if gstring.StartWith(strings.ToLower(s), "mssql://") {
		res.Driver = Mssql
	} else if gstring.StartWith(strings.ToLower(s), "oracle://") {
		res.Driver = Oracle
	} else {
		return nil, gerrors.New("unrecognized dirver for dsn %s", s)
	}

	switch res.Driver {
	case DragonDB:
		as, err := gnet.ParseUrl(s)
		if err != nil {
			return nil, err
		}
		res.Host = as.Host.Domain + as.Host.IP
		res.Port = uint16(as.Host.Port)
		res.User = as.Auth.User
		res.Password = as.Auth.Password
		res.Path = as.Path.Str
		return &res, nil

	case MongoDB:
		mgConf, err := mongo.ParseURL(s)
		if err != nil {
			return nil, err
		}
		host, portStr, err := net.SplitHostPort(mgConf.Host)
		if err != nil {
			return nil, err
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, err
		}
		res.Host = host
		res.Port = uint16(port)
		res.User = mgConf.User
		res.Password = mgConf.Password
		res.Database = mgConf.Database
		res.Options = mgConf.Options
		return &res, nil

	case SQLite:
		slConf, err := sqlite.ParseURL(s)
		if err != nil {
			return nil, err
		}
		res.Host = ""
		res.Port = 0
		res.User = ""
		res.Password = ""
		res.Database = slConf.Database
		res.Options = slConf.Options
		return &res, nil

	case MySQL:
		myConf, err := mysql.ParseURL(s)
		if err != nil {
			return nil, err
		}
		host, portStr, err := net.SplitHostPort(myConf.Host)
		if err != nil {
			return nil, err
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, err
		}
		res.Host = host
		res.Port = uint16(port)
		res.Database = myConf.Database
		res.User = myConf.User
		res.Password = myConf.Password
		res.Options = myConf.Options
		return &res, nil

	case PgSQL:
		pgConf, err := pgx.ParseConnectionString(s)
		if err != nil {
			return nil, err
		}
		res.Host = pgConf.Host
		res.Port = pgConf.Port
		res.User = pgConf.User
		res.Password = pgConf.Password
		res.Options = nil
		return &res, nil

	case TiDB:
		s = gstring.HeadToLowerASCII(s, len(res.Driver.String()))
		s = strings.Replace(s, "tidb://", "mysql://", 1)
		res, err := Parse(s)
		if err != nil {
			return nil, err
		}
		res.Driver = TiDB
		return res, nil

	default:
		return nil, gerrors.Wrap(ErrUnsupportedDriver, res.Driver.String())
	}
}

func ParseForSpecificDriver(s string, dvr Driver) (*DSN, error) {
	res, err := Parse(s)
	if err != nil {
		return nil, err
	}
	if res.Driver != dvr {
		return nil, gerrors.New("Driver %s required but %s got", dvr.String(), res.Driver.String())
	}
	return res, nil
}

func (ci *DSN) HostAndPort() string {
	return ci.Host + ":" + strconv.FormatInt(int64(ci.Port), 10)
}

func (ci *DSN) Build(dvr Driver) (string, error) {
	switch dvr {
	case MySQL:
		optStr := ""
		for k, v := range ci.Options {
			optStr += k + "=" + v + "&"
		}
		optStr = gstring.RemoveTail(optStr, 1)
		return fmt.Sprintf("%s:%s@%s/%s?%s", ci.User, ci.Password, ci.Host, ci.Database, optStr), nil
	default:
		return "", gerrors.New("unsupported driver %s", dvr.String())
	}
}
