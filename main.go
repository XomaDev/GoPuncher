package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type ReplyPacket struct {
	IP   string `json:"IP"`
	PORT int    `json:"PORT"`
}

func main() {
	const port = 6688
	addr, err := net.ResolveUDPAddr("udp4", "0.0.0.0:"+fmt.Sprint(port))
	if err != nil {
		fmt.Println("Error ResolveUDPAddr:", err)
		os.Exit(1)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error ListenUDP:", err)
		os.Exit(1)
	}

	defer conn.Close()
	fmt.Println("Listening on port", port)

	buffer := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error ReadFromUDP:", err)
			continue
		}

		fmt.Printf("Received %s from %s\n", string(buffer[:n]), addr.String())

		reply := ReplyPacket{
			IP:   addr.IP.String(),
			PORT: addr.Port,
		}

		jsonData, err := json.Marshal(reply)
		if err != nil {
			fmt.Println("Error Marshal:", err)
			continue
		}
		fmt.Println(string(jsonData))
		_, err = conn.WriteToUDP(jsonData, addr)
		if err != nil {
			fmt.Println("Error WriteToUDP:", err)
			continue
		}
	}
}
