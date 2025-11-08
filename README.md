# subnet proxy

A simple socks5 proxy server that utilize the given subnet to create outbound requests

```bash
$ ./subnetproxy
Usage of ./subnetproxy:
  -listen string
        Address for the socks5 server to listen on (default ":1080")
  -proxyprotocol
        Enable PROXY protocol support to get correct client ip
  -strategy string
        Strategy for selecting local IP: hash (hash client ip) or random (default "hash")
  -subnet string
        Comma-separated subnets for IP pool (e.g., 192.168.1.0/24,10.0.0.0/8)
  -version
        Show version and exit
```

Docker image is available at [fallenbreath/subnetproxy](https://hub.docker.com/r/fallenbreath/subnetproxy)
