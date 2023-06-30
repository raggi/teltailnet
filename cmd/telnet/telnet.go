package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/ziutek/telnet"
)

var connect = flag.String("connect", ":23", "connect address")

func main() {
	log.SetFlags(0)
	flag.Parse()

	t, err := telnet.Dial("tcp", *connect)
	if err != nil {
		log.Fatal(err)
	}
	defer t.Close()

	// TODO: proper half-close dance
	go io.Copy(os.Stdout, t)
	io.Copy(t, os.Stdin)
}
