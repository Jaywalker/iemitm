package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/Jaywalker/iemitm/ie"
	"github.com/Jaywalker/iemitm/interprocess"
)

func processJMPacket(packet interprocess.PacketData, header ie.IEHeader) (forward bool) {
	defer fmt.Fprintln(rl, "--------------------------")
	printDebug("FULL: " + packet.Source + " => " + packet.Dest + header.String() + " - " + hex.EncodeToString(packet.Data[:packet.Size]))
	forward = true // Our default action is to forward the packet

	jmPacket, err := ie.NewJMPacket(packet.Data, packet.Size)
	if err != nil {
		fmt.Fprintln(rl, err.Error())
		return
	}

	var decompressed []byte
	if jmPacket.IsCompressed() {
		if jmPacket.PacketLength() > 0 {
			decompressed, err = decompress(jmPacket.PacketData())
			if err != nil {
				fmt.Fprintln(rl, "ERROR: Failed to decompress data:", err)
				return
			}
		}
	} else {
		decompressed = jmPacket.PacketData()
	}

	if jmPacket.IsSpecMsg() {
		printDebug("Spec Message! 0x%x", packet.Data[ie.JMHeaderSize:ie.JMHeaderSize+1])
		switch jmPacket.SpecType() {
		case ie.IE_SPEC_MSG_TYPE_MPSETTINGS:
			switch jmPacket.SpecSubType() {
			case ie.IE_SPEC_MSG_SUBTYPE_UPDATE_SERVER_ARBITRATION_INFO:
				var servStatus ie.IEMPSettingsFullSet
				if err := binary.Read(bytes.NewReader(decompressed), binary.BigEndian, &servStatus); err != nil {
					fmt.Fprintln(rl, "binary.Read failed:", err)
					fmt.Fprintln(rl, packet.Source, " => ", packet.Dest, ": ", jmPacket.String()+" - ", hex.EncodeToString(decompressed))
				} else {
					fmt.Fprintln(rl, servStatus.String())
				}
			case ie.IE_SPEC_MSG_SUBTYPE_TOGGLE_CHAR_READY:
				var charReady ie.IEMPSettingsToggleCharReady
				if err := binary.Read(bytes.NewReader(decompressed), binary.BigEndian, &charReady); err != nil {
					fmt.Fprintln(rl, "binary.Read failed:", err)
					fmt.Fprintln(rl, packet.Source, " => ", packet.Dest, ": ", jmPacket.String(), " - ", hex.EncodeToString(decompressed))
				} else {
					fmt.Fprintf(rl, "Player 0x%x Indicates %s\n", jmPacket.FromPlayerID(), charReady.String())
				}
			default:
				fmt.Fprintln(rl, "Unknown JM Spec Msg Type: ", packet.Source, " => ", packet.Dest, ": ", jmPacket.String(), " - ", hex.EncodeToString(decompressed))
			}
		case ie.IE_SPEC_MSG_TYPE_VERSION:
			switch jmPacket.SpecSubType() {
			case ie.IE_SPEC_MSG_SUBTYPE_VERSION_SERVER:
				var introHeader ie.IEVersionHeader
				if err := binary.Read(bytes.NewReader(decompressed[:ie.IEVersionHeaderSize]), binary.BigEndian, &introHeader); err != nil {
					fmt.Fprintln(rl, "binary.Read header failed:", err)
					fmt.Fprintln(rl, packet.Source, " => ", packet.Dest, ": ", jmPacket.String(), " - ", hex.EncodeToString(decompressed))
					return
				}
				var introFooter ie.IEVersionFooter
				if err := binary.Read(bytes.NewReader(decompressed[(ie.IEVersionHeaderSize+int(introHeader.VersionStringLen)):]), binary.BigEndian, &introFooter); err != nil {
					fmt.Fprintln(rl, "binary.Read footer failed:", err)
					fmt.Fprintln(rl, packet.Source, " => ", packet.Dest, ": ", jmPacket.String(), " - ", hex.EncodeToString(decompressed))
					return
				}
				intro := ie.IEVersion{introHeader, string(decompressed[ie.IEVersionHeaderSize:(ie.IEVersionHeaderSize + int(introHeader.VersionStringLen))]), introFooter}
				fmt.Fprintln(rl, packet.Source, " => ", packet.Dest, ": ", intro.String())

			default:
				fmt.Fprintln(rl, "Unknown JM Spec Msg Type: ", packet.Source, " => ", packet.Dest, ": ", jmPacket.String(), " - ", hex.EncodeToString(decompressed))
			}
		default:
			fmt.Fprintln(rl, "Unknown JM Spec Msg Type: ", packet.Source, " => ", packet.Dest, ": ", jmPacket.String()+" - ", hex.EncodeToString(decompressed))
		}
	} else {
		printDebug("Not a Spec Message! 0x%x", packet.Data[ie.JMHeaderSize:ie.JMHeaderSize+1])
		// Non-Spec messages are just messages from players
		ieMsg := ie.IEMsg{decompressed[0], string(decompressed[1:jmPacket.DataLength()])}
		fmt.Fprintln(rl, "Got Message: "+ieMsg.String())
	}
	return
}
