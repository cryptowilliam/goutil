package gnet

import (
	"fmt"
	"net"
	"strings"
)

// TODO: New Url parse module

//scheme://user:pass@host(domain/ip):port/dir1/.../dirN?params
//mailto:user@domain
//mysql://user:pass@host(domain/ip):port/database?options

type (
	Auth struct {
		User string
		Pass string
	}

	Host struct {
		Domain string
		IP     net.IP
	}

	Url struct {
		Proto   string
		Auth    Auth
		Host    Host
		Port    int
		Paths   []string
		Queries map[string]string
	}
)

func (h *Host) String() string {
	if h.Domain != "" {
		return h.Domain
	}
	if h.IP != nil {
		return h.IP.String()
	}
	return ""
}

func NewUrl() *Url {
	return &Url{}
}

func (u *Url) SetProto(proto string) *Url {
	u.Proto = proto
	return u
}

func (u *Url) SetUser(user string) *Url {
	u.Auth.User = user
	return u
}

func (u *Url) SetPass(pass string) *Url {
	u.Auth.Pass = pass
	return u
}

func (u *Url) SetHost(host string) *Url {
	if ip := net.ParseIP(host); ip != nil {
		u.Host.Domain = ""
		u.Host.IP = ip
	} else {
		u.Host.Domain = host
		u.Host.IP = nil
	}
	return u
}

func (u *Url) SetPort(port int) *Url {
	u.Port = port
	return u
}

func (u *Url) SetPath(paths []string) *Url {
	u.Paths = paths
	return u
}

func (u *Url) SetQueries(queries map[string]string) *Url {
	u.Queries = queries
	return u
}

func (u *Url) String() string {
	if strings.ToLower(u.Proto) == "mailto" {
		return fmt.Sprintf("%s://%s@%s", strings.ToLower(u.Proto), u.Auth.User, u.Host.String())
	} else {
		//TODO
		return ""
	}
}
