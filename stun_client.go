// Standard Stun Client (RFC 5389)
//

package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

const MagicCookie int = 0x2112A442

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

	for {
		var address string

		fmt.Print("Enter address: ")
		_, err := fmt.Scan(&address)
		if err != nil {
			fmt.Println("Err Scan:", err)
			return
		}

		myAddr, publicPort := stunMe(conn, address)
		fmt.Printf("IP: %s, Port: %d\n", myAddr, publicPort)
	}
}

func stunMe(conn *net.UDPConn, address string) (net.IP, uint16) {

	remoteAddr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		fmt.Println("Err ResolveUDPAddr:", err)
		os.Exit(1)
	}

	// [ [00 bits reserved]
	// [Message Type 14 bits]
	// [Message Length 16 bits (2^4)]
	// [Magic Cookie 32 bits Int64]
	// [12 rand bytes (12^8 bits) ] ]

	buff := make([]byte, 20)
	buff[1] = 1
	buff[4] = byte(MagicCookie >> 24)
	buff[5] = byte((MagicCookie >> 16) & 0xff)
	buff[6] = byte((MagicCookie >> 8) & 0xff)
	buff[7] = byte(MagicCookie & 0xff)

	for i := 8; i < 20; i++ {
		buff[i] = uint8(rand.Intn(256))
	}

	_, err = conn.WriteToUDP(buff, remoteAddr)
	if err != nil {
		fmt.Println("Err WriteToUDP:", err)
		return nil, 0
	}

	conn.SetDeadline(time.Now().Add(5 * time.Second))
	// Read reply!
	rbuff := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(rbuff)
	if err != nil {
		fmt.Println("Err ReadFromUDP:", err)
	}

	conn.SetDeadline(time.Time{})
	for i := 20; i < n; {
		addrType := ((rbuff[i] & 255) | (rbuff[i+1] & 255)) & 0xff
		attrSize := ((rbuff[i+2] & 255) | (rbuff[i+3] & 255)) & 0xff

		i += 4
		if addrType != 0x0020 {
			// not an XOR-MAPPED_ADDRESS
			i += int(attrSize)
			continue
		}
		// Ignore fist byte + family byte
		i += 2
		publicPort := (uint16(rbuff[i]) << 8) | uint16(rbuff[i+1]) ^ uint16(MagicCookie>>16)
		i += 2

		addrBuff := make([]byte, 4)
		for j := 0; j < 4; j++ {
			addrBuff[j] = rbuff[i+j] ^ byte(MagicCookie>>(24-j*8))
		}
		i += 4

		return addrBuff, publicPort
	}
	return nil, 0
}
