package ie

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
)

//NOTE: -Ish suffix imply a guess, not fact

type IEHeader struct {
	PlayerIDFrom  uint32
	PlayerIDTo    uint32 //Server seems to be 0001 //
	FrameKind     byte
	FrameNumber   uint16
	FrameExpected uint16
	Compressed    byte // This tells us if we're using Zlib compression (Best Compression flag)
	CRC32         uint32
}

func (header IEHeader) String() string {
	return fmt.Sprintf("IEHead PlayerFrom: %x PlayerTo: %x FrameKind: %x FrameNumber: %x FrameExpected: %x Compressed?: %x CRC32: %x", header.PlayerIDFrom, header.PlayerIDTo, header.FrameKind, header.FrameNumber, header.FrameExpected, header.Compressed, header.CRC32)
}

const IEHeaderSize int = 18

const IE_SPEC_MSG_TYPE_CHAR_ARBITRATION byte = 0x4d

type JMPacketHeader struct {
	IEHeader
	JM             [2]byte
	Unknown1       byte // 00
	Unknown2       byte // 01
	PacketLength   uint16
	SpecMsgFlag    byte // TODO: change all bytes to uint8s?
	SpecMsgType    byte
	SpecMsgSubtype byte
}

func (jmHeader JMPacketHeader) String() string {
	return fmt.Sprintf("IEHead PlayerFrom: %x PlayerTo: %x FrameKind: %x FrameNumber: %x FrameExpected: %x Compressed?: %x CRC32: %x - %c%c Unk1: %x Unk2: %x Len: %d SpecMsgFlag: %x SpecMsgType: %x SpecMsgSubtype: %x", jmHeader.PlayerIDFrom, jmHeader.PlayerIDTo, jmHeader.FrameKind, jmHeader.FrameNumber, jmHeader.FrameExpected, jmHeader.Compressed, jmHeader.CRC32, jmHeader.JM[0], jmHeader.JM[1], jmHeader.Unknown1, jmHeader.Unknown2, jmHeader.PacketLength, jmHeader.SpecMsgFlag, jmHeader.SpecMsgType, jmHeader.SpecMsgSubtype)
}

const JMPacketHeaderSize int = IEHeaderSize + 9

type JMPacketCompressed struct {
	JMPacketHeader
	DecompressedSize uint32
	Data             []byte
}

func (jmCompressed JMPacketCompressed) String() string {
	return fmt.Sprintf("IEHead PlayerFrom: %x PlayerTo: %x FrameKind: %x FrameNumber: %x FrameExpected: %x Compressed?: %x CRC32: %x - %c%c Unk1: %x Unk2: %x Len: %d SpecMsgFlag: %x SpecMsgType: %x SpecMsgSubtype: %x DecompressedSize: %x", jmCompressed.PlayerIDFrom, jmCompressed.PlayerIDTo, jmCompressed.FrameKind, jmCompressed.FrameNumber, jmCompressed.FrameExpected, jmCompressed.Compressed, jmCompressed.CRC32, jmCompressed.JM[0], jmCompressed.JM[1], jmCompressed.Unknown1, jmCompressed.Unknown2, jmCompressed.PacketLength, jmCompressed.SpecMsgFlag, jmCompressed.SpecMsgType, jmCompressed.SpecMsgSubtype, jmCompressed.DecompressedSize)
}

type IEMsgPacket struct {
	IEHeader
	JM            [2]byte
	Unknown1      byte // 00
	Unknown2      byte // 01
	PacketLength  uint16
	MessageLength byte
	Message       string
}

func (iemsg IEMsgPacket) Serialize() ([]byte, error) {
	var serialbuf []byte
	var err error
	serialbuf, err = binary.Append(serialbuf, binary.BigEndian, iemsg.PlayerIDFrom)
	if err != nil {
		return nil, err
	}
	serialbuf, err = binary.Append(serialbuf, binary.BigEndian, iemsg.PlayerIDTo)
	if err != nil {
		return nil, err
	}
	serialbuf, err = binary.Append(serialbuf, binary.BigEndian, iemsg.FrameKind)
	if err != nil {
		return nil, err
	}
	serialbuf, err = binary.Append(serialbuf, binary.BigEndian, iemsg.FrameNumber)
	if err != nil {
		return nil, err
	}
	serialbuf, err = binary.Append(serialbuf, binary.BigEndian, iemsg.FrameExpected)
	if err != nil {
		return nil, err
	}
	serialbuf, err = binary.Append(serialbuf, binary.BigEndian, iemsg.Compressed)
	if err != nil {
		return nil, err
	}
	serialbuf, err = binary.Append(serialbuf, binary.BigEndian, iemsg.CRC32)
	if err != nil {
		return nil, err
	}
	serialbuf, err = binary.Append(serialbuf, binary.BigEndian, iemsg.JM)
	if err != nil {
		return nil, err
	}
	serialbuf, err = binary.Append(serialbuf, binary.BigEndian, iemsg.Unknown1)
	if err != nil {
		return nil, err
	}
	serialbuf, err = binary.Append(serialbuf, binary.BigEndian, iemsg.Unknown2)
	if err != nil {
		return nil, err
	}
	serialbuf, err = binary.Append(serialbuf, binary.BigEndian, iemsg.PacketLength)
	if err != nil {
		return nil, err
	}
	serialbuf, err = binary.Append(serialbuf, binary.BigEndian, iemsg.MessageLength)
	if err != nil {
		return nil, err
	}
	serialbuf, err = binary.Append(serialbuf, binary.BigEndian, []byte(iemsg.Message))
	if err != nil {
		return nil, err
	}
	return serialbuf, nil
}

const IE_SPEC_MSG_SUBTYPE_TOGGLE_CHAR_READY byte = 0x72

type IECharArbToggleCharReady struct {
	JMPacketHeader
	CharacterNum byte
	ReadyStatus  uint32
}

const IECharArbToggleCharReadySize int = JMPacketHeaderSize + 5

func (charReady IECharArbToggleCharReady) String() string {
	if charReady.ReadyStatus == 0 {
		return "Character " + strconv.Itoa(int(charReady.CharacterNum)) + " is not ready."
	} else if charReady.ReadyStatus == 1 {
		return "Character " + strconv.Itoa(int(charReady.CharacterNum)) + " is ready."
	}
	return "Character " + strconv.Itoa(int(charReady.CharacterNum)) + " is ready (" + strconv.Itoa(int(charReady.ReadyStatus)) + ")."
}

const IE_SPEC_MSG_SUBTYPE_UPDATE_SERVER_ARBITRATION_INFO byte = 0x53

type IECharArbServerStatus struct {
	Unknown1                    [2]byte
	DefaultPermBuyAndSell       byte
	DefaultPermTravel           byte
	DefaultPermDialog           byte
	DefaultPermViewCharacters   byte
	DefaultPermPause            byte
	DefaultPermHasBeenLeaderIsh byte
	DefaultPermLeader           byte
	DefaultPermModifyCharacters byte
	Unknown2                    [20]byte
	Player1ID                   uint32
	Player1PermBuyAndSell       byte
	Player1PermTravel           byte
	Player1PermDialog           byte
	Player1PermViewCharacters   byte
	Player1PermPause            byte
	Player1PermHasBeenLeaderIsh byte
	Player1PermLeader           byte
	Player1PermModifyCharacters byte
	Unknown3                    [89]byte
	CharIsReady                 [6]byte
	Unknown4                    [6]byte
	CharOwnerPlayerID           [6]uint32
	ImportCharSettings          byte
	RestrictStores              byte
	ListenToJoinRequests        byte
	Unknown5                    [29]byte
}

func (charArbServStatus IECharArbServerStatus) String() string {
	ret := "Unk1: " + hex.EncodeToString(charArbServStatus.Unknown1[:])

	ret += "\nDefaultPerms:"

	if charArbServStatus.DefaultPermBuyAndSell == 0 {
		ret += "\n\tBuyAndSell: no"
	} else if charArbServStatus.DefaultPermBuyAndSell == 1 {
		ret += "\n\tBuyAndSell: yes"
	} else {
		ret += fmt.Sprintf("\n\tBuyAndSell: yes? 0x%x", charArbServStatus.DefaultPermBuyAndSell)
	}

	if charArbServStatus.DefaultPermTravel == 0 {
		ret += "\n\tTravel: no"
	} else if charArbServStatus.DefaultPermTravel == 1 {
		ret += "\n\tTravel: yes"
	} else {
		ret += fmt.Sprintf("\n\tTravel: yes? 0x%x", charArbServStatus.DefaultPermTravel)
	}

	if charArbServStatus.DefaultPermDialog == 0 {
		ret += "\n\tDialog: no"
	} else if charArbServStatus.DefaultPermDialog == 1 {
		ret += "\n\tDialog: yes"
	} else {
		ret += fmt.Sprintf("\n\tDialog: yes? 0x%x", charArbServStatus.DefaultPermDialog)
	}

	if charArbServStatus.DefaultPermViewCharacters == 0 {
		ret += "\n\tView Characters: no"
	} else if charArbServStatus.DefaultPermViewCharacters == 1 {
		ret += "\n\tView Characters: yes"
	} else {
		ret += fmt.Sprintf("\n\tView Characters: yes? 0x%x", charArbServStatus.DefaultPermViewCharacters)
	}

	if charArbServStatus.DefaultPermPause == 0 {
		ret += "\n\tPause: no"
	} else if charArbServStatus.DefaultPermPause == 1 {
		ret += "\n\tPause: yes"
	} else {
		ret += fmt.Sprintf("\n\tPause: yes? 0x%x", charArbServStatus.DefaultPermPause)
	}

	if charArbServStatus.DefaultPermHasBeenLeaderIsh == 0 {
		ret += "\n\tMaybe 'HasBeenLeader': no"
	} else if charArbServStatus.DefaultPermHasBeenLeaderIsh == 1 {
		ret += "\n\tMaybe 'HasBeenLeader': yes"
	} else {
		ret += fmt.Sprintf("\n\tMaybe 'HasBeenLeader': yes? 0x%x", charArbServStatus.DefaultPermHasBeenLeaderIsh)
	}

	if charArbServStatus.DefaultPermLeader == 0 {
		ret += "\n\tLeader: no"
	} else if charArbServStatus.DefaultPermLeader == 1 {
		ret += "\n\tLeader: yes"
	} else {
		ret += fmt.Sprintf("\n\tLeader: yes? 0x%x", charArbServStatus.DefaultPermLeader)
	}

	if charArbServStatus.DefaultPermModifyCharacters == 0 {
		ret += "\n\tModify Characters: no"
	} else if charArbServStatus.DefaultPermModifyCharacters == 1 {
		ret += "\n\tModify Characters: yes"
	} else {
		ret += fmt.Sprintf("\n\tModify Characters: yes? 0x%x", charArbServStatus.DefaultPermModifyCharacters)
	}

	ret += "\nUnk2: " + hex.EncodeToString(charArbServStatus.Unknown2[:])

	ret += fmt.Sprintf("\nPlayer1Perms (%x):", charArbServStatus.Player1ID)

	if charArbServStatus.Player1PermBuyAndSell == 0 {
		ret += "\n\tBuyAndSell: no"
	} else if charArbServStatus.Player1PermBuyAndSell == 1 {
		ret += "\n\tBuyAndSell: yes"
	} else {
		ret += fmt.Sprintf("\n\tBuyAndSell: yes? 0x%x", charArbServStatus.Player1PermBuyAndSell)
	}

	if charArbServStatus.Player1PermTravel == 0 {
		ret += "\n\tTravel: no"
	} else if charArbServStatus.Player1PermTravel == 1 {
		ret += "\n\tTravel: yes"
	} else {
		ret += fmt.Sprintf("\n\tTravel: yes? 0x%x", charArbServStatus.Player1PermTravel)
	}

	if charArbServStatus.Player1PermDialog == 0 {
		ret += "\n\tDialog: no"
	} else if charArbServStatus.Player1PermDialog == 1 {
		ret += "\n\tDialog: yes"
	} else {
		ret += fmt.Sprintf("\n\tDialog: yes? 0x%x", charArbServStatus.Player1PermDialog)
	}

	if charArbServStatus.Player1PermViewCharacters == 0 {
		ret += "\n\tView Characters: no"
	} else if charArbServStatus.Player1PermViewCharacters == 1 {
		ret += "\n\tView Characters: yes"
	} else {
		ret += fmt.Sprintf("\n\tView Characters: yes? 0x%x", charArbServStatus.Player1PermViewCharacters)
	}

	if charArbServStatus.Player1PermPause == 0 {
		ret += "\n\tPause: no"
	} else if charArbServStatus.Player1PermPause == 1 {
		ret += "\n\tPause: yes"
	} else {
		ret += fmt.Sprintf("\n\tPause: yes? 0x%x", charArbServStatus.Player1PermPause)
	}

	if charArbServStatus.Player1PermHasBeenLeaderIsh == 0 {
		ret += "\n\tMaybe 'HasBeenLeader': no"
	} else if charArbServStatus.Player1PermHasBeenLeaderIsh == 1 {
		ret += "\n\tMaybe 'HasBeenLeader': yes"
	} else {
		ret += fmt.Sprintf("\n\tMaybe 'HasBeenLeader': yes? 0x%x", charArbServStatus.Player1PermHasBeenLeaderIsh)
	}

	if charArbServStatus.Player1PermLeader == 0 {
		ret += "\n\tLeader: no"
	} else if charArbServStatus.Player1PermLeader == 1 {
		ret += "\n\tLeader: yes"
	} else {
		ret += fmt.Sprintf("\n\tLeader: yes? 0x%x", charArbServStatus.Player1PermLeader)
	}

	if charArbServStatus.Player1PermModifyCharacters == 0 {
		ret += "\n\tModify Characters: no"
	} else if charArbServStatus.Player1PermModifyCharacters == 1 {
		ret += "\n\tModify Characters: yes"
	} else {
		ret += fmt.Sprintf("\n\tModify Characters: yes? 0x%x", charArbServStatus.Player1PermModifyCharacters)
	}

	ret += "\nUnk3: " + hex.EncodeToString(charArbServStatus.Unknown3[:])

	for k, v := range charArbServStatus.CharIsReady {
		if v == 0 {
			ret += "\nCharacter " + strconv.Itoa(k) + " is not ready"
		} else if v == 1 {
			ret += "\nCharacter " + strconv.Itoa(k) + " is ready"
		} else {
			ret += fmt.Sprintf("\nCharacter "+strconv.Itoa(k)+" is ready? 0x%x", v)
		}

	}

	ret += "\nUnk4: " + hex.EncodeToString(charArbServStatus.Unknown4[:])

	for k, v := range charArbServStatus.CharOwnerPlayerID {
		ret += fmt.Sprintf("\nCharacter "+strconv.Itoa(k)+" is controlled by %x", v)
	}
	ret += "\nImport Character: "
	if charArbServStatus.ImportCharSettings == 1 {
		ret += "Statistics"
	} else if charArbServStatus.ImportCharSettings == 3 {
		ret += "Statistics, Experience"
	} else if charArbServStatus.ImportCharSettings == 7 {
		ret += "Statistics, Experience, Items"
	} else {
		ret += fmt.Sprintf("Unknown Byte Value (0x%x)", charArbServStatus.ImportCharSettings)
	}
	ret += "\nRestrict Stores: "
	if charArbServStatus.RestrictStores == 0 {
		ret += "No"
	} else if charArbServStatus.RestrictStores == 1 {
		ret += "Yes"
	} else {
		ret += fmt.Sprintf("Unknown Byte Value (0x%x)", charArbServStatus.RestrictStores)
	}
	if charArbServStatus.ListenToJoinRequests == 0 {
		ret += "\nListenToJoinRequests: No"
	} else if charArbServStatus.ListenToJoinRequests == 1 {
		ret += "\nListenToJoinRequests: Yes"
	} else {
		ret += fmt.Sprintf("\nListenToJoinRequests: Yes? 0x%x", charArbServStatus.ListenToJoinRequests)
	}
	ret += "\nUnk5: " + hex.EncodeToString(charArbServStatus.Unknown5[:])
	return ret
}

const IECharArbServerStatusSize int = JMPacketHeaderSize + 398

//Not sure what to call this.
//The 43 byte  long packet with the server version
//Sent from server to client
//JoinLobby
//129 Host
type IEServerIntro1 struct {
	PlayerIDFromIsh uint32
	//					  JM			         v1.3.5521
	//01000000 00000000 000000ffff00642cc774 4a4d 00 01 00 13 ff 56 73 0309(76312e332e3535323100)1e000000
	//01000000 00000000 000000ffff00642cc774 4a4d 00 01 00 13 ff 56 73 0309(76312e332e3535323100)1e000000
}

//The 166-168 byte long packet after the server version
//Sent from server to client
//JoinLobby
//129 Host
type IEServerIntro2 struct {
	PlayerIDFrom uint32
	PlayerIDTo   uint32
	//0100000072442c00000001ffff01b32df6d14a4d0001003fff4d53000000c778da6364648000464630010550c122171d06983c0326a0b61844141b6667640c636160886042284b016236c72023330303060600c362037e000000000001000000004a4d0001003fff4d53000000c778da6364648000464630010550c122171d06983c0326a0b61844141b6667640c636160886042284b016236c72023330303060600c362037e
	//010000003a381100000001ffff01c53896604a4d00010040ff4d53000000c778da6364648000464630010550412b0b4106983c0326a0a618239cc486d91919c35818182298108a538098cd31c8c8ccc080810100852b03209caa6b0b9c016b0b00004a4d00010040ff4d53000000c778da6364648000464630010550412b0b4106983c0326a0a618239cc486d91919c35818182298108a538098cd31c8c8ccc080810100852b0320
}

//Not sure what to call this.
//Probably because I'm not sure what it does
//It's the packet that is 88-89 bytes long just before the pings start.
//JoinLobby
//129 Host
type IEPrePing struct {
	IEHeader
	//01d3e978 fb4a4d00 0100 40 ff4d53000000c778da6364648000464630010550 c122171d 06983c121b2ec68009 2811838862c3ec8c8c612c0c0c114ca89ad91c838ccc0c0c181800485b03fc
	//018268b2 974a4d00 0100 41 ff4d53000000c778da6364648000464630010550 412b0b41 06983c121b2ec68009 c815638493d8303b2363180b03430413aa6636c72023330303060600d014033f
}

//These seem to be ping types.
//The types start off by some mechinism I've yet to discover, but
//They seem to iterate as a sort of "ACK" when a message is received
//Only the receiver iterates the ping type, not the sender.
//01000000 7adb1f02 01000000 06 00 39 9b 88 23
//7adb1f02 01000000 01000000 07 00 f2 c7 5b 86
//7adb1f02 01000000 01000000 08 00 03 91 e9 53
//72442c00 01000000 01000000 09 00 c8 cd 3a f6
//72442c00 01000000 01000000 0a 00 4e 59 48 58
//3a381100 01000000 01000000 0b 00 85 05 9b fd
//                           0c 00 98 00 ab 45
//                           0d 00 53 5c 78 e0
//3a381100 01000000 01000000 0e 00 d5 c8 0a 4e
//3a381100 01000000 01000000 0f 00 1e 94 d9 eb
//3a381100 01000000 01000000 10 00 ec 14 69 a5
//ef6a8e02 01000000 01000000 11 00 27 48 ba 00
//ef6a8e02 01000000 01000000 12 00 a1 dc c8 ae
//ef6a8e02 01000000 01000000 13 00 6a 80 1b 0b
//ef6a8e02 01000000 01000000 14 00 77 85 2b b3
//ef6a8e02 01000000 01000000 15 00 bc d9 f8 16
//ef6a8e02 01000000 01000000 16 00 3a 4d 8a b8
//ef6a8e02 01000000 01000000 17 00 f1 11 59 1d
//ef6a8e02 01000000 01000000 18 00 00 47 eb c8
//ef6a8e02 01000000 01000000 19 00 cb 1b 38 6d
//ef6a8e02 01000000 01000000 1a 00 4d 8f 4a c3
//ef6a8e02 01000000 01000000 1b 00 86 d3 99 66
//ef6a8e02 01000000 01000000 1c 00 9b d6 a9 de
//ef6a8e02 01000000 01000000 1d 00 50 8a 7a 7b
//ef6a8e02 01000000 01000000 1e 00 d6 1e 08 d5
//ef6a8e02 01000000 01000000 1f 00 1d 42 db 70
//ef6a8e02 01000000 01000000 20 00 e8 6e 6e 08
//ef6a8e02 01000000 01000000 21 00 23 32 bd ad
//ef6a8e02 01000000 01000000 22 00 a5 a6 cf 03
//ef6a8e02 01000000 01000000 23 00 6e fa 1c a6
//ef6a8e02 01000000 01000000
//ef6a8e02 01000000 01000000 87 00 ec 77 4f 5e
//ef6a8e02 01000000 01000000 88 00 1d 21 fd 8b
type IEPing struct {
	IEHeader
	Unknown1           uint32
	PacketTypeIsh      uint16
	PacketSignatureIsh uint32
}

//??
//JoinLobby, Second Message via PlayerTestName: Hey I'm 129
//JoinLobby, First Message via PlayerTwoName: Hello There
type IEMsg struct {
	IEHeader
	//00 00 08 00
	//00 00 0a 00
	//00 00 08 00
	Unknown1 uint32
	//07 00
	//08 00
	//09 00
	PacketTypeIsh uint16
	//d1 91 b8 0c
	//e0 44 52 0d
	//96 50 65 78
	Unknown2 uint64
	//4a 4d 00 01
	//4a 4d 00 01
	//4a 4d 00 01
	Unknown3 uint32
	//00
	//00
	//00
	Unknown4 uint8
	//16
	//1f
	//1e
	MessageLengthPlusOne uint8
	//15
	//1e
	//1d
	MessageLength uint8
	//5b 50 6c 61 79 65 72 54 65 73 74 4e 61 6d 65 5d 3a 20 20 48 69
	//5b 50 6c 61 79 65 72 54 65 73 74 4e 61 6d 65 5d 3a 20 20 48 65 79 20 49 27 6d 20 31 32 39
	//5b 50 6c 61 79 65 72 54 77 6f 4e 61 6d 65 5d 3a 20 20 48 65 6c 6c 6f 20 74 68 65 72 65
	Message string
}

//============= 29/30 toggle before 18 len packets ==========
//JoinLobby, Send Two Messages.pcap
//129 Host, 131 Server.pcap
//===29
//72442c00 01000000 0000 0100 0100 e57afa7f 4a4d0001 00 05ff5044 0000
//3a381100 01000000 0000 0200 0300 8b0bdd50 4a4d0001 00 05ff5044 0100
//===30
//01000000 72442c00 0000 0300 0100 9b4bd578 4a4d0001 00 06ff5064 000000
//01000000 3a381100 0000 0400 0200 914bc4be 4a4d0001 00 06ff5064 010000
//===29
//72442c00 01000000 0000 0200 0300 8b0bdd50 4a4d0001 00 05ff5044 0100
//3a381100 01000000 0000 0300 0400 57567280 4a4d0001 00 05ff5044 0200
//===30
//01000000 72442c00 0000 0400 0200 914bc4be 4a4d0001 00 06ff5064 010000
//01000000 3a381100 0000 0500 0300 1938f740 4a4d0001 00 06ff5064 020000
//===29
//72442c00 01000000 0000 0300 0400 57567280 4a4d0001 00 05ff5044 0200
//3a381100 01000000 0000 0400 0500 83d503c9 4a4d0001 00 05ff5044 0300
//===30
//01000000 72442c00 0000 0500 0300 1938f740 4a4d0001 00 06ff5064 020000
//01000000 3a381100 0000 0600 0400 567ddad9 4a4d0001 00 06ff5064 030000
//===29
//72442c00 01000000 0000 0400 0500 83d503c9 4a4d0001 00 05ff5044 0300
//3a381100 01000000 0000 0500 0600 48ec4ed2 4a4d0001 00 05ff5044 0400
//===30
//01000000 72442c00 0000 0700 0500 d90741fb 4a4d0001 00 06ff5064 040000
//01000000 3a381100 0000 0700 0500 d90741fb 4a4d0001 00 06ff5064 040000
//===29
//72442c00 01000000 0000 0600 0700 75073279 4a4d0001 00 05ff5044 0500
//3a381100 01000000 0000 0600 0700 75073279 4a4d0001 00 05ff5044 0500
//===30
//0100000072442c000000080006004edc868f4a4d00010006ff5064050000
//010000003a3811000000080006004edc868f4a4d00010006ff5064050000
//===32
//72442c000100000000000700080015f5f4144a4d00010008ff4d7972442c0001
//3a38110001000000000007000800702c8e7b4a4d00010008ff4d793a38110001
//===88/89
//0100000072442c00000009000701d3e978fb4a4d00010040ff4d53000000c778da6364648000464630010550c122171d06983c121b2ec680092811838862c3ec8c8c612c0c0c114ca89ad91c838ccc0c0c181800485b03fc
//010000003a3811000000090007018268b2974a4d00010041ff4d53000000c778da6364648000464630010550412b0b4106983c121b2ec68009c815638493d8303b2363180b03430413aa6636c72023330303060600d014033f

//**************************************************************
//====================Begin the 18 length packets===============
//72442c00 01000000 01000000 0900 c8cd3af6
//3a381100 01000000 01000000 0900 c8cd3af6
