package config

import (
	"flag"

	"github.com/hsfzxjy/pipe"
	"github.com/hsfzxjy/zunagi/internal/role"
)

var HostAddress = pipe.NewControllerCM("", true)

func Setup(r role.Role) {
	var defaultHostAddress string
	switch r {
	case role.Host:
		defaultHostAddress = "127.0.0.1:2783"
	case role.Guest:
		defaultHostAddress = "10.0.2.2:2783"
	}
	host_address := flag.String("host_address", defaultHostAddress, "")
	flag.Parse()
	HostAddress.Send(*host_address)
}
