package main

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/Jaywalker/iemitm/crc"
	"github.com/Jaywalker/iemitm/ie"
	"github.com/Jaywalker/iemitm/interprocess"
)

var crcDecoder *crc.CRC

func decompress(data []byte, from, to int) ([]byte, error) {
	b := bytes.NewReader(data[from : from+to])
	z, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer z.Close()
	p, err := ioutil.ReadAll(z)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func processPacket(packet interprocess.PacketData) {
	var header ie.IEHeader
	if err := binary.Read(bytes.NewReader(packet.Data), binary.BigEndian, &header); err != nil {
		fmt.Println("binary.Read failed:", err)
		return
	}

	crcChecker := crc.New()
	validatedCrc := crcChecker.Calculate(packet.Data[8:packet.Size], uint64(packet.Size-8))
	// fmt.Println("DEBUG:", hex.EncodeToString(buf[8:n]))

	if packet.Size == 36 { // Pre-Name, Post-Auth Ping
		fmt.Println("DPlay Ping/Pong")
	} else if packet.Size == ie.IEHeaderSize && header.FrameKind == 1 { // Ping!
		// fmt.Printf("IEHead PlayerFrom: %x PlayerTo: %x FrameKind: %x FrameNumber: %x FrameExpected: %x CRC32: %x - ", header.PlayerIDFrom, header.PlayerIDTo, header.FrameKind, header.FrameNumber, header.FrameExpected, header.CRC32) //, crc)
		// fmt.Println("")
		// fmt.Println("Ping!")
	} else {
		twoLetterIdent := string(packet.Data[ie.IEHeaderSize : ie.IEHeaderSize+2])
		if twoLetterIdent == "JM" {

			var jmHeader ie.JMPacketHeader
			if err := binary.Read(bytes.NewReader(packet.Data), binary.BigEndian, &jmHeader); err != nil {
				fmt.Println("binary.Read failed:", err)
				return
			}
			if header.Compressed == 1 {
				jmCompressed := ie.JMPacketCompressed{jmHeader, []byte{}}
				jmCompressed.Data = append(packet.Data[ie.JMPacketHeaderSize:packet.Size])
				fmt.Printf("IEHead PlayerFrom: %x PlayerTo: %x FrameKind: %x FrameNumber: %x FrameExpected: %x Compressed?: %x CRC32: %x vs %x - %c%c Unk1: %x Unk2: %x Unk3: %x Len: %d - ", jmCompressed.PlayerIDFrom, jmCompressed.PlayerIDTo, jmCompressed.FrameKind, jmCompressed.FrameNumber, jmCompressed.FrameExpected, jmCompressed.Compressed, jmCompressed.CRC32, validatedCrc, jmCompressed.JM[0], jmCompressed.JM[1], jmCompressed.Unknown1, jmCompressed.Unknown2, jmCompressed.Unknown3, uint8(jmCompressed.PacketLength))
				looksPromising := false
				compressedOffset := -1
				for i, v := range jmCompressed.Data {
					if !looksPromising && v == 0x78 {
						looksPromising = true
					} else if looksPromising {
						if v == 0x01 || v == 0x5e || v == 0x9c || v == 0xda {
							compressedOffset = i - 1
							break
						}
					}

				}
				if compressedOffset > -1 {
					if compressedOffset != 0 {
						fmt.Printf("Found data before compression ")
						for i := 0; i < compressedOffset; i++ {
							fmt.Printf("%x ", jmCompressed.Data[i])
						}
						fmt.Printf("- ")
					}

					decompressed, err := decompress(jmCompressed.Data, compressedOffset, (int(jmCompressed.PacketLength) - compressedOffset))
					if err != nil {
						fmt.Println("Failed to decompress data:", err)
						fmt.Println(packet.Source, " => ", packet.Dest, hex.EncodeToString(jmCompressed.Data[compressedOffset:(int(jmCompressed.PacketLength)-compressedOffset)]))
					} else {
						fmt.Println(packet.Source, " => ", packet.Dest, hex.EncodeToString(decompressed))
					}
				} else {
					fmt.Println("Could not find compressed data in the packet data!!")
				}
			} else {
				fmt.Printf("IEHead PlayerFrom: %x PlayerTo: %x FrameKind: %x FrameNumber: %x FrameExpected: %x Compressed?: %x CRC32: %x vs %x - %c%c Unk1: %x Unk2: %x Unk3: %x Len: %d - ", jmHeader.PlayerIDFrom, jmHeader.PlayerIDTo, jmHeader.FrameKind, jmHeader.FrameNumber, jmHeader.FrameExpected, jmHeader.Compressed, jmHeader.CRC32, validatedCrc, jmHeader.JM[0], jmHeader.JM[1], jmHeader.Unknown1, jmHeader.Unknown2, jmHeader.Unknown3, uint8(jmHeader.PacketLength))
				fmt.Println(packet.Source, " => ", packet.Dest, hex.EncodeToString(packet.Data[ie.JMPacketHeaderSize:packet.Size]))
			}
		} else {
			fmt.Println("Unhandled Two Letter Ident")
			fmt.Printf("IEHead PlayerFrom: %x PlayerTo: %x FrameKind: %x FrameNumber: %x FrameExpected: %x Compressed?: %x CRC32: %x vs %x - ", header.PlayerIDFrom, header.PlayerIDTo, header.FrameKind, header.FrameNumber, header.FrameExpected, header.Compressed, header.CRC32, validatedCrc)
			fmt.Println(packet.Source, " => ", packet.Dest, hex.EncodeToString(packet.Data[ie.IEHeaderSize:packet.Size]))
		}
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

	fmt.Println("Got connection to iemitm.. Data parsing will commence!")

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
