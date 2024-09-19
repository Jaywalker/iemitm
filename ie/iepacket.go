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

const IE_SPEC_MSG_TYPE_INTRO byte = 0x56
const IE_SPEC_MSG_SUBTYPE_INTRO byte = 0x73

type IEIntroHeader struct {
	JMPacketHeader
	Unknown3         uint8 // 03
	VersionStringLen uint8
}

const IEIntroHeaderSize int = JMPacketHeaderSize + 2

type IEIntro struct {
	IEIntroHeader
	VersionString string
	IEIntroFooter
}

func (ieSrvIntro IEIntro) String() string {
	return fmt.Sprintf("Client Version: %s Unk3: 0x%x - VersionStrLen: 0x%x - Unk4: 0x%x", ieSrvIntro.VersionString, ieSrvIntro.Unknown3, ieSrvIntro.VersionStringLen, ieSrvIntro.Unknown4)
}

type IEIntroFooter struct {
	VersionStringNull uint8  // 00
	Unknown4          uint32 // 1e000000
}

const IEIntroFooterSize int = 5

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

const IE_SPEC_MSG_TYPE_CHAR_ARBITRATION byte = 0x4d
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
