package ie

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
)

func specMsgSubtypeToString(msgSubType uint8) string {
	switch msgSubType {
	}
	return "UNKNOWN"
}

type IEHeader struct {
	PlayerIDFrom  uint32
	PlayerIDTo    uint32 //Server seems to be 0001 //
	FrameKind_    uint8
	FrameNum      uint16
	FrameExpected uint16
	Compressed    uint8 // This tells us if we're using Zlib compression (Best Compression flag)
	CRC32         uint32
}

const IEHeaderSize int = 18

func (header IEHeader) String() string {
	return fmt.Sprintf("IEHead PlayerFrom: 0x%x PlayerTo: 0x%x FrameKind: 0x%x FrameNumber: 0x%x FrameExpected: 0x%x Compressed?: 0x%x CRC32: 0x%x", header.PlayerIDFrom, header.PlayerIDTo, header.FrameKind_, header.FrameNum, header.FrameExpected, header.Compressed, header.CRC32)
}

type JMPacket interface {
	String() string
	FromPlayerID() uint32
	ToPlayerID() uint32
	FrameKind() uint8
	FrameNumber() uint16
	FrameNumberExpected() uint16
	CRC() uint32
	IsSpecMsg() bool
	SpecType() uint8
	SpecSubType() uint8
	IsCompressed() bool
	DecompressedSize() uint32
	PacketLength() uint16
	DataLength() int
	PacketData() []byte
}

func NewJMPacket(data []byte, size int) (JMPacket, error) {
	var jmHeader JMHeader
	if err := binary.Read(bytes.NewReader(data), binary.BigEndian, &jmHeader); err != nil {
		return nil, errors.New("ERROR: JMHeader binary.Read failed: " + err.Error())
	}

	if data[JMHeaderSize] == 0xff {
		if jmHeader.Compressed == 1 {
			var jmSpecHeaderCompressed JMSpecHeaderCompressed
			if err := binary.Read(bytes.NewReader(data), binary.BigEndian, &jmSpecHeaderCompressed); err != nil {
				return nil, errors.New("ERROR: JMSpecHeaderCompressed binary.Read failed: " + err.Error())
			}
			jmSpecCompressed := JMSpecCompressed{jmSpecHeaderCompressed, []byte{}}
			jmSpecCompressed.Data = append(data[JMSpecHeaderCompressedSize:size])
			return jmSpecCompressed, nil
		} else {
			var jmSpecHeader JMSpecHeader
			if err := binary.Read(bytes.NewReader(data), binary.BigEndian, &jmSpecHeader); err != nil {
				return nil, errors.New("ERROR: JMSpecHeader binary.Read failed: " + err.Error())
			}
			jmSpec := JMSpec{jmSpecHeader, []byte{}}
			if JMSpecHeaderSize > size {
				return nil, errors.New("ERROR: JMSpecHeaderSize > size")
			}
			jmSpec.Data = append(data[JMSpecHeaderSize:size])
			return jmSpec, nil
		}
	} else {
		if jmHeader.Compressed == 1 {
			var jmHeaderCompressed JMHeaderCompressed
			if err := binary.Read(bytes.NewReader(data), binary.BigEndian, &jmHeaderCompressed); err != nil {
				return nil, errors.New("ERROR: JMHeaderCompressed binary.Read failed: " + err.Error())
			}
			jmCompressed := JMCompressed{jmHeaderCompressed, []byte{}}
			jmCompressed.Data = append(data[JMHeaderCompressedSize:size])
			return jmCompressed, nil
		} else {
			jm := JM{jmHeader, []byte{}}
			jm.Data = append(data[JMHeaderSize:size])
			return jm, nil
		}
	}

	return nil, errors.New("ERROR: Not a known JM Packet type")
}

type JMHeader struct {
	IEHeader
	JM        [2]byte
	Unknown1  uint8 // 00
	Unknown2  uint8 // 01
	PacketLen uint16
}

const JMHeaderSize int = IEHeaderSize + 6

func (jmHeader JMHeader) String() string {
	return fmt.Sprintf("IEHead PlayerFrom: 0x%x PlayerTo: 0x%x FrameKind: 0x%x FrameNumber: 0x%x FrameExpected: 0x%x Compressed?: 0x%x CRC32: 0x%x - %c%c Unk1: 0x%x Unk2: 0x%x Len: %d", jmHeader.PlayerIDFrom, jmHeader.PlayerIDTo, jmHeader.FrameKind_, jmHeader.FrameNum, jmHeader.FrameExpected, jmHeader.Compressed, jmHeader.CRC32, jmHeader.JM[0], jmHeader.JM[1], jmHeader.Unknown1, jmHeader.Unknown2, jmHeader.PacketLen)
}

func (jmHeader JMHeader) FromPlayerID() uint32 {
	return jmHeader.PlayerIDFrom
}
func (jmHeader JMHeader) ToPlayerID() uint32 {
	return jmHeader.PlayerIDTo
}
func (jmHeader JMHeader) FrameKind() uint8 {
	return jmHeader.FrameKind_
}
func (jmHeader JMHeader) FrameNumber() uint16 {
	return jmHeader.FrameNum
}
func (jmHeader JMHeader) FrameNumberExpected() uint16 {
	return jmHeader.FrameExpected
}
func (jmHeader JMHeader) CRC() uint32 {
	return jmHeader.CRC32
}

func (jmHeader JMHeader) IsSpecMsg() bool {
	return false
}

func (jmHeader JMHeader) SpecType() uint8 {
	return 0
}

func (jmHeader JMHeader) SpecSubType() uint8 {
	return 0
}

func (jmHeader JMHeader) IsCompressed() bool {
	return false
}

func (jmHeader JMHeader) DecompressedSize() uint32 {
	return 0
}

func (jmHeader JMHeader) PacketLength() uint16 {
	return jmHeader.PacketLen
}

func (jmHeader JMHeader) DataLength() int {
	return 0
}

func (jmHeader JMHeader) PacketData() []byte {
	return nil
}

type JMHeaderCompressed struct {
	JMHeader
	DecompressedSize_ uint32
}

const JMHeaderCompressedSize int = JMHeaderSize + 4

type JMCompressed struct {
	JMHeaderCompressed
	Data []byte
}

func (jmCompressed JMCompressed) String() string {
	return fmt.Sprintf("IEHead PlayerFrom: 0x%x PlayerTo: 0x%x FrameKind: 0x%x FrameNumber: 0x%x FrameExpected: 0x%x Compressed?: 0x%x CRC32: 0x%x - %c%c Unk1: 0x%x Unk2: 0x%x Len: %d DecpmpressedSize: %d", jmCompressed.PlayerIDFrom, jmCompressed.PlayerIDTo, jmCompressed.FrameKind_, jmCompressed.FrameNum, jmCompressed.FrameExpected, jmCompressed.Compressed, jmCompressed.CRC32, jmCompressed.JM[0], jmCompressed.JM[1], jmCompressed.Unknown1, jmCompressed.Unknown2, jmCompressed.PacketLen, jmCompressed.DecompressedSize_)
}

func (jmCompressed JMCompressed) FromPlayerID() uint32 {
	return jmCompressed.PlayerIDFrom
}
func (jmCompressed JMCompressed) ToPlayerID() uint32 {
	return jmCompressed.PlayerIDTo
}
func (jmCompressed JMCompressed) FrameKind() uint8 {
	return jmCompressed.FrameKind_
}
func (jmCompressed JMCompressed) FrameNumber() uint16 {
	return jmCompressed.FrameNum
}
func (jmCompressed JMCompressed) FrameNumberExpected() uint16 {
	return jmCompressed.FrameExpected
}
func (jmCompressed JMCompressed) CRC() uint32 {
	return jmCompressed.CRC32
}

func (jmCompressed JMCompressed) IsSpecMsg() bool {
	return false
}

func (jmCompressed JMCompressed) SpecType() uint8 {
	return 0
}

func (jmCompressed JMCompressed) SpecSubType() uint8 {
	return 0
}

func (jmCompressed JMCompressed) IsCompressed() bool {
	return true
}

func (jmCompressed JMCompressed) DecompressedSize() uint32 {
	return jmCompressed.DecompressedSize_
}

func (jmCompressed JMCompressed) PacketLength() uint16 {
	return jmCompressed.PacketLen
}

func (jmCompressed JMCompressed) DataLength() int {
	return binary.Size(jmCompressed.Data)
}

func (jmCompressed JMCompressed) PacketData() []byte {
	return jmCompressed.Data
}

type JMSpecHeaderCompressed struct {
	JMHeader
	SpecMsgFlag       uint8
	SpecMsgType       uint8
	SpecMsgSubtype    uint8
	DecompressedSize_ uint32
}

const JMSpecHeaderCompressedSize int = JMHeaderSize + 7

type JMSpecCompressed struct {
	JMSpecHeaderCompressed
	Data []byte
}

func (jmSpecCompressed JMSpecCompressed) String() string {
	specType, specSubType := specMsgTypeToString(jmSpecCompressed.SpecMsgType, jmSpecCompressed.SpecMsgSubtype)
	return fmt.Sprintf("IEHead PlayerFrom: 0x%x PlayerTo: 0x%x FrameKind: 0x%x FrameNumber: 0x%x FrameExpected: 0x%x Compressed?: 0x%x CRC32: 0x%x - %c%c Unk1: 0x%x Unk2: 0x%x Len: %d SpecMsgFlag: 0x%x SpecMsgType: %s (%d) SpecMsgSubtype: %s (%d) DecompressedSize: 0x%x", jmSpecCompressed.PlayerIDFrom, jmSpecCompressed.PlayerIDTo, jmSpecCompressed.FrameKind_, jmSpecCompressed.FrameNum, jmSpecCompressed.FrameExpected, jmSpecCompressed.Compressed, jmSpecCompressed.CRC32, jmSpecCompressed.JM[0], jmSpecCompressed.JM[1], jmSpecCompressed.Unknown1, jmSpecCompressed.Unknown2, jmSpecCompressed.PacketLen, jmSpecCompressed.SpecMsgFlag, specType, jmSpecCompressed.SpecMsgType, specSubType, jmSpecCompressed.SpecMsgSubtype, jmSpecCompressed.DecompressedSize_)
}

func (jmSpecCompressed JMSpecCompressed) FromPlayerID() uint32 {
	return jmSpecCompressed.PlayerIDFrom
}
func (jmSpecCompressed JMSpecCompressed) ToPlayerID() uint32 {
	return jmSpecCompressed.PlayerIDTo
}
func (jmSpecCompressed JMSpecCompressed) FrameKind() uint8 {
	return jmSpecCompressed.FrameKind_
}
func (jmSpecCompressed JMSpecCompressed) FrameNumber() uint16 {
	return jmSpecCompressed.FrameNum
}
func (jmSpecCompressed JMSpecCompressed) FrameNumberExpected() uint16 {
	return jmSpecCompressed.FrameExpected
}
func (jmSpecCompressed JMSpecCompressed) CRC() uint32 {
	return jmSpecCompressed.CRC32
}

func (jmSpecCompressed JMSpecCompressed) IsSpecMsg() bool {
	return true
}

func (jmSpecCompressed JMSpecCompressed) SpecType() uint8 {
	return jmSpecCompressed.SpecMsgType
}
func (jmSpecCompressed JMSpecCompressed) SpecSubType() uint8 {
	return jmSpecCompressed.SpecMsgSubtype
}

func (jmSpecCompressed JMSpecCompressed) IsCompressed() bool {
	return true
}

func (jmSpecCompressed JMSpecCompressed) DecompressedSize() uint32 {
	return jmSpecCompressed.DecompressedSize_
}

func (jmSpecCompressed JMSpecCompressed) PacketLength() uint16 {
	return jmSpecCompressed.PacketLen
}

func (jmSpecCompressed JMSpecCompressed) DataLength() int {
	return binary.Size(jmSpecCompressed.Data)
}

func (jmSpecCompressed JMSpecCompressed) PacketData() []byte {
	return jmSpecCompressed.Data
}

type JMSpecHeader struct {
	JMHeader
	SpecMsgFlag    uint8
	SpecMsgType    uint8
	SpecMsgSubtype uint8
}

const JMSpecHeaderSize int = JMHeaderSize + 3

type JMSpec struct {
	JMSpecHeader
	Data []byte
}

func (jmSpec JMSpec) String() string {
	specType, specSubType := specMsgTypeToString(jmSpec.SpecMsgType, jmSpec.SpecMsgSubtype)
	return fmt.Sprintf("IEHead PlayerFrom: 0x%x PlayerTo: 0x%x FrameKind: 0x%x FrameNumber: 0x%x FrameExpected: 0x%x Compressed?: 0x%x CRC32: 0x%x - %c%c Unk1: 0x%x Unk2: 0x%x Len: %d SpecMsgFlag: 0x%x SpecMsgType: %s (%d) SpecMsgSubtype: %s (%d)", jmSpec.PlayerIDFrom, jmSpec.PlayerIDTo, jmSpec.FrameKind_, jmSpec.FrameNum, jmSpec.FrameExpected, jmSpec.Compressed, jmSpec.CRC32, jmSpec.JM[0], jmSpec.JM[1], jmSpec.Unknown1, jmSpec.Unknown2, jmSpec.PacketLen, jmSpec.SpecMsgFlag, specType, jmSpec.SpecMsgType, specSubType, jmSpec.SpecMsgSubtype)
}

func (jmSpec JMSpec) FromPlayerID() uint32 {
	return jmSpec.PlayerIDFrom
}
func (jmSpec JMSpec) ToPlayerID() uint32 {
	return jmSpec.PlayerIDTo
}
func (jmSpec JMSpec) FrameKind() uint8 {
	return jmSpec.FrameKind_
}
func (jmSpec JMSpec) FrameNumber() uint16 {
	return jmSpec.FrameNum
}
func (jmSpec JMSpec) FrameNumberExpected() uint16 {
	return jmSpec.FrameExpected
}
func (jmSpec JMSpec) CRC() uint32 {
	return jmSpec.CRC32
}

func (jmSpec JMSpec) IsSpecMsg() bool {
	return true
}

func (jmSpec JMSpec) SpecType() uint8 {
	return jmSpec.SpecMsgType
}
func (jmSpec JMSpec) SpecSubType() uint8 {
	return jmSpec.SpecMsgSubtype
}

func (jmSpec JMSpec) IsCompressed() bool {
	return false
}

func (jmSpec JMSpec) DecompressedSize() uint32 {
	return 0
}

func (jmSpec JMSpec) PacketLength() uint16 {
	return jmSpec.PacketLen
}

func (jmSpec JMSpec) DataLength() int {
	return binary.Size(jmSpec.Data)
}

func (jmSpec JMSpec) PacketData() []byte {
	return jmSpec.Data
}

type JM struct {
	JMHeader
	Data []byte
}

func (jm JM) String() string {
	return fmt.Sprintf("IEHead PlayerFrom: 0x%x PlayerTo: 0x%x FrameKind: 0x%x FrameNumber: 0x%x FrameExpected: 0x%x Compressed?: 0x%x CRC32: 0x%x - %c%c Unk1: 0x%x Unk2: 0x%x Len: %d", jm.PlayerIDFrom, jm.PlayerIDTo, jm.FrameKind_, jm.FrameNum, jm.FrameExpected, jm.Compressed, jm.CRC32, jm.JM[0], jm.JM[1], jm.Unknown1, jm.Unknown2, jm.PacketLen)
}

func (jm JM) FromPlayerID() uint32 {
	return jm.PlayerIDFrom
}
func (jm JM) ToPlayerID() uint32 {
	return jm.PlayerIDTo
}
func (jm JM) FrameKind() uint8 {
	return jm.FrameKind_
}
func (jm JM) FrameNumber() uint16 {
	return jm.FrameNum
}
func (jm JM) FrameNumberExpected() uint16 {
	return jm.FrameExpected
}
func (jm JM) CRC() uint32 {
	return jm.CRC32
}

func (jm JM) IsSpecMsg() bool {
	return false
}

func (jm JM) SpecType() uint8 {
	return 0
}
func (jm JM) SpecSubType() uint8 {
	return 0
}

func (jm JM) IsCompressed() bool {
	return false
}

func (jm JM) DecompressedSize() uint32 {
	return 0
}

func (jm JM) PacketLength() uint16 {
	return jm.PacketLen
}

func (jm JM) DataLength() int {
	return binary.Size(jm.Data)
}

func (jm JM) PacketData() []byte {
	return jm.Data
}

type IEMsg struct {
	MessageLength uint8
	Message       string
}

func (iemsg IEMsg) String() string {
	return iemsg.Message
}

type IEMsgDecompressed struct {
	JMHeaderCompressed
	IEMsg
}

type IEMsgPacket struct {
	JMHeader
	IEMsg
}

func (iemsg IEMsgPacket) String() string {
	return iemsg.Message
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

const IE_SPEC_MSG_TYPE_VERSION uint8 = 86
const IE_SPEC_MSG_SUBTYPE_VERSION_SERVER uint8 = 115

type IEVersionHeader struct {
	// JMHeader
	NumFields        uint8 // 03
	VersionStringLen uint8
}

const IEVersionHeaderSize int = 2

type IEVersion struct {
	IEVersionHeader
	VersionString string
	IEVersionFooter
}

func (ieVersion IEVersion) String() string {
	return fmt.Sprintf("Client Version: %s NumFields: 0x%x - VersionStrLen: 0x%x - Expansion: 0x%x TimerUpdatesPerSecond: %d", ieVersion.VersionString, ieVersion.NumFields, ieVersion.VersionStringLen, ieVersion.ExpansionPack, ieVersion.TimerUpdatesPerSecond)
}

type IEVersionFooter struct {
	ExpansionPack         uint8
	TimerUpdatesPerSecond uint32 // 1e000000
}

const IEVersionFooterSize int = 5

const IE_SPEC_MSG_TYPE_MPSETTINGS uint8 = 77
const IE_SPEC_MSG_SUBTYPE_TOGGLE_CHAR_READY uint8 = 114

type IEMPSettingsToggleCharReady struct {
	// JMSpec
	CharacterNum uint8
	ReadyStatus  uint32
}

const IEMPSettingsToggleCharReadySize int = 5

func (charReady IEMPSettingsToggleCharReady) String() string {
	if charReady.ReadyStatus == 0 {
		return "Character " + strconv.Itoa(int(charReady.CharacterNum)) + " is not ready."
	} else if charReady.ReadyStatus == 1 {
		return "Character " + strconv.Itoa(int(charReady.CharacterNum)) + " is ready."
	}
	return "Character " + strconv.Itoa(int(charReady.CharacterNum)) + " is ready (" + strconv.Itoa(int(charReady.ReadyStatus)) + ")."
}

const IE_SPEC_MSG_SUBTYPE_UPDATE_SERVER_ARBITRATION_INFO uint8 = 83

type IEMPSettingsFullSet struct {
	Unknown1                    [2]byte
	DefaultPermBuyAndSell       uint8
	DefaultPermTravel           uint8
	DefaultPermDialog           uint8
	DefaultPermViewCharacters   uint8
	DefaultPermPause            uint8
	DefaultPermHasBeenLeaderIsh uint8
	DefaultPermLeader           uint8
	DefaultPermModifyCharacters uint8
	Unknown2                    [20]byte
	Player1ID                   uint32
	Player1PermBuyAndSell       uint8
	Player1PermTravel           uint8
	Player1PermDialog           uint8
	Player1PermViewCharacters   uint8
	Player1PermPause            uint8
	Player1PermHasBeenLeaderIsh uint8
	Player1PermLeader           uint8
	Player1PermModifyCharacters uint8
	Unknown3                    [89]byte
	CharIsReady                 [6]uint8
	Unknown4                    [6]byte
	CharOwnerPlayerID           [6]uint32
	ImportCharSettings          uint8
	RestrictStores              uint8
	ListenToJoinRequests        uint8
	Unknown5                    [29]byte
}

const IEMPSettingsFullSetSize int = 199

func (charArbServStatus IEMPSettingsFullSet) String() string {
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

	ret += fmt.Sprintf("\nPlayer1Perms (0x%x):", charArbServStatus.Player1ID)

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
		ret += fmt.Sprintf("\nCharacter "+strconv.Itoa(k)+" is controlled by 0x%x", v)
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

func specMsgTypeToString(msgType, msgSubType uint8) (string, string) {
	switch msgType {
	case 66:
		msgTypeStr := "PROGRESSBAR"
		switch msgSubType {
		case 83:
			msgSubTypeStr := "PROGRESSBAR_STATUS"
			return msgTypeStr, msgSubTypeStr
		}
	case 67:
		msgTypeStr := "CMESSAGE"
		switch msgSubType {
		case 0:
			msgSubTypeStr := "CMESSAGE_ADD_ACTION"
			return msgTypeStr, msgSubTypeStr
		case 1:
			msgSubTypeStr := "CMESSAGE_ADD_EFFECT"
			return msgTypeStr, msgSubTypeStr
		case 3:
			msgSubTypeStr := "CMESSAGE_ANIMATION_CHANGE"
			return msgTypeStr, msgSubTypeStr
		case 4:
			msgSubTypeStr := "CMESSAGE_CHANGE_DIRECTION"
			return msgTypeStr, msgSubTypeStr
		case 5:
			msgSubTypeStr := "CMESSAGE_CLEAR_ACTIONS"
			return msgTypeStr, msgSubTypeStr
		case 6:
			msgSubTypeStr := "CMESSAGE_CLEAR_DIALOG_ACTIONS"
			return msgTypeStr, msgSubTypeStr
		case 7:
			msgSubTypeStr := "CMESSAGE_CLEAR_GROUP_SLOT"
			return msgTypeStr, msgSubTypeStr
		case 8:
			msgSubTypeStr := "CMESSAGE_CLEAR_TRIGGERS"
			return msgTypeStr, msgSubTypeStr
		case 9:
			msgSubTypeStr := "CMESSAGE_COLOR_CHANGE"
			return msgTypeStr, msgSubTypeStr
		case 10:
			msgSubTypeStr := "CMESSAGE_COLOR_RESET"
			return msgTypeStr, msgSubTypeStr
		case 11:
			msgSubTypeStr := "CMESSAGE_COLOR_UPDATE"
			return msgTypeStr, msgSubTypeStr
		case 12:
			msgSubTypeStr := "CMESSAGE_CONTAINER_ADD_ITEM"
			return msgTypeStr, msgSubTypeStr
		case 13:
			msgSubTypeStr := "CMESSAGE_CONTAINER_ITEMS"
			return msgTypeStr, msgSubTypeStr
		case 14:
			msgSubTypeStr := "CMESSAGE_CONTAINER_STATUS"
			return msgTypeStr, msgSubTypeStr
		case 15:
			msgSubTypeStr := "CMESSAGE_CUT_SCENE_MODE_STATUS"
			return msgTypeStr, msgSubTypeStr
		case 16:
			msgSubTypeStr := "CMESSAGE_DISPLAY_TEXT"
			return msgTypeStr, msgSubTypeStr
		case 17:
			msgSubTypeStr := "CMESSAGE_DISPLAY_TEXTREF"
			return msgTypeStr, msgSubTypeStr
		case 19:
			msgSubTypeStr := "CMESSAGE_DOOR_STATUS"
			return msgTypeStr, msgSubTypeStr
		case 20:
			msgSubTypeStr := "CMESSAGE_DROP_PATH"
			return msgTypeStr, msgSubTypeStr
		case 21:
			msgSubTypeStr := "CMESSAGE_ENTER_DIALOG"
			return msgTypeStr, msgSubTypeStr
		case 23:
			msgSubTypeStr := "CMESSAGE_ENTER_STORE_MODE"
			return msgTypeStr, msgSubTypeStr
		case 24:
			msgSubTypeStr := "CMESSAGE_EXIT_DIALOG_MODE"
			return msgTypeStr, msgSubTypeStr
		case 25:
			msgSubTypeStr := "CMESSAGE_EXIT_STORE_MODE"
			return msgTypeStr, msgSubTypeStr
		case 26:
			msgSubTypeStr := "CMESSAGE_FIRE_PROJECTILE"
			return msgTypeStr, msgSubTypeStr
		case 27:
			msgSubTypeStr := "CMESSAGE_INSERT_ACTION"
			return msgTypeStr, msgSubTypeStr
		case 29:
			msgSubTypeStr := "CMESSAGE_LEAVE_PARTY"
			return msgTypeStr, msgSubTypeStr
		case 31:
			msgSubTypeStr := "CMESSAGE_PARTY_GOLD"
			return msgTypeStr, msgSubTypeStr
		case 32:
			msgSubTypeStr := "CMESSAGE_PLAY_SOUND"
			return msgTypeStr, msgSubTypeStr
		case 33:
			msgSubTypeStr := "CMESSAGE_PLAY_SOUND_REF"
			return msgTypeStr, msgSubTypeStr
		case 35:
			msgSubTypeStr := "CMESSAGE_REMOVE_REPLIES"
			return msgTypeStr, msgSubTypeStr
		case 36:
			msgSubTypeStr := "CMESSAGE_REPUTATION_CHANGE"
			return msgTypeStr, msgSubTypeStr
		case 37:
			msgSubTypeStr := "CMESSAGE_SET_ACTIVE"
			return msgTypeStr, msgSubTypeStr
		case 38:
			msgSubTypeStr := "CMESSAGE_SET_AISPEED"
			return msgTypeStr, msgSubTypeStr
		case 39:
			msgSubTypeStr := "CMESSAGE_SET_COMMAND_PAUSE"
			return msgTypeStr, msgSubTypeStr
		case 40:
			msgSubTypeStr := "CMESSAGE_SET_DIALOG_WAIT"
			return msgTypeStr, msgSubTypeStr
		case 41:
			msgSubTypeStr := "CMESSAGE_SET_DIRECTION"
			return msgTypeStr, msgSubTypeStr
		case 42:
			msgSubTypeStr := "CMESSAGE_SET_DRAW_POLY"
			return msgTypeStr, msgSubTypeStr
		case 43:
			msgSubTypeStr := "CMESSAGE_SET_FORCE_ACTION_PICK"
			return msgTypeStr, msgSubTypeStr
		case 44:
			msgSubTypeStr := "CMESSAGE_SET_HAPPINESS"
			return msgTypeStr, msgSubTypeStr
		case 45:
			msgSubTypeStr := "CMESSAGE_SET_IN_CUT_SCENE"
			return msgTypeStr, msgSubTypeStr
		case 46:
			msgSubTypeStr := "CMESSAGE_SET_LAST_ATTACKER"
			return msgTypeStr, msgSubTypeStr
		case 48:
			msgSubTypeStr := "CMESSAGE_SET_NUM_TIMES_TALKED_TO"
			return msgTypeStr, msgSubTypeStr
		case 49:
			msgSubTypeStr := "CMESSAGE_SET_PATH"
			return msgTypeStr, msgSubTypeStr
		case 50:
			msgSubTypeStr := "CMESSAGE_SET_SEQUENCE"
			return msgTypeStr, msgSubTypeStr
		case 52:
			msgSubTypeStr := "CMESSAGE_SET_TRIGGER"
			return msgTypeStr, msgSubTypeStr
		case 47:
			msgSubTypeStr := "CMESSAGE_SET_LAST_OBJECT"
			return msgTypeStr, msgSubTypeStr
		case 54:
			msgSubTypeStr := "CMESSAGE_SPRITE_DEATH"
			return msgTypeStr, msgSubTypeStr
		case 55:
			msgSubTypeStr := "CMESSAGE_SPRITE_EQUIPMENT"
			return msgTypeStr, msgSubTypeStr
		case 56:
			msgSubTypeStr := "CMESSAGE_SPRITE_PETRIFY"
			return msgTypeStr, msgSubTypeStr
		case 57:
			msgSubTypeStr := "CMESSAGE_SPRITE_UPDATE"
			return msgTypeStr, msgSubTypeStr
		case 58:
			msgSubTypeStr := "CMESSAGE_START_FOLLOW"
			return msgTypeStr, msgSubTypeStr
		case 59:
			msgSubTypeStr := "CMESSAGE_START_SCROLL"
			return msgTypeStr, msgSubTypeStr
		case 60:
			msgSubTypeStr := "CMESSAGE_STOP_ACTIONS"
			return msgTypeStr, msgSubTypeStr
		case 61:
			msgSubTypeStr := "CMESSAGE_STOP_FOLLOW"
			return msgTypeStr, msgSubTypeStr
		case 62:
			msgSubTypeStr := "CMESSAGE_TRIGGER_STATUS"
			return msgTypeStr, msgSubTypeStr
		case 63:
			msgSubTypeStr := "CMESSAGE_UNLOCK"
			return msgTypeStr, msgSubTypeStr
		case 64:
			msgSubTypeStr := "CMESSAGE_UPDATE_REACTION"
			return msgTypeStr, msgSubTypeStr
		case 65:
			msgSubTypeStr := "CMESSAGE_VERBAL_CONSTANT"
			return msgTypeStr, msgSubTypeStr
		case 66:
			msgSubTypeStr := "CMESSAGE_VISIBILITY_MAP_MOVE"
			return msgTypeStr, msgSubTypeStr
		case 67:
			msgSubTypeStr := "CMESSAGE_VISUAL_EFFECT"
			return msgTypeStr, msgSubTypeStr
		case 68:
			msgSubTypeStr := "CMESSAGE_SET_DIALOG_RESREF"
			return msgTypeStr, msgSubTypeStr
		case 69:
			msgSubTypeStr := "CMESSAGE_ESCAPE_AREA"
			return msgTypeStr, msgSubTypeStr
		case 70:
			msgSubTypeStr := "CMESSAGE_DISPLAY_TEXTREF_SEND"
			return msgTypeStr, msgSubTypeStr
		case 71:
			msgSubTypeStr := "CMESSAGE_SET_CURRENT_ACTION_ID"
			return msgTypeStr, msgSubTypeStr
		case 72:
			msgSubTypeStr := "CMESSAGE_MOVE_GLOBAL"
			return msgTypeStr, msgSubTypeStr
		case 73:
			msgSubTypeStr := "CMESSAGE_FADE_COLOR"
			return msgTypeStr, msgSubTypeStr
		case 74:
			msgSubTypeStr := "CMESSAGE_START_TEXT_SCREEN"
			return msgTypeStr, msgSubTypeStr
		case 75:
			msgSubTypeStr := "CMESSAGE_SPAWNPT_ACTIVATE"
			return msgTypeStr, msgSubTypeStr
		case 76:
			msgSubTypeStr := "CMESSAGE_SPAWNPT_SPAWN"
			return msgTypeStr, msgSubTypeStr
		case 77:
			msgSubTypeStr := "CMESSAGE_STATIC_START"
			return msgTypeStr, msgSubTypeStr
		case 78:
			msgSubTypeStr := "CMESSAGE_STORE_ADD_ITEM"
			return msgTypeStr, msgSubTypeStr
		case 79:
			msgSubTypeStr := "CMESSAGE_STORE_REMOVE_ITEM"
			return msgTypeStr, msgSubTypeStr
		case 80:
			msgSubTypeStr := "CMESSAGE_FAMILIAR_ADD"
			return msgTypeStr, msgSubTypeStr
		case 81:
			msgSubTypeStr := "CMESSAGE_FAMILIAR_REMOVE_RESREF"
			return msgTypeStr, msgSubTypeStr
		case 82:
			msgSubTypeStr := "CMESSAGE_STOP_ESCAPE_AREA"
			return msgTypeStr, msgSubTypeStr
		case 85:
			msgSubTypeStr := "CMESSAGE_SET_TIME_STOP"
			return msgTypeStr, msgSubTypeStr
		case 87:
			msgSubTypeStr := "CMESSAGE_STORE_RELEASE"
			return msgTypeStr, msgSubTypeStr
		case 90:
			msgSubTypeStr := "CMESSAGE_90"
			return msgTypeStr, msgSubTypeStr
		case 91:
			msgSubTypeStr := "CMESSAGE_FLOAT_TEXT"
			return msgTypeStr, msgSubTypeStr
		case 92:
			msgSubTypeStr := "CMESSAGE_92"
			return msgTypeStr, msgSubTypeStr
		case 93:
			msgSubTypeStr := "CMESSAGE_SET_PROTAGONIST"
			return msgTypeStr, msgSubTypeStr
		case 94:
			msgSubTypeStr := "CMESSAGE_START_COMBAT_MUSIC"
			return msgTypeStr, msgSubTypeStr
		case 99:
			msgSubTypeStr := "CMESSAGE_SCREENSHAKE"
			return msgTypeStr, msgSubTypeStr
		case 100:
			msgSubTypeStr := "CMESSAGE_STORE_DEMAND"
			return msgTypeStr, msgSubTypeStr
		case 101:
			msgSubTypeStr := "CMESSAGE_101"
			return msgTypeStr, msgSubTypeStr
		case 102:
			msgSubTypeStr := "CMESSAGE_WEAPON_IMMUNITIES_UPDATE"
			return msgTypeStr, msgSubTypeStr
		case 103:
			msgSubTypeStr := "CMESSAGE_103"
			return msgTypeStr, msgSubTypeStr
		case 106:
			msgSubTypeStr := "CMESSAGE_TOGGLE_INTERFACE"
			return msgTypeStr, msgSubTypeStr
		case 107:
			msgSubTypeStr := "CMESSAGE_107"
			return msgTypeStr, msgSubTypeStr
		case 109:
			msgSubTypeStr := "CMESSAGE_SET_AREA_TYPE"
			return msgTypeStr, msgSubTypeStr
		case 110:
			msgSubTypeStr := "CMESSAGE_SET_AREA_REST_ENCOUNTER"
			return msgTypeStr, msgSubTypeStr
		case 114:
			msgSubTypeStr := "CMESSAGE_SET_AREA_EXPLORED"
			return msgTypeStr, msgSubTypeStr
		case 116:
			msgSubTypeStr := "CMESSAGE_CHANGE_STAT"
			return msgTypeStr, msgSubTypeStr
		case 120:
			msgSubTypeStr := "CMESSAGE_120"
			return msgTypeStr, msgSubTypeStr
		case 121:
			msgSubTypeStr := "CMESSAGE_END_GAME"
			return msgTypeStr, msgSubTypeStr
		}
	case 68:
		msgTypeStr := "DIALOG"
		switch msgSubType {
		case 82:
			msgSubTypeStr := "DIALOG_PERMIT_REQUEST"
			return msgTypeStr, msgSubTypeStr
		case 114:
			msgSubTypeStr := "DIALOG_PERMIT_REPLY"
			return msgTypeStr, msgSubTypeStr
		case 67:
			msgSubTypeStr := "DIALOG_CANCEL_REQUEST"
			return msgTypeStr, msgSubTypeStr
		case 75:
			msgSubTypeStr := "DIALOG_KILL_OR_USE"
			return msgTypeStr, msgSubTypeStr
		}
	case 73:
		msgTypeStr := "SWAPITEM"
		switch msgSubType {
		case 82:
			msgSubTypeStr := "SWAPITEM_REQUEST"
			return msgTypeStr, msgSubTypeStr
		case 114:
			msgSubTypeStr := "SWAPITEM_REPLY"
			return msgTypeStr, msgSubTypeStr
		}
	case 106:
		msgTypeStr := "JOURNAL"
		switch msgSubType {
		case 69:
			msgSubTypeStr := "JOURNAL_ADD_ENTRY"
			return msgTypeStr, msgSubTypeStr
		case 65:
			msgSubTypeStr := "JOURNAL_ANNOUNCE"
			return msgTypeStr, msgSubTypeStr
		case 117:
			msgSubTypeStr := "JOURNAL_ADD_USER_ENTRY"
			return msgTypeStr, msgSubTypeStr
		case 85:
			msgSubTypeStr := "JOURNAL_ANNOUNCE_USER_ENTRY"
			return msgTypeStr, msgSubTypeStr
		case 99:
			msgSubTypeStr := "JOURNAL_CHANGE_ENTRY"
			return msgTypeStr, msgSubTypeStr
		case 67:
			msgSubTypeStr := "JOURNAL_ANNOUNCE_CHANGE"
			return msgTypeStr, msgSubTypeStr
		}
	case 98:
		msgTypeStr := "BIOGRAPHY"
		switch msgSubType {
		case 99:
			msgSubTypeStr := "BIOGRAPHY_CHANGE"
			return msgTypeStr, msgSubTypeStr
		case 67:
			msgSubTypeStr := "BIOGRAPHY_CHANGE_ANNOUNCE"
			return msgTypeStr, msgSubTypeStr
		}
	case 75:
		msgTypeStr := "KICK_PLAYER"
		switch msgSubType {
		case 82:
			msgSubTypeStr := "KICK_PLAYER_REQUEST"
			return msgTypeStr, msgSubTypeStr
		case 83:
			msgSubTypeStr := "KICK_PLAYER_SERVER_SUPPORT"
			return msgTypeStr, msgSubTypeStr
		case 72:
			msgSubTypeStr := "KICK_PLAYER_HOOFED_OUT"
			return msgTypeStr, msgSubTypeStr
		}
	case 77:
		msgTypeStr := "MPSETTINGS"
		switch msgSubType {
		case 68:
			msgSubTypeStr := "MPSETTINGS_FULLDEMAND"
			return msgTypeStr, msgSubTypeStr
		case 83:
			msgSubTypeStr := "MPSETTINGS_FULLSET"
			return msgTypeStr, msgSubTypeStr
		case 112:
			msgSubTypeStr := "MPSETTINGS_PERMISSION"
			return msgTypeStr, msgSubTypeStr
		case 121:
			msgSubTypeStr := "MPSETTINGS_PLAYER_READY"
			return msgTypeStr, msgSubTypeStr
		case 114:
			msgSubTypeStr := "MPSETTINGS_CHAR_READY"
			return msgTypeStr, msgSubTypeStr
		case 99:
			msgSubTypeStr := "MPSETTINGS_CHAR_CONTROL"
			return msgTypeStr, msgSubTypeStr
		case 105:
			msgSubTypeStr := "MPSETTINGS_IMPORTING"
			return msgTypeStr, msgSubTypeStr
		case 106:
			msgSubTypeStr := "MPSETTINGS_LISTEN_JOIN"
			return msgTypeStr, msgSubTypeStr
		case 115:
			msgSubTypeStr := "MPSETTINGS_SLOT_STATUS"
			return msgTypeStr, msgSubTypeStr
		case 108:
			msgSubTypeStr := "MPSETTINGS_LOCK_STATUS"
			return msgTypeStr, msgSubTypeStr
		case 97:
			msgSubTypeStr := "MPSETTINGS_LOCK_ALLOW_INPUT"
			return msgTypeStr, msgSubTypeStr
		case 76:
			msgSubTypeStr := "MPSETTINGS_LOCK_REQUEST"
			return msgTypeStr, msgSubTypeStr
		case 54:
			msgSubTypeStr := "MPSETTINGS_GORE_LEVEL"
			return msgTypeStr, msgSubTypeStr
		case 57:
			msgSubTypeStr := "MPSETTINGS_RESTRICT_STORE"
			return msgTypeStr, msgSubTypeStr
		case 78:
			msgSubTypeStr := "MPSETTINGS_NIGHTMAREMODE"
			return msgTypeStr, msgSubTypeStr
		case 110:
			msgSubTypeStr := "MPSETTINGS_DEMAND_NIGHTMAREMODE"
			return msgTypeStr, msgSubTypeStr
		}
	case 79:
		msgTypeStr := "OBJECT"
		switch msgSubType {
		case 65:
			msgSubTypeStr := "OBJECT_ADD"
			return msgTypeStr, msgSubTypeStr
		}
	case 80:
		msgTypeStr := "PLAYERCHAR"
		switch msgSubType {
		case 85:
			msgSubTypeStr := "PLAYERCHAR_UPDATE_DEMAND"
			return msgTypeStr, msgSubTypeStr
		case 117:
			msgSubTypeStr := "PLAYERCHAR_UPDATE_REPLY"
			return msgTypeStr, msgSubTypeStr
		case 68:
			msgSubTypeStr := "PLAYERCHAR_DEMAND_SLOT"
			return msgTypeStr, msgSubTypeStr
		case 100:
			msgSubTypeStr := "PLAYERCHAR_DEMAND_REPLY"
			return msgTypeStr, msgSubTypeStr
		case 70:
			msgSubTypeStr := "PLAYERCHAR_70"
			return msgTypeStr, msgSubTypeStr
		case 102:
			msgSubTypeStr := "PLAYERCHAR_102"
			return msgTypeStr, msgSubTypeStr
		}
	case 81:
		msgTypeStr := "PAUSING"
		switch msgSubType {
		case 82:
			msgSubTypeStr := "PAUSING_PERMIT_REQUEST"
			return msgTypeStr, msgSubTypeStr
		case 65:
			msgSubTypeStr := "PAUSING_ANNOUNCE"
			return msgTypeStr, msgSubTypeStr
		}
	case 82:
		msgTypeStr := "RESOURCE"
		switch msgSubType {
		case 68:
			msgSubTypeStr := "RESOURCE_DEMAND"
			return msgTypeStr, msgSubTypeStr
		}
	case 83:
		msgTypeStr := "SIGNAL"
		switch msgSubType {
		case 83:
			msgSubTypeStr := "SIGNAL"
			return msgTypeStr, msgSubTypeStr
		case 82:
			msgSubTypeStr := "SIGNAL_REQUEST"
			return msgTypeStr, msgSubTypeStr
		}
	case 115:
		msgTypeStr := "MPSYNCH"
		switch msgSubType {
		case 82:
			msgSubTypeStr := "MPSYNCH_REQUEST"
			return msgTypeStr, msgSubTypeStr
		case 80:
			msgSubTypeStr := "MPSYNCH_REPLY"
			return msgTypeStr, msgSubTypeStr
		}
	case 86:
		msgTypeStr := "VERSION"
		switch msgSubType {
		case 115:
			msgSubTypeStr := "VERSION_SERVER"
			return msgTypeStr, msgSubTypeStr
		}
	case 88:
		msgTypeStr := "LEAVEAREALUA"
		switch msgSubType {
		case 82:
			msgSubTypeStr := "LEAVEAREALUA_PERMIT_REQUEST"
			return msgTypeStr, msgSubTypeStr
		}
	case 120:
		msgTypeStr := "LEAVEAREANAME"
		switch msgSubType {
		case 82:
			msgSubTypeStr := "LEAVEAREANAME_PERMIT_REQUEST"
			return msgTypeStr, msgSubTypeStr
		}
	}
	return "UNKNOWN", "UNKNOWN"
}
