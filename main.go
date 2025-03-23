package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/ristretto"
	"net"
	"os"
	"strconv"
	"time"
)

type RequestPacket struct {
	ID   string `json:"ID"`
	FIND bool   `json:"FIND"`
}

type ReplyPacket struct {
	SUCCESS bool   `json:"SUCCESS"`
	IP      string `json:"IP"`
	PORT    int    `json:"PORT"`
}

func main() {
	port := 8242
	if len(os.Args) > 1 {
		if p, err := strconv.Atoi(os.Args[1]); err == nil {
			port = p
		}
	}
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,       // Number of keys to track frequency
		MaxCost:     100 << 20, // 100MB max cache size
		BufferItems: 64,        // Number of keys per eviction buffer
	})
	if err != nil {
		fmt.Println("Err NewCache:", err)
		os.Exit(1)
	}

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

	buffer := make([]byte, 50)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error ReadFromUDP:", err)
			continue
		}

		var request RequestPacket
		err = json.Unmarshal(buffer[:n], &request)
		if err != nil {
			fmt.Println("Received bad JSON:", err)
			continue
		}

		var reply ReplyPacket
		var success bool

		if request.FIND {
			// We gotta connect to 'em!
			if val, found := cache.Get(request.ID); found {
				cache.Del(request.ID)
				reply = val.(ReplyPacket)
				success = true
			}
		} else {
			// Someone wants to connect to us!
			reply = ReplyPacket{
				SUCCESS: true,
				IP:      addr.IP.String(),
				PORT:    addr.Port,
			}
			success = true

			cache.SetWithTTL(request.ID, reply, 1, 15*time.Second)
		}
		if !success {
			reply = ReplyPacket{
				SUCCESS: false,
				IP:      "",
				PORT:    0,
			}
		}

		jsonData, err := json.Marshal(reply)
		if err != nil {
			fmt.Println("Error Marshal:", err)
			continue
		}

		_, err = conn.WriteToUDP(jsonData, addr)
		if err != nil {
			fmt.Println("Error WriteToUDP:", err)
			continue
		}
	}
}
