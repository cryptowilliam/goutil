package gfirewall

import (
	"fmt"
	"runtime"
)

func IsInstalled() bool {
	return false
}

// 防火墙会影响到其他机器访问本机的telnet、tcping等监听端口的服务，但是本机访问本机的服务时不受影响
func IsRunning() bool {
	// service iptables status // centos7/redhat7必须改用systemctl
	// pfctl
	// netsh firewall
	return false
}

func Start() error {
	return nil
}

func Stop() error {
	return nil
}

func printPortRedirHelp(port int) {
	fmt.Printf("*Helpful hint on how to redirect port %d -> 443*\n", port)
	if runtime.GOOS == "windows" {
		fmt.Println("Windows instuctions")
		fmt.Printf("\tnetsh interface portproxy add v4tov4 connectport=%d listenport=443 connectaddress=127.0.0.1 listenaddress=127.0.0.1\n", port)
	} else if runtime.GOOS == "darwin" {
		fmt.Println("OSX instuctions")
		fmt.Printf("\techo \"rdr pass on lo0 inet proto tcp from any to any port 443 -> 127.0.0.1 port %d\" | sudo pfctl -ef -\n", port)
	} else if runtime.GOOS == "linux" {
		fmt.Println("Linux instuctions")
		fmt.Printf("\tsudo iptables -t nat -A PREROUTING -p tcp --dport 443 -j REDIRECT --to-port %d", port)
	}
}
