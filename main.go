package main

import (
	"flag"
	"fmt"
	"net"
)

var debug bool

func run() error {
	var pProto = flag.String("p", "udp", "listen protocol <udp, udp4, udp6>")
	var pAddr = flag.String("l", ":7777", "listen [address]<:port>")
	flag.Parse()
	udpAddr, err := net.ResolveUDPAddr(*pProto, *pAddr)
	if err != nil {
		return err
	}
	udpConn, err := net.ListenUDP(*pProto, udpAddr)
	if err != nil {
		return err
	}
	defer udpConn.Close()

	var buf = make([]byte, 1024)
	for {
		rn, from, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			if !err.(*net.OpError).Temporary() {
				return err
			}
			if debug {
				fmt.Printf("Temporary error while receiving request: %s\n", err.Error())
			}
		}
		if debug {
			fmt.Printf("Request from: %s (%d bytes)\n", from.String(), rn)
		}
		_, err = udpConn.WriteToUDP(buf[:rn], from)
		if err != nil {
			if !err.(*net.OpError).Temporary() {
				return err
			}
			if debug {
				fmt.Printf("Temporary error while sending reply: %s\n", err.Error())
			}
		}
	}
}

func main() {
	flag.BoolVar(&debug, "d", false, "debug")
	if err := run(); err != nil {
		if debug {
			panic(err)
		} else {
			fmt.Println(err)
		}
	}
}
