package server

import (
	"context"
	"log"
	"net"

	"github.com/Fallen-Breath/subnetproxy/internal/utils"
	"github.com/armon/go-socks5"
	"github.com/pires/go-proxyproto"
)

type Server struct {
	Listen        string
	Subnet        *utils.IpPool
	ProxyProtocol bool
	Strategy      string
}

func (s *Server) Serve() {
	ln, err := net.Listen("tcp", s.Listen)
	if err != nil {
		log.Fatal(err)
	}
	if s.ProxyProtocol {
		ln = &proxyproto.Listener{Listener: ln}
		log.Println("Proxy protocol support is enabled")
	}
	log.Printf("Starting socks5 server on %s", s.Listen)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}

		go func() {
			defer func() {
				_ = conn.Close()
			}()
			s.handleConnection(conn)
		}()
	}
}

type sock5Dialer = func(ctx context.Context, network, addr string) (net.Conn, error)

type sock5Resolver struct {
	netType string
}

func (r *sock5Resolver) Resolve(ctx context.Context, name string) (context.Context, net.IP, error) {
	addr, err := net.ResolveIPAddr(r.netType, name)
	if err != nil {
		return ctx, nil, err
	}
	return ctx, addr.IP, err
}

func (s *Server) handleConnection(conn net.Conn) {
	clientAddr, ok := conn.RemoteAddr().(*net.TCPAddr)
	if !ok {
		log.Printf("Non-TCP connection from %s, skipping", conn.RemoteAddr().String())
		return
	}
	clientIP := clientAddr.IP.String()

	var dialer sock5Dialer
	var resolver socks5.NameResolver
	if s.Subnet != nil {
		var localIP net.IP
		if s.Strategy == "hash" {
			localIP = s.Subnet.GetByKey(clientIP)
		} else if s.Strategy == "random" {
			localIP = s.Subnet.GetRandomly()
		} else {
			panic("unknown strategy " + s.Strategy)
		}
		log.Printf("%s --(%s)-> outbound", clientIP, localIP)

		var network, netType string
		if localIP.To4() == nil {
			network = "tcp6"
			netType = "ip6"
		} else {
			network = "tcp4"
			netType = "ip4"
		}

		dialer = func(ctx context.Context, _, addr string) (net.Conn, error) {
			d := net.Dialer{
				LocalAddr: &net.TCPAddr{IP: localIP},
			}
			return d.DialContext(ctx, network, addr)
		}
		resolver = &sock5Resolver{netType: netType}
	} else {
		log.Printf("Proxy: %s --(default)-> outbound", clientIP)
	}

	conf := &socks5.Config{Dial: dialer, Resolver: resolver}
	server, err := socks5.New(conf)
	if err != nil {
		log.Printf("Failed to create server for conn: %v", err)
		return
	}

	if err := server.ServeConn(conn); err != nil {
		log.Printf("ServeConn error: %v", err)
	}
}
