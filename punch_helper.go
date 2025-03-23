package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp4", "0.0.0.0:0")
	if err != nil {
		fmt.Println("Err ResolveUDPAddr:", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		fmt.Println("Err ListenUDP:", err)
		os.Exit(1)
	}
	defer conn.Close()
}
