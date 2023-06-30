package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net"
	"strings"

	"github.com/ziutek/telnet"
	"tailscale.com/tsnet"
	"tailscale.com/types/logger"
)

var (
	connectAddr = flag.String("connect", "127.0.0.1:23", "connect address")
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	var ts tsnet.Server
	ts.Logf = logger.Filtered(log.Printf, logFilter)
	status, err := ts.Up(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer ts.Close()

	listenAddr := status.TailscaleIPs[0].String() + ":23"

	listener, err := ts.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Printf("teltailnet: listening on: %s", listenAddr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	log.Printf("Accepted connection from %v", conn.RemoteAddr())

	remote, err := telnet.Dial("tcp", *connectAddr)
	if err != nil {
		log.Printf("teltailnet: dial error: %v, closing %v", err, conn.RemoteAddr())
		return
	}
	defer remote.Close()

	log.Printf("teltailnet: connected to %v", remote.RemoteAddr())

	// TODO: half-close dance.
	go func() {
		defer remote.Close()
		io.Copy(remote, conn)
	}()

	io.Copy(conn, remote)
}

func logFilter(s string) bool {
	return strings.Contains(s, "login.tailscale.com")
}
