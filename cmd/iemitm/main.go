package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/Jaywalker/iemitm/dplay"
	"github.com/Jaywalker/iemitm/interprocess"
)

var listenerAddr string
var srvStrAddr string
var clientStrAddr string

var decoderSock *net.TCPConn

func UDPProxyListener(port string) {
	udpListenAddr, err := net.ResolveUDPAddr("udp", listenerAddr+port)
	if err != nil {
		panic(err)
	}

	udpListener, err := net.ListenUDP("udp", udpListenAddr)
	if err != nil {
		panic(err)
	}

	fmt.Println("UDP "+listenerAddr+port, " Listener Started")
	defer udpListener.Close()

	//May be needed for two way relay??
	cltAddr, err := net.ResolveUDPAddr("udp", clientStrAddr+port)
	if err != nil {
		panic(err)
	}
	clientOutSock, err := net.DialUDP("udp", nil, cltAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer clientOutSock.Close()

	srvAddr, err := net.ResolveUDPAddr("udp", srvStrAddr+port)
	if err != nil {
		panic(err)
	}
	srvOutSock, err := net.DialUDP("udp", nil, srvAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer srvOutSock.Close()

	buf := make([]byte, 0xffff)
	for {
		n, addr, err := udpListener.ReadFromUDP(buf)
		if err != nil {
			panic(err)
		}
		fmt.Println("UDP", addr, " => "+listenerAddr+port, " Received ", n, " bytes")

		b := buf[:n]

		forwardPacket := true
		forwardRespBuf := false
		var respPacket interprocess.RespPacketData
		if port != ":2350" { // DPlay ports. 2350 seems to be BG (or maybe IE) specific
			packet := dplay.NewDPlayPacket(b)
			if packet == nil {
				return
			}
			fmt.Println(packet)
			/*
				var header DPSP_MSG_HEADER
				if err := binary.Read(bytes.NewReader(buf), binary.LittleEndian, &header); err != nil {
					fmt.Println("binary.Read failed:", err)
					return
				}

				//Fix the Port, which for some reason is BigEndian
				header.SockAddr.Port = (header.SockAddr.Port >> 8) | (header.SockAddr.Port << 8)

				fmt.Println("	Sign: ", string(header.Signature[:]))
				fmt.Println("	Address Family: ", header.SockAddr.AddressFamily)
				fmt.Println("	IP: ", header.SockAddr.Address)
				fmt.Println("	Port: ", header.SockAddr.Port)
				fmt.Println("	Version: ", header.Version)
				fmt.Println("	Command: ", CommandToString(header.Command))
			*/

			if port == ":47624" {
				// TODO: Check if these listeners are already active and dont enable them if so
				go TCPProxyListener(":" + strconv.Itoa(packet.Port()))
				go UDPProxyListener(":" + strconv.Itoa(packet.Port()))
			}
		} else { // BG Port
			if decoderSock != nil {
				var networkBytes bytes.Buffer
				enc := gob.NewEncoder(&networkBytes)
				source := ""
				dest := ""
				if addr.IP.String() == srvStrAddr {
					source = srvStrAddr
					dest = listenerAddr
				} else {
					source = listenerAddr
					dest = srvStrAddr
				}
				data := &interprocess.PacketData{Source: source, Dest: dest, Port: port, Size: n, Data: buf}
				err := enc.Encode(data)
				if err != nil {
					fmt.Println("encode error:", err)
				} else {
					_, err = decoderSock.Write(networkBytes.Bytes())
					if err != nil {
						decoderSock.Close()
						decoderSock = nil
					} else {
						dec := gob.NewDecoder(decoderSock)

						err = dec.Decode(&respPacket)
						if err != nil {
							fmt.Println("decode error:", err)
						} else {
							forwardPacket = respPacket.Forward
							/*
								if len(respPacket.Data) > 0 {
									forwardRespBuf = true
								}
							*/
							if respPacket.Dest != "" {
								forwardRespBuf = true
							}
						}
					}
				}
			}
		}

		if forwardPacket {
			if addr.IP.String() == srvStrAddr {
				clientOutSock.Write(buf[:n])
				// fmt.Println("UDP", srvStrAddr+port, " => "+listenerAddr+port, " Sent", n, " bytes: ", hex.EncodeToString(buf[:n]))
			} else {
				srvOutSock.Write(buf[:n])
				// fmt.Println("UDP "+listenerAddr+port, " => ", srvStrAddr+port, " Sent", n, " bytes: ", hex.EncodeToString(buf[:n]))
			}
		}

		if forwardRespBuf {
			fmt.Println("ForwardBuf Found")
			if respPacket.Dest == "client" {
				clientOutSock.Write(respPacket.Data)
			} else {
				srvOutSock.Write(respPacket.Data)
			}
		}
	}
}

func TCPSocketRelay(src, dst *net.TCPConn, port string) {
	fmt.Println("TCP", src.RemoteAddr().String(), " => ", src.LocalAddr().String(), " - ", dst.LocalAddr().String(), " => ", dst.RemoteAddr().String(), " Relay Started")
	buf := make([]byte, 0xffff)
	for {
		n, err := src.Read(buf)
		if err != nil {
			//fmt.Printf("	Read failed '%s'\n", err)
			return
		}

		fmt.Println("TCP", src.RemoteAddr().String(), " => ", src.LocalAddr().String(), " Received", n, " bytes")
		b := buf[:n]

		//All TCP packets we've seen so far have been DPlay only
		packet := dplay.NewDPlayPacket(b)
		fmt.Println(packet)

		/*
			var header DPSP_MSG_HEADER
			if err := binary.Read(bytes.NewReader(buf), binary.LittleEndian, &header); err != nil {
				fmt.Println("binary.Read failed:", err)
				return
			}

			//Fix the Port, which for some reason is BigEndian
			header.SockAddr.Port = (header.SockAddr.Port >> 8) | (header.SockAddr.Port << 8)

			fmt.Println("	Sign: ", string(header.Signature[:]))
			fmt.Println("	Address Family: ", header.SockAddr.AddressFamily)
			fmt.Println("	IP: ", header.SockAddr.Address)
			fmt.Println("	Port: ", header.SockAddr.Port)
			fmt.Println("	Version: ", header.Version)
			fmt.Println("	Command: ", CommandToString(header.Command))
		*/

		/*
			if port == ":2300" && header.Command == DPSP_MSG_TYPE_ENUMSESSIONSREPLY {
				if gotEnumSessionReply {
					//Got both of our enum session replys, lets create our backchannel
					//From here on out, our relay is changed such that data on the server->:2300 channel gets sent via the client->:2300 channel
					tcpRemoteAddr, err := net.ResolveTCPAddr("tcp", clientStrAddr
					clientBackChannel2300 =
				}
				gotEnumSessionReply = true // Skip the first one as BG1 sends enums on both udp and tcp
			}
		*/

		/*
			if port == ":47624" {
				go TCPStaticListener(":" + strconv.Itoa(int(header.SockAddr.Port)))
				time.Sleep(500 * time.Millisecond)
			}
		*/

		//write out result
		n, err = dst.Write(b)
		if err != nil {
			fmt.Printf("	Write failed '%s'\n", err)
			return
		}
		fmt.Println("TCP", dst.LocalAddr().String(), " => ", dst.RemoteAddr().String(), " Sent", n, " bytes")
	}
}

var haxCounter int

func TCPConnHandler(src *net.TCPConn, port string) {
	fmt.Println("TCP", src.RemoteAddr().String(), " Handler Started")
	tcpRemoteAddr, err := net.ResolveTCPAddr("tcp", srvStrAddr+port)
	if err != nil {
		panic(err)
	}

	if port == ":9988" {
		// This is our decoder tool, handle it differently
		decoderSock = src
		// TCPDecoderTool(src,
		return
	}

	if strings.Split(src.RemoteAddr().String(), ":")[0] == srvStrAddr {
		//If the connection is from the server, we connect to the client
		tcpRemoteAddr, err = net.ResolveTCPAddr("tcp", clientStrAddr+port)
		if err != nil {
			panic(err)
		}
	}

	dst, err := net.DialTCP("tcp", nil, tcpRemoteAddr)
	// This whole bit is necessary for DPlay I guess? Our first connection to 47624 gets closed. Whyever, easy enough
	if port == ":47624" && haxCounter == 0 {
		fmt.Println("HaxCounter invoked. Goodbye")
		haxCounter++
		src.Close() // I guess in some versions this is important to do? I didn't have this in an old working version, I stopped working on this for a few years, I come back, reinstall what I believe is the exact same env, this no longer works. Will look into it further WAY-FUTURE-TODO
		dst.Close()
		/*
			go func() {
				time.Sleep(3 * time.Second)
				haxCounter = 0 // Reset our haxcounter so we can do connections in the future
			}()
		*/
		return
	}
	defer dst.Close()
	if err != nil {
		panic(err)
	}

	// Relay between src<->dst
	go TCPSocketRelay(src, dst, port)
	TCPSocketRelay(dst, src, port)
}

func TCPProxyListener(port string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", listenerAddr+port)
	if err != nil {
		panic(err)
	}

	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("TCP "+listenerAddr+port, " Listener Started")
	for {
		conn, err := tcpListener.AcceptTCP()
		if err != nil {
			panic(err)
		}

		fmt.Println("TCP", conn.RemoteAddr().String(), " => "+listenerAddr+":"+port, "- Got Connection!")
		go TCPConnHandler(conn, port)
	}
}

func main() {
	decoderSock = nil
	haxCounter = 0
	//	usingBackChannel = false
	//	gotEnumSessionReply = false
	//	clientBackChannel2300 = nil
	if len(os.Args) != 4 {
		os.Exit(0)
	}
	listenerAddr = os.Args[1]
	srvStrAddr = os.Args[3]
	clientStrAddr = os.Args[2]
	fmt.Println("DPlay MitM Activating...")
	fmt.Println("Fowarding", clientStrAddr, "to", srvStrAddr)
	go TCPProxyListener(":47624")
	go TCPProxyListener(":9988")
	go UDPProxyListener(":2350")
	UDPProxyListener(":47624")
	//go TCPProxyListener(":2300")
	//for {
	//}
}
