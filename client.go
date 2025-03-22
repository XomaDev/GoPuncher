package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	for {
		var address string
		var port int

		fmt.Print("Enter address: ")
		_, err := fmt.Scan(&address)
		if err != nil {
			fmt.Println("Err Scan:", err)
			return
		}

		fmt.Print("Enter port: ")
		_, err = fmt.Scan(&port)
		if err != nil {
			fmt.Println("Err Scan:", err)
			return
		}

		stunRequest(address, port)
	}

}

func stunRequest(IP string, PORT int) {
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
