package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"log"
	"net"

	"github.com/Jaywalker/iemitm/crc"
	"github.com/Jaywalker/iemitm/ie"
	"github.com/Jaywalker/iemitm/interprocess"
)

var crcDecoder *crc.CRC

func processPacket(packet interprocess.PacketData) {
	var header ie.IEHeader
	if err := binary.Read(bytes.NewReader(packet.Data), binary.BigEndian, &header); err != nil {
		fmt.Println("binary.Read failed:", err)
		return
	}
	// crc := calculateCRC(buf[8:n], uint64(n-8))
	fmt.Printf("IEHead PlayerFrom: %x PlayerTo: %x FrameKind: %x FrameNumber: %x FrameExpected: %x CRC32: %x - ", header.PlayerIDFrom, header.PlayerIDTo, header.FrameKind, header.FrameNumber, header.FrameExpected, header.CRC32) //, crc)
	// fmt.Println("DEBUG:", hex.EncodeToString(buf[8:n]))

	if packet.Size == 36 { // Pre-Name, Post-Auth Ping
		fmt.Println("DPlay Ping/Pong")
	} else if packet.Size == 18 && header.FrameKind == 1 { // Ping!
		fmt.Println("Ping!")
	} else {
		fmt.Println("")
		/*
			if addr.IP.String() == srvStrAddr {
				fmt.Println("UDP", srvStrAddr+port, " => "+listenerAddr+port, hex.EncodeToString(buf[18:n]))
			} else {
				fmt.Println("UDP "+listenerAddr+port, " => ", srvStrAddr+port, hex.EncodeToString(buf[18:n]))
			}
		*/
	}

}

func main() {
	crcDecoder = crc.New()

	tcpRemoteAddr, err := net.ResolveTCPAddr("tcp", "192.168.122.1:9988")
	if err != nil {
		panic(err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpRemoteAddr)
	if err != nil {
		panic(err)
	}

	for {
		dec := gob.NewDecoder(conn)

		var packet interprocess.PacketData
		err = dec.Decode(&packet)
		if err != nil {
			log.Fatal("decode error:", err)
		}
		processPacket(packet)
	}
}
