package gaddr

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/grand"
	"net"
	"strconv"
	"strings"
)

func IsValidPort(port int) bool {
	return port >= 0 && port <= 65535
}

// Privileged port needs root
func IsPrivilegedPort(port int) bool {
	return port >= 0 && port <= 1023
}

func IsLocalPortUsing(nettype string, port int) (bool, error) {
	if nettype == "tcp" {
		ln, err := net.Listen("tcp", "0.0.0.0:"+strconv.FormatInt(int64(port), 10))
		if err != nil {
			return true, nil
		} else {
			_ = ln.Close()
			return false, nil
		}
	} else if nettype == "udp" {
		addr, _ := net.ResolveUDPAddr("udp", "0.0.0.0:"+strconv.FormatInt(int64(port), 10))
		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			return true, nil
		} else {
			_ = conn.Close()
			return false, nil
		}
	} else {
		return false, gerrors.New("unsupported net type " + nettype)
	}
}

/*
// github.com/wheelcomplex/lsof pure go lsof, but support linux only
func whoIsListeningLocalPort(nettype string, port int) (xproc.ProcId, error) {
	if runtime.GOOS == "windows" {
		return xproc.InvalidProcId, gerrors.New("Windows unsupported for now")
	} else {
		// 下面这个命令其实就是lsof -i:****的升级版，而且限定了TCP，其实lsof还可以检测UDP的
		cmd := fmt.Sprintf("lsof -nP -iTCP:%v -sTCP:LISTEN | awk -v i=2 -v j=2 'FNR == i {print $j}'", port)
		result := xcmd.ExecWait(cmd, false)
		result = strings.Trim(result, "\n")
		if result == "" {
			return xproc.InvalidProcId, gerrors.New("port is not using")
		}
		pid, err := strconv.ParseInt(result, 10, 32)
		if err != nil {
			return xproc.InvalidProcId, err
		}
		return xproc.ProcId(pid), nil
	}
}
*/

// check exchange port, 但有的没lsof，有的报错
/*
	ERRO[2017-06-10 16:47:37.665] strconv.ParseInt: parsing "lsof: status error on |: No such file or directory\nlsof: status error on awk: No such file or directory\nlsof: status error on -v: No such file or directory\nlsof: status error on i=2: No such file or directory\nlsof: status error on -v: No such file or directory\nlsof: status error on j=2: No such file or directory\nlsof: status error on 'FNR: No such file or directory\nlsof: status error on ==: No such file or directory\nlsof: status error on i: No such file or directory\nlsof: status error on {print: No such file or directory\nlsof: status error on $j}': No such file or directory\nlsof 4.89\n latest revision: ftp://lsof.itap.purdue.edu/pub/tools/unix/lsof/\n latest FAQ: ftp://lsof.itap.purdue.edu/pub/tools/unix/lsof/FAQ\n latest man page: ftp://lsof.itap.purdue.edu/pub/tools/unix/lsof/lsof_man\n usage: [-?abhKlnNoOPRtUvVX] [+|-c c] [+|-d s] [+D D] [+|-E] [+|-e s] [+|-f[gG]]\n [-F [f]] [-g [s]] [-i [i]] [+|-L [l]] [+m [m]] [+|-M] [-o [o]] [-p s]\n [+|-r [t]] [-s [p:s]] [-S [t]] [-T [t]] [-u s] [+|-w] [-x [fl]] [--] [names]\nUse the ``-h'' option to get more help information.": invalid syntax
*/
// Check process id by using port doesn't works well on linux, also not support other OS
/*
func CheckLocalPort(nettype string, port int) (isUsing bool, usingByWho xproc.ProcId, err error) {
	nettype = strings.ToLower(nettype)
	if !IsValidPort(port) {
		return false, xproc.InvalidProcId, gerrors.New("invalid input port " + strconv.FormatInt(int64(port), 10))
	}
	if nettype != "tcp" && nettype != "udp" {
		return false, xproc.InvalidProcId, gerrors.New("unsupported network type " + nettype)
	}

	// Is port using
	isPortUsing, err := IsLocalPortUsing(nettype, port)
	if err != nil {
		return false, xproc.InvalidProcId, err
	}

	// who is using port
	if isPortUsing {
		listeningPortPid, err := whoIsListeningLocalPort(nettype, port)
		if err != nil {
			return false, xproc.InvalidProcId, err
		}
		return isPortUsing, listeningPortPid, nil
	}

	return isPortUsing, xproc.InvalidProcId, nil
}*/

// 端口是否通畅
// 如果发起检测端不是你的程序，需要开关端口监听多次才能确定，因为，如果端口可访问，也可能是被映射到别的电脑上而那个电脑的该端口也被监听了
func IsLocalInboundPortClear(nettype string, port int) (bool, error) {
	return false, nil
}

// 只能检测端口是否正开启，如果未开启，可能是多种原因，比如主机关机、防火墙、端口未映射、端口没有服务在监听等
func IsRemotePortOpen(nettype string, port int) (bool, error) {
	return false, nil
}

func GetRandomAvailablePort(nettype string) (int, error) {
	nettype = strings.ToLower(nettype)
	if nettype != "tcp" && nettype != "udp" {
		return 0, gerrors.New("unsupported network type " + nettype)
	}
	for {
		port := grand.RandomInt(1024, 65535)
		isUsing, err := IsLocalPortUsing(nettype, port)
		if err != nil || isUsing {
			continue
		}
		return port, nil
	}
}

//portquiz.net
// 所谓检测outgoing port，不太理解测的是个啥
