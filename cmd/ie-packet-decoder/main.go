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

	// crcChecker := crc.New()
	// validatedCrc := crcChecker.Calculate(packet.Data[8:packet.Size], uint64(packet.Size-8))

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
				jmCompressed := ie.JMPacketCompressed{jmHeader, 0, []byte{}}
				jmCompressed.Data = append(packet.Data[ie.JMPacketHeaderSize:packet.Size])
				compressedOffset := 1 // May not actually be 1 all the time, but the only compressed msg I get now that's not 0xFF it is
				if jmCompressed.SpecMsgFlag == 0xff {
					compressedOffset = 4
				}
				if compressedOffset != 1 {
					if err := binary.Read(bytes.NewReader(jmCompressed.Data[0:4]), binary.BigEndian, &jmCompressed.DecompressedSize); err != nil {
						fmt.Println("binary.Read failed:", err)
						fmt.Printf(jmCompressed.String() + " - ")
						fmt.Println(packet.Source, " => ", packet.Dest, hex.EncodeToString(jmCompressed.Data[compressedOffset:(int(jmCompressed.PacketLength)-compressedOffset)]))
						return
					}
					jmCompressed.Data = append(jmCompressed.Data[:compressedOffset], jmCompressed.Data[compressedOffset:]...)
				}

				decompressed, err := decompress(jmCompressed.Data, compressedOffset, (int(jmCompressed.PacketLength) - compressedOffset))

				if err != nil {
					fmt.Println("Failed to decompress data:", err)
					fmt.Printf(jmCompressed.String() + " - ")
					fmt.Println(packet.Source, " => ", packet.Dest, hex.EncodeToString(jmCompressed.Data[compressedOffset:(int(jmCompressed.PacketLength)-compressedOffset)]))
				} else {
					switch jmHeader.SpecMsgType {
					case ie.IE_SPEC_MSG_TYPE_CHAR_ARBITRATION:
						switch jmHeader.SpecMsgSubtype {
						case ie.IE_SPEC_MSG_SUBTYPE_UPDATE_SERVER_ARBITRATION_INFO:
							var servStatus ie.IECharArbServerStatus
							if err := binary.Read(bytes.NewReader(decompressed), binary.BigEndian, &servStatus); err != nil {
								fmt.Println("binary.Read failed:", err)
								fmt.Printf(jmCompressed.String() + " - ")
								fmt.Println(packet.Source, " => ", packet.Dest, hex.EncodeToString(decompressed))
							} else {
								fmt.Println("--------------------------")
								fmt.Println(servStatus.String())
								fmt.Println("--------------------------")
							}
						default:
							fmt.Printf(jmCompressed.String() + " - ")
							fmt.Println(packet.Source, " => ", packet.Dest, hex.EncodeToString(decompressed))
						}
					default:
						fmt.Printf(jmCompressed.String() + " - ")
						fmt.Println(packet.Source, " => ", packet.Dest, hex.EncodeToString(decompressed))
					}
				}
			} else {
				switch jmHeader.SpecMsgType {
				case ie.IE_SPEC_MSG_TYPE_CHAR_ARBITRATION:
					switch jmHeader.SpecMsgSubtype {
					case ie.IE_SPEC_MSG_SUBTYPE_TOGGLE_CHAR_READY:
						var charReady ie.IECharArbToggleCharReady
						if err := binary.Read(bytes.NewReader(packet.Data), binary.BigEndian, &charReady); err != nil {
							fmt.Println("binary.Read failed:", err)
							fmt.Printf(jmHeader.String())
							fmt.Println(packet.Source, " => ", packet.Dest, hex.EncodeToString(packet.Data[ie.JMPacketHeaderSize:packet.Size]))
						} else {
							fmt.Printf("Player %x Indicates %s\n", jmHeader.PlayerIDFrom, charReady.String())
						}
					default:
						fmt.Printf(jmHeader.String() + " - ")
						fmt.Println(packet.Source, " => ", packet.Dest, hex.EncodeToString(packet.Data[ie.JMPacketHeaderSize:packet.Size]))
					}
				default:
					fmt.Printf(jmHeader.String() + " - ")
					fmt.Println(packet.Source, " => ", packet.Dest, hex.EncodeToString(packet.Data[ie.JMPacketHeaderSize:packet.Size]))
				}
			}
		} else {
			fmt.Println("Unhandled Two Letter Ident")
			fmt.Printf(header.String() + " - ")
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
