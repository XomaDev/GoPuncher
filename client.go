package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	const IP = "37.27.51.34"
	const PORT = 6688

	remoteAddr, err := net.ResolveUDPAddr("udp4", IP+":"+fmt.Sprint(PORT))
	if err != nil {
		fmt.Println("Err ResolveUDPAddr:", err)
		os.Exit(1)
	}

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

	_, err = conn.WriteToUDP([]byte("Hi!"), remoteAddr)
	if err != nil {
		fmt.Println("Err WriteToUDP:", err)
		return
	}

	// Read reply!
	buffer := make([]byte, 1024)
	n, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Err ReadFromUDP:", err)
	}
	response := string(buffer[:n])
	fmt.Println(response)
}
