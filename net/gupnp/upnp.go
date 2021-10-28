package gupnp

import (
	"fmt"
	"github.com/prestonTao/upnp"
)

/*
https://github.com/prestonTao/upnp
https://github.com/huin/goupnp
https://github.com/ianr0bkny/go-sonos
https://github.com/jackpal/go-nat-pmp
https://github.com/hlandau/portmap
nat-pmp和upnp都要支持，目前看上述库都没有卵用
*/

func PortMap() bool {

	mapping := new(upnp.Upnp)
	if err := mapping.AddPortMapping(2048, 2048, "TCP"); err == nil {
		fmt.Println("success !")
		// remove port mapping in gatway
		mapping.Reclaim()
	} else {
		fmt.Println("fail !")
	}
	return false
}
