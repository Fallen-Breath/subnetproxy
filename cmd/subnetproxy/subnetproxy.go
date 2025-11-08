package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/Fallen-Breath/subnetproxy/internal/server"
	"github.com/Fallen-Breath/subnetproxy/internal/utils"
)

const Version = "0.1.0"

func main() {
	listenAddr := flag.String("listen", ":1080", "Address for the socks5 server to listen on")
	subnetStr := flag.String("subnet", "", "Comma-separated subnets for IP pool (e.g., 192.168.1.0/24,10.0.0.0/8)")
	enableProxyProtocol := flag.Bool("proxyprotocol", false, "Enable PROXY protocol support to get correct client ip")
	strategy := flag.String("strategy", "hash", "Strategy for selecting local IP: hash (hash client ip) or random")
	showVersion := flag.Bool("version", false, "Show version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("subnetproxy v%s\n", Version)
		return
	}

	log.Printf("CONFIG: listen = %s", *listenAddr)
	log.Printf("CONFIG: subnet = %s", *subnetStr)
	log.Printf("CONFIG: proxyprotocol = %t", *enableProxyProtocol)
	log.Printf("CONFIG: strategy = %s", *strategy)

	var subnet *utils.IpPool
	if *subnetStr == "" {
		log.Printf("WARN: subnet is not provided, the default outbound address will be used")
	} else {
		var err error
		subnet, err = utils.NewIpPool(strings.Split(*subnetStr, ","))
		if err != nil {
			log.Fatalf("Failed to create IP pool: %v", err)
		}
	}
	if *strategy != "hash" && *strategy != "random" {
		log.Fatal("Invalid strategy: must be 'hash' or 'random'")
	}

	svr := server.Server{
		Listen:        *listenAddr,
		Subnet:        subnet,
		ProxyProtocol: *enableProxyProtocol,
		Strategy:      *strategy,
	}
	svr.Serve()
}
