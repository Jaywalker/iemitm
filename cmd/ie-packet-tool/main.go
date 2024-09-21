package main

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/Jaywalker/iemitm/crc"
	"github.com/Jaywalker/iemitm/ie"
	"github.com/Jaywalker/iemitm/interprocess"
	"github.com/chzyer/readline"
)

var rl *readline.Instance

var crcDecoder *crc.CRC

var debug bool

var clientID uint32
var serverID uint32
var forwardPings bool
var forwardDplayPings bool
var clientFrameNumber uint16
var serverFrameNumber uint16
var clientExpectedFrameNumber uint16
var serverExpectedFrameNumber uint16

var crcChecker *crc.CRC

var sendBuf []byte
var sendBufTo string
var sendBufWaiting bool

func printDebug(str string, args ...any) {
	if debug {
		if !strings.HasSuffix(str, "\n") {
			str += "\n"
		}
		fmt.Fprintf(rl, "DEBUG: "+str, args...)
	}
}

func decompress(data []byte) ([]byte, error) {
	b := bytes.NewReader(data)
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

func processPacket(packet interprocess.PacketData) (forward bool) {
	forward = true
	var header ie.IEHeader
	if err := binary.Read(bytes.NewReader(packet.Data), binary.BigEndian, &header); err != nil {
		fmt.Fprintln(rl, "binary.Read failed:", err)
		return
	}

	if clientID == 0 {
		if header.PlayerIDFrom != serverID {
			clientID = header.PlayerIDFrom
		} else {
			clientID = header.PlayerIDTo
		}
	}

	if header.PlayerIDFrom != serverID {
		if clientFrameNumber < header.FrameNumber {
			clientFrameNumber = header.FrameNumber
			printDebug("Client frame number updated to %x", clientFrameNumber)
		} else if clientExpectedFrameNumber < header.FrameExpected {
			clientExpectedFrameNumber = header.FrameExpected
			printDebug("Client expected frame number updated to %x", clientExpectedFrameNumber)
		}
	} else if header.PlayerIDFrom == serverID {
		if serverFrameNumber < header.FrameNumber {
			serverFrameNumber = header.FrameNumber
			printDebug("Server frame number updated to %x", serverFrameNumber)
		} else if serverExpectedFrameNumber < header.FrameExpected {
			serverExpectedFrameNumber = header.FrameExpected
			printDebug("Server expected frame number updated to %x", serverExpectedFrameNumber)
		}
	}

	if packet.Size == 36 { // Pre-Name, Post-Auth Ping
		fmt.Fprintln(rl, "DPlay Ping/Pong")
		return forwardDplayPings
	} else if packet.Size == ie.IEHeaderSize && header.FrameKind == 1 { // Ping!
		// fmt.Fprintf(rl, "IEHead PlayerFrom: %x PlayerTo: %x FrameKind: %x FrameNumber: %x FrameExpected: %x CRC32: %x - Ping\n", header.PlayerIDFrom, header.PlayerIDTo, header.FrameKind, header.FrameNumber, header.FrameExpected, header.CRC32) //, crc)
		// fmt.Fprintln(rl, "")
		// fmt.Fprintln(rl, "Ping!")
		return forwardPings
	} else {
		twoLetterIdent := string(packet.Data[ie.IEHeaderSize : ie.IEHeaderSize+2])
		if twoLetterIdent == "JM" {
			return processJMPacket(packet, header)
		} else {
			fmt.Fprintln(rl, "Unhandled Two Letter Ident")
			fmt.Fprintln(rl, packet.Source, " => ", packet.Dest, header.String(), " - ", hex.EncodeToString(packet.Data[ie.IEHeaderSize:packet.Size]))
		}
	}
	return
}

func processJMPacket(packet interprocess.PacketData, header ie.IEHeader) (forward bool) {
	printDebug("FULL: " + packet.Source + " => " + packet.Dest + header.String() + " - " + hex.EncodeToString(packet.Data[:packet.Size]))
	forward = true
	var jmHeader ie.JMHeader
	if err := binary.Read(bytes.NewReader(packet.Data), binary.BigEndian, &jmHeader); err != nil {
		fmt.Fprintln(rl, "JMHeader binary.Read failed:", err)
		return
	}

	if packet.Data[ie.JMHeaderSize] == 0xff {
		printDebug("Spec Message!", packet.Data[ie.JMHeaderSize:ie.JMHeaderSize+1])
		if header.Compressed == 1 {
			printDebug("Compressed")
			var jmSpecHeaderCompressed ie.JMSpecHeaderCompressed
			if err := binary.Read(bytes.NewReader(packet.Data), binary.BigEndian, &jmSpecHeaderCompressed); err != nil {
				fmt.Fprintln(rl, "ERROR: JMSpecHeaderCompressed binary.Read failed:", err)
				return
			}
			jmSpecCompressed := ie.JMSpecCompressed{jmSpecHeaderCompressed, []byte{}}
			jmSpecCompressed.Data = append(packet.Data[ie.JMSpecHeaderCompressedSize:packet.Size])
		}
		/*
			switch jmHeader.SpecMsgType {
			case ie.IE_SPEC_MSG_TYPE_CHAR_ARBITRATION:
				switch jmHeader.SpecMsgSubtype {
				case ie.IE_SPEC_MSG_SUBTYPE_UPDATE_SERVER_ARBITRATION_INFO:
					var servStatus ie.IECharArbServerStatus
					if err := binary.Read(bytes.NewReader(decompressed), binary.BigEndian, &servStatus); err != nil {
						fmt.Fprintln(rl, "binary.Read failed:", err)
						fmt.Fprintln(rl, packet.Source, " => ", packet.Dest, jmCompressed.String()+" - ", hex.EncodeToString(decompressed))
					} else {
						fmt.Fprintln(rl, "--------------------------")
						fmt.Fprintln(rl, servStatus.String())
						fmt.Fprintln(rl, "--------------------------")
					}
				case ie.IE_SPEC_MSG_SUBTYPE_TOGGLE_CHAR_READY:
					var charReady ie.IECharArbToggleCharReady
					if err := binary.Read(bytes.NewReader(packet.Data), binary.BigEndian, &charReady); err != nil {
						fmt.Fprintln(rl, "binary.Read failed:", err)
						fmt.Fprintln(rl, packet.Source, " => ", packet.Dest, jmHeader.String(), " - ", hex.EncodeToString(packet.Data[ie.JMHeaderSize:packet.Size]))
					} else {
						fmt.Fprintf(rl, "Player %x Indicates %s\n", jmHeader.PlayerIDFrom, charReady.String())
					}
				case ie.IE_SPEC_MSG_TYPE_INTRO:
					switch jmHeader.SpecMsgSubtype {
					case ie.IE_SPEC_MSG_SUBTYPE_INTRO:
						var introHeader ie.IEIntroHeader
						if err := binary.Read(bytes.NewReader(packet.Data[:ie.IEIntroHeaderSize]), binary.BigEndian, &introHeader); err != nil {
							fmt.Fprintln(rl, "binary.Read header failed:", err)
							fmt.Fprintln(rl, packet.Source, " => ", packet.Dest, jmHeader.String(), " - ", hex.EncodeToString(packet.Data[ie.JMHeaderSize:packet.Size]))
							return
						}
						var introFooter ie.IEIntroFooter
						if err := binary.Read(bytes.NewReader(packet.Data[(ie.IEIntroHeaderSize+int(introHeader.VersionStringLen)):]), binary.BigEndian, &introFooter); err != nil {
							fmt.Fprintln(rl, "binary.Read footer failed:", err)
							fmt.Fprintln(rl, packet.Source, " => ", packet.Dest, jmHeader.String(), " - ", hex.EncodeToString(packet.Data[ie.JMHeaderSize:packet.Size]))
							return
						}
						intro := ie.IEIntro{introHeader, string(packet.Data[ie.IEIntroHeaderSize:(ie.IEIntroHeaderSize + int(introHeader.VersionStringLen))]), introFooter}
						fmt.Fprintln(rl, packet.Source, " => ", packet.Dest, intro.String())

					default:
						fmt.Fprintln(rl, packet.Source, " => ", packet.Dest, jmHeader.String(), " - ", hex.EncodeToString(packet.Data[ie.JMHeaderSize:packet.Size]))
					}
				default:
					fmt.Fprintln(rl, packet.Source, " => ", packet.Dest, jmCompressed.String()+" - ", hex.EncodeToString(decompressed))
				}
			default:
				fmt.Fprintln(rl, packet.Source, " => ", packet.Dest, jmCompressed.String()+" - ", hex.EncodeToString(decompressed))
				fmt.Fprintln(rl, "FULL:", packet.Source, " => ", packet.Dest, hex.EncodeToString(packet.Data[:packet.Size]))
			}
		*/
	} else {
		if header.Compressed == 1 {
			printDebug("Compressed")
			var jmHeaderCompressed ie.JMHeaderCompressed
			if err := binary.Read(bytes.NewReader(packet.Data), binary.BigEndian, &jmHeaderCompressed); err != nil {
				fmt.Fprintln(rl, "ERROR: JMHeaderCompressed binary.Read failed:", err)
				return
			}
			jmCompressed := ie.JMCompressed{jmHeaderCompressed, []byte{}}
			jmCompressed.Data = append(packet.Data[ie.JMHeaderCompressedSize:packet.Size])
			decompressed, err := decompress(jmCompressed.Data)
			if err != nil {
				fmt.Fprintln(rl, "ERROR: Failed to decompress data:", err)
				// fmt.Fprintln(rl, packet.Source, " => ", packet.Dest, jmCompressed.String()+" - ", hex.EncodeToString(jmCompressed.Data[compressedOffset:(jmCompressed.PacketLength-uint16(compressedOffset))]))
				fmt.Fprintln(rl, "FULL:", packet.Source, " => ", packet.Dest, header.String(), " - ", hex.EncodeToString(packet.Data[:packet.Size]))
				return
			}
			ieMsg := ie.IEMsg{decompressed[1], string(decompressed[1:])}
			ieMsgCompressed := ie.IEMsgDecompressed{jmHeaderCompressed, ieMsg}
			fmt.Fprintln(rl, "Got Message: "+ieMsgCompressed.String())
		} else {
			var jmHeader ie.JMHeader
			if err := binary.Read(bytes.NewReader(packet.Data), binary.BigEndian, &jmHeader); err != nil {
				fmt.Fprintln(rl, "ERROR: JMHeader binary.Read failed:", err)
				return
			}
			ieMsg := ie.IEMsg{packet.Data[ie.JMHeaderSize], string(packet.Data[(ie.JMHeaderSize + 1):packet.Size])}
			fmt.Fprintln(rl, "Got Message: "+ieMsg.String())
		}
	}
	/*
			jmCompressed := ie.JMCompressed{ie.JMHeaderCompressed{jmHeader, 0}, []byte{}}
			jmCompressed.Data = append(packet.Data[ie.JMHeaderSize:packet.Size])
			compressedOffset := 1 // May not actually be 1 all the time, but the only compressed msg I get now that's not 0xFF it is
			if jmCompressed.Data[0] == 0xff {
				compressedOffset = 4
			}
			if compressedOffset != 1 {
				if err := binary.Read(bytes.NewReader(jmCompressed.Data[0:4]), binary.BigEndian, &jmCompressed.DecompressedSize); err != nil {
					fmt.Fprintln(rl, "binary.Read failed:", err)
					fmt.Fprintln(rl, packet.Source, " => ", packet.Dest, jmCompressed.String()+" - ", hex.EncodeToString(jmCompressed.Data[compressedOffset:(jmCompressed.PacketLength-uint16(compressedOffset))]))
					return
				}
				jmCompressed.Data = append(jmCompressed.Data[:compressedOffset], jmCompressed.Data[compressedOffset:]...)
			}

			decompressed, err := decompress(jmCompressed.Data, compressedOffset, int(jmCompressed.PacketLength-uint16(compressedOffset)))

			if err != nil {
				fmt.Fprintln(rl, "Failed to decompress data:", err)
				fmt.Fprintln(rl, packet.Source, " => ", packet.Dest, jmCompressed.String()+" - ", hex.EncodeToString(jmCompressed.Data[compressedOffset:(jmCompressed.PacketLength-uint16(compressedOffset))]))
			} else {
				switch jmHeader.SpecMsgType {
				case ie.IE_SPEC_MSG_TYPE_CHAR_ARBITRATION:
					switch jmHeader.SpecMsgSubtype {
					case ie.IE_SPEC_MSG_SUBTYPE_UPDATE_SERVER_ARBITRATION_INFO:
						var servStatus ie.IECharArbServerStatus
						if err := binary.Read(bytes.NewReader(decompressed), binary.BigEndian, &servStatus); err != nil {
							fmt.Fprintln(rl, "binary.Read failed:", err)
							fmt.Fprintln(rl, packet.Source, " => ", packet.Dest, jmCompressed.String()+" - ", hex.EncodeToString(decompressed))
						} else {
							fmt.Fprintln(rl, "--------------------------")
							fmt.Fprintln(rl, servStatus.String())
							fmt.Fprintln(rl, "--------------------------")
						}
					default:
						fmt.Fprintln(rl, packet.Source, " => ", packet.Dest, jmCompressed.String()+" - ", hex.EncodeToString(decompressed))
					}
				default:
					fmt.Fprintln(rl, packet.Source, " => ", packet.Dest, jmCompressed.String()+" - ", hex.EncodeToString(decompressed))
					fmt.Fprintln(rl, "FULL:", packet.Source, " => ", packet.Dest, hex.EncodeToString(packet.Data[:packet.Size]))
				}
			}
		} else {
			switch jmHeader.SpecMsgType {
			case ie.IE_SPEC_MSG_TYPE_CHAR_ARBITRATION:
				switch jmHeader.SpecMsgSubtype {
				case ie.IE_SPEC_MSG_SUBTYPE_TOGGLE_CHAR_READY:
					var charReady ie.IECharArbToggleCharReady
					if err := binary.Read(bytes.NewReader(packet.Data), binary.BigEndian, &charReady); err != nil {
						fmt.Fprintln(rl, "binary.Read failed:", err)
						fmt.Fprintln(rl, packet.Source, " => ", packet.Dest, jmHeader.String(), " - ", hex.EncodeToString(packet.Data[ie.JMHeaderSize:packet.Size]))
					} else {
						fmt.Fprintf(rl, "Player %x Indicates %s\n", jmHeader.PlayerIDFrom, charReady.String())
					}
				default:
					fmt.Fprintln(rl, packet.Source, " => ", packet.Dest, jmHeader.String(), " - ", hex.EncodeToString(packet.Data[ie.JMHeaderSize:packet.Size]))
				}
			case ie.IE_SPEC_MSG_TYPE_INTRO:
				switch jmHeader.SpecMsgSubtype {
				case ie.IE_SPEC_MSG_SUBTYPE_INTRO:
					var introHeader ie.IEIntroHeader
					if err := binary.Read(bytes.NewReader(packet.Data[:ie.IEIntroHeaderSize]), binary.BigEndian, &introHeader); err != nil {
						fmt.Fprintln(rl, "binary.Read header failed:", err)
						fmt.Fprintln(rl, packet.Source, " => ", packet.Dest, jmHeader.String(), " - ", hex.EncodeToString(packet.Data[ie.JMHeaderSize:packet.Size]))
						return
					}
					var introFooter ie.IEIntroFooter
					if err := binary.Read(bytes.NewReader(packet.Data[(ie.IEIntroHeaderSize+int(introHeader.VersionStringLen)):]), binary.BigEndian, &introFooter); err != nil {
						fmt.Fprintln(rl, "binary.Read footer failed:", err)
						fmt.Fprintln(rl, packet.Source, " => ", packet.Dest, jmHeader.String(), " - ", hex.EncodeToString(packet.Data[ie.JMHeaderSize:packet.Size]))
						return
					}
					intro := ie.IEIntro{introHeader, string(packet.Data[ie.IEIntroHeaderSize:(ie.IEIntroHeaderSize + int(introHeader.VersionStringLen))]), introFooter}
					fmt.Fprintln(rl, packet.Source, " => ", packet.Dest, intro.String())

				default:
					fmt.Fprintln(rl, packet.Source, " => ", packet.Dest, jmHeader.String(), " - ", hex.EncodeToString(packet.Data[ie.JMHeaderSize:packet.Size]))
				}
			default:
				fmt.Fprintln(rl, packet.Source, " => ", packet.Dest, jmHeader.String(), " - ", hex.EncodeToString(packet.Data[ie.JMHeaderSize:packet.Size]))
				fmt.Fprintln(rl, "FULL:", packet.Source, " => ", packet.Dest, hex.EncodeToString(packet.Data[:packet.Size]))
			}
		}
	*/
	return
}

var completer = readline.NewPrefixCompleter(
	readline.PcItem("sendraw",
		readline.PcItem("client"),
		readline.PcItem("server"),
	),
	readline.PcItem("sendmsg",
		readline.PcItem("client"),
		readline.PcItem("server"),
	),
	readline.PcItem("dplay",
		readline.PcItem("pings",
			readline.PcItem("enable"),
			readline.PcItem("disable"),
		),
	),
	readline.PcItem("pings",
		readline.PcItem("enable"),
		readline.PcItem("disable"),
	),
	readline.PcItem("debug"),
	readline.PcItem("exit"),
	readline.PcItem("quit"),
	readline.PcItem("bye"),
)

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func main() {
	crcChecker = crc.New()

	sendBufWaiting = false
	debug = false
	forwardPings = true
	clientID = 0
	serverID = 0x1000000
	clientFrameNumber = 0
	serverFrameNumber = 0
	clientExpectedFrameNumber = 0
	serverExpectedFrameNumber = 0

	var err error
	rl, err = readline.NewEx(&readline.Config{
		UniqueEditLine:  true,
		Prompt:          "\033[31mie-packet-tool »\033[0m ",
		HistoryFile:     "/tmp/ie-packet-injector.tmp",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	done := make(chan struct{})

	// =================================================================================================

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

	// =================================================================================================

	go func() {
	loop:
		for {
			dec := gob.NewDecoder(conn)

			var packet interprocess.PacketData
			err = dec.Decode(&packet)
			if err != nil {
				fmt.Fprintln(rl, "decode error:", err)
				break
			} else {
				forward := processPacket(packet)

				encoder := gob.NewEncoder(conn)
				resp := &interprocess.RespPacketData{}
				resp.Forward = forward
				if sendBufWaiting {
					sendBufWaiting = false
					resp.Dest = sendBufTo
					resp.Data = append([]byte{}, sendBuf...)
					sendBuf = []byte{}
					sendBufTo = ""
				}
				encoder.Encode(resp)
				if err != nil {
					fmt.Fprintln(rl, "encode error:", err)
				}
			}

			select {
			case <-time.After(1 * time.Millisecond):
			case <-done:
				break loop
			}
		}
		done <- struct{}{}
	}()

	// =================================================================================================

	for {
		ln := rl.Line()
		if ln.CanContinue() {
			continue
		} else if ln.CanBreak() {
			break
		}
		line := strings.TrimSpace(ln.Line)
		switch {
		case strings.HasPrefix(line, "sendraw "):
			if strings.HasPrefix(line[8:], "server") {
				fmt.Fprintln(rl, "Sending to server")
			} else if strings.HasPrefix(line[8:], "client") {
				fmt.Fprintln(rl, "Sending to client")
			} else {
				fmt.Fprintln(rl, "Invalid target. Valid targets are: client server")
			}
		case strings.HasPrefix(line, "sendmsg "):
			if strings.HasPrefix(line[8:], "server") {
				fmt.Fprintln(rl, "Sending to server")

				serverFrameNumber += 1 // TODO: Do I need to add 1 to this before?

				msgPacket := ie.IEMsgPacket{}
				msgPacket.PlayerIDFrom = serverID
				msgPacket.PlayerIDTo = clientID
				msgPacket.FrameKind = 0
				msgPacket.FrameNumber = serverFrameNumber
				msgPacket.FrameExpected = serverExpectedFrameNumber
				msgPacket.Compressed = 0
				msgPacket.CRC32 = 0
				msgPacket.JM[0] = 'J'
				msgPacket.JM[1] = 'M'
				msgPacket.Unknown1 = 0
				msgPacket.Unknown2 = 1
				msgPacket.Message = strings.Replace(line, "sendmsg server ", "", 1)
				msgPacket.MessageLength = byte(len(msgPacket.Message))
				msgPacket.PacketLength = uint16(msgPacket.MessageLength) + 1

				serialbuf, err := msgPacket.Serialize()
				if err != nil {
					fmt.Fprintln(rl, "Error: failed to serialize ", err)
				} else {
					size := len(serialbuf)
					validatedCrc := crcChecker.Calculate(serialbuf[8:], uint32(size-8))
					msgPacket.CRC32 = validatedCrc
					serialbuf, err := msgPacket.Serialize()
					if err != nil {
						fmt.Fprintln(rl, "Error: failed to serialize ", err)
					} else {
						sendBuf = append([]byte{}, serialbuf...)
						sendBufTo = "server"
						sendBufWaiting = true
						fmt.Fprintln(rl, "Serialized: ", hex.EncodeToString(serialbuf))
					}
				}
			} else if strings.HasPrefix(line[8:], "client") {
				fmt.Fprintln(rl, "Sending to client")

				clientFrameNumber += 1 // TODO: Do I need to add 1 to this before?

				msgPacket := ie.IEMsgPacket{}
				msgPacket.PlayerIDFrom = serverID
				msgPacket.PlayerIDTo = clientID
				msgPacket.FrameKind = 0
				msgPacket.FrameNumber = clientFrameNumber
				msgPacket.FrameExpected = clientExpectedFrameNumber
				msgPacket.Compressed = 0
				msgPacket.CRC32 = 0
				msgPacket.JM[0] = 'J'
				msgPacket.JM[1] = 'M'
				msgPacket.Unknown1 = 0
				msgPacket.Unknown2 = 1
				msgPacket.Message = strings.Replace(line, "sendmsg client ", "", 1)
				msgPacket.MessageLength = byte(len(msgPacket.Message))
				msgPacket.PacketLength = uint16(msgPacket.MessageLength) + 1

				serialbuf, err := msgPacket.Serialize()
				if err != nil {
					fmt.Fprintln(rl, "Error: failed to serialize ", err)
				} else {
					size := len(serialbuf)
					validatedCrc := crcChecker.Calculate(serialbuf[8:], uint32(size-8))
					msgPacket.CRC32 = validatedCrc
					serialbuf, err := msgPacket.Serialize()
					if err != nil {
						fmt.Fprintln(rl, "Error: failed to serialize ", err)
					} else {
						sendBuf = append([]byte{}, serialbuf...)
						sendBufTo = "client"
						sendBufWaiting = true
						fmt.Fprintln(rl, "Serialized: ", hex.EncodeToString(serialbuf))
					}
				}
			} else {
				fmt.Fprintln(rl, "Invalid target. Valid targets are: client server")
			}
		case line == "pings disable":
			forwardPings = false
			fmt.Fprintln(rl, "Pings disabled.")
		case line == "pings enable":
			forwardPings = true
			fmt.Fprintln(rl, "Pings enabled.")
		case line == "dplay pings disable":
			forwardDplayPings = false
			fmt.Fprintln(rl, "Dplay pings disabled.")
		case line == "dplay pings enable":
			forwardDplayPings = true
			fmt.Fprintln(rl, "Dplay pings enabled.")
		case line == "debug":
			if !debug {
				fmt.Fprintln(rl, "Debug Enabled")
			} else {
				fmt.Fprintln(rl, "Debug Disabled")
			}
			debug = !debug
		// case line == "set char ready 0":
		case line == "exit":
			fallthrough
		case line == "quit":
			fallthrough
		case line == "bye":
			goto exit
		case line == "":
		default:
			fmt.Fprintln(rl, "Invalid Command:", strconv.Quote(line))
		}
	}
exit:
	rl.Clean()
	done <- struct{}{}
	<-done
}
