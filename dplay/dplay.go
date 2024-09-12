package dplay

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"unicode/utf16"
	"unicode/utf8"
)

type DPPacketType uint16

const (
	DPSP_MSG_TYPE_ENUMSESSIONSREPLY DPPacketType = iota + 1
	DPSP_MSG_TYPE_ENUMSESSIONS
	DPSP_MSG_TYPE_ENUMPLAYERSREPLY
	DPSP_MSG_TYPE_ENUMPLAYER
	DPSP_MSG_TYPE_REQUESTPLAYERID
	DPSP_MSG_TYPE_REQUESTGROUPID
	DPSP_MSG_TYPE_REQUESTPLAYERREPLY
	DPSP_MSG_TYPE_CREATEPLAYER
	DPSP_MSG_TYPE_CREATEGROUP
	DPSP_MSG_TYPE_PLAYERMESSAGE
	DPSP_MSG_TYPE_DELETEPLAYER
	DPSP_MSG_TYPE_DELETEGROUP
	DPSP_MSG_TYPE_ADDPLAYERTOGROUP
	DPSP_MSG_TYPE_DELETEPLAYERFROMGROUP
	DPSP_MSG_TYPE_PLAYERDATACHANGED
	DPSP_MSG_TYPE_PLAYERNAMECHANGED
	DPSP_MSG_TYPE_GROUPDATACHANGED
	DPSP_MSG_TYPE_GROUPNAMECHANGED
	DPSP_MSG_TYPE_ADDFORWARDREQUEST
	DPSP_MSG_TYPE_UNK1
	DPSP_MSG_TYPE_PACKET
	DPSP_MSG_TYPE_PING
	DPSP_MSG_TYPE_PINGREPLY
	DPSP_MSG_TYPE_YOUAREDEAD
	DPSP_MSG_TYPE_PLAYERWRAPPER
	DPSP_MSG_TYPE_SESSIONDESCCHANGED
	DPSP_MSG_TYPE_UNK2
	DPSP_MSG_TYPE_CHALLENGE
	DPSP_MSG_TYPE_ACCESSGRANTED
	DPSP_MSG_TYPE_LOGONDENIED
	DPSP_MSG_TYPE_AUTHERROR
	DPSP_MSG_TYPE_NEGOTIATE
	DPSP_MSG_TYPE_CHALLENGERESPONSE
	DPSP_MSG_TYPE_SIGNED
	DPSP_MSG_TYPE_UNK3
	DPSP_MSG_TYPE_ADDFORWARDREPLY
	DPSP_MSG_TYPE_ASK4MULTICAST
	DPSP_MSG_TYPE_ASK4MULTICASTGUARANTEED
	DPSP_MSG_TYPE_ADDSHORTCUTTOGROUP
	DPSP_MSG_TYPE_DELETEGROUPFROMGROUP
	DPSP_MSG_TYPE_SUPERENUMPLAYERSREPLY
	DPSP_MSG_TYPE_UNK4
	DPSP_MSG_TYPE_KEYEXCHANGE
	DPSP_MSG_TYPE_KEYEXCHANGEREPLY
	DPSP_MSG_TYPE_CHAT
	DPSP_MSG_TYPE_ADDFORWARD
	DPSP_MSG_TYPE_ADDFORWARDACK
	DPSP_MSG_TYPE_PACKET2_DATA
	DPSP_MSG_TYPE_PACKET2_ACK
	DPSP_MSG_TYPE_UNK5
	DPSP_MSG_TYPE_UNK6
	DPSP_MSG_TYPE_UNK7
	DPSP_MSG_TYPE_IAMNAMESERVER
	DPSP_MSG_TYPE_VOICE
	DPSP_MSG_TYPE_MULTICASTDELIVERY
	DPSP_MSG_TYPE_CREATEPLAYERVERIFY
)

func commandToString(cmd DPPacketType) string {
	switch cmd {
	case DPSP_MSG_TYPE_ENUMSESSIONSREPLY:
		return "DPSP_MSG_TYPE_ENUMSESSIONSREPLY"
	case DPSP_MSG_TYPE_ENUMSESSIONS:
		return "DPSP_MSG_TYPE_ENUMSESSIONS"
	case DPSP_MSG_TYPE_ENUMPLAYERSREPLY:
		return "DPSP_MSG_TYPE_ENUMPLAYERSREPLY"
	case DPSP_MSG_TYPE_ENUMPLAYER:
		return "DPSP_MSG_TYPE_ENUMPLAYER"
	case DPSP_MSG_TYPE_REQUESTPLAYERID:
		return "DPSP_MSG_TYPE_REQUESTPLAYERID"
	case DPSP_MSG_TYPE_REQUESTGROUPID:
		return "DPSP_MSG_TYPE_REQUESTGROUPID"
	case DPSP_MSG_TYPE_REQUESTPLAYERREPLY:
		return "DPSP_MSG_TYPE_REQUESTPLAYERREPLY"
	case DPSP_MSG_TYPE_CREATEPLAYER:
		return "DPSP_MSG_TYPE_CREATEPLAYER"
	case DPSP_MSG_TYPE_CREATEGROUP:
		return "DPSP_MSG_TYPE_CREATEGROUP"
	case DPSP_MSG_TYPE_PLAYERMESSAGE:
		return "DPSP_MSG_TYPE_PLAYERMESSAGE"
	case DPSP_MSG_TYPE_DELETEPLAYER:
		return "DPSP_MSG_TYPE_DELETEPLAYER"
	case DPSP_MSG_TYPE_DELETEGROUP:
		return "DPSP_MSG_TYPE_DELETEGROUP"
	case DPSP_MSG_TYPE_ADDPLAYERTOGROUP:
		return "DPSP_MSG_TYPE_ADDPLAYERTOGROUP"
	case DPSP_MSG_TYPE_DELETEPLAYERFROMGROUP:
		return "DPSP_MSG_TYPE_DELETEPLAYERFROMGROUP"
	case DPSP_MSG_TYPE_PLAYERDATACHANGED:
		return "DPSP_MSG_TYPE_PLAYERDATACHANGED"
	case DPSP_MSG_TYPE_PLAYERNAMECHANGED:
		return "DPSP_MSG_TYPE_PLAYERNAMECHANGED"
	case DPSP_MSG_TYPE_GROUPDATACHANGED:
		return "DPSP_MSG_TYPE_GROUPDATACHANGED"
	case DPSP_MSG_TYPE_GROUPNAMECHANGED:
		return "DPSP_MSG_TYPE_GROUPNAMECHANGED"
	case DPSP_MSG_TYPE_ADDFORWARDREQUEST:
		return "DPSP_MSG_TYPE_ADDFORWARDREQUEST"
	case DPSP_MSG_TYPE_UNK1:
		return "DPSP_MSG_TYPE_UNK1"
	case DPSP_MSG_TYPE_PACKET:
		return "DPSP_MSG_TYPE_PACKET"
	case DPSP_MSG_TYPE_PING:
		return "DPSP_MSG_TYPE_PING"
	case DPSP_MSG_TYPE_PINGREPLY:
		return "DPSP_MSG_TYPE_PINGREPLY"
	case DPSP_MSG_TYPE_YOUAREDEAD:
		return "DPSP_MSG_TYPE_YOUAREDEAD"
	case DPSP_MSG_TYPE_PLAYERWRAPPER:
		return "DPSP_MSG_TYPE_PLAYERWRAPPER"
	case DPSP_MSG_TYPE_SESSIONDESCCHANGED:
		return "DPSP_MSG_TYPE_SESSIONDESCCHANGED"
	case DPSP_MSG_TYPE_UNK2:
		return "DPSP_MSG_TYPE_UNK2"
	case DPSP_MSG_TYPE_CHALLENGE:
		return "DPSP_MSG_TYPE_CHALLENGE"
	case DPSP_MSG_TYPE_ACCESSGRANTED:
		return "DPSP_MSG_TYPE_ACCESSGRANTED"
	case DPSP_MSG_TYPE_LOGONDENIED:
		return "DPSP_MSG_TYPE_LOGONDENIED"
	case DPSP_MSG_TYPE_AUTHERROR:
		return "DPSP_MSG_TYPE_AUTHERROR"
	case DPSP_MSG_TYPE_NEGOTIATE:
		return "DPSP_MSG_TYPE_NEGOTIATE"
	case DPSP_MSG_TYPE_CHALLENGERESPONSE:
		return "DPSP_MSG_TYPE_CHALLENGERESPONSE"
	case DPSP_MSG_TYPE_SIGNED:
		return "DPSP_MSG_TYPE_SIGNED"
	case DPSP_MSG_TYPE_UNK3:
		return "DPSP_MSG_TYPE_UNK3"
	case DPSP_MSG_TYPE_ADDFORWARDREPLY:
		return "DPSP_MSG_TYPE_ADDFORWARDREPLY"
	case DPSP_MSG_TYPE_ASK4MULTICAST:
		return "DPSP_MSG_TYPE_ASK4MULTICAST"
	case DPSP_MSG_TYPE_ASK4MULTICASTGUARANTEED:
		return "DPSP_MSG_TYPE_ASK4MULTICASTGUARANTEED"
	case DPSP_MSG_TYPE_ADDSHORTCUTTOGROUP:
		return "DPSP_MSG_TYPE_ADDSHORTCUTTOGROUP"
	case DPSP_MSG_TYPE_DELETEGROUPFROMGROUP:
		return "DPSP_MSG_TYPE_DELETEGROUPFROMGROUP"
	case DPSP_MSG_TYPE_SUPERENUMPLAYERSREPLY:
		return "DPSP_MSG_TYPE_SUPERENUMPLAYERSREPLY"
	case DPSP_MSG_TYPE_UNK4:
		return "DPSP_MSG_TYPE_UNK4"
	case DPSP_MSG_TYPE_KEYEXCHANGE:
		return "DPSP_MSG_TYPE_KEYEXCHANGE"
	case DPSP_MSG_TYPE_KEYEXCHANGEREPLY:
		return "DPSP_MSG_TYPE_KEYEXCHANGEREPLY"
	case DPSP_MSG_TYPE_CHAT:
		return "DPSP_MSG_TYPE_CHAT"
	case DPSP_MSG_TYPE_ADDFORWARD:
		return "DPSP_MSG_TYPE_ADDFORWARD"
	case DPSP_MSG_TYPE_ADDFORWARDACK:
		return "DPSP_MSG_TYPE_ADDFORWARDACK"
	case DPSP_MSG_TYPE_PACKET2_DATA:
		return "DPSP_MSG_TYPE_PACKET2_DATA"
	case DPSP_MSG_TYPE_PACKET2_ACK:
		return "DPSP_MSG_TYPE_PACKET2_ACK"
	case DPSP_MSG_TYPE_UNK5:
		return "DPSP_MSG_TYPE_UNK5"
	case DPSP_MSG_TYPE_UNK6:
		return "DPSP_MSG_TYPE_UNK6"
	case DPSP_MSG_TYPE_UNK7:
		return "DPSP_MSG_TYPE_UNK7"
	case DPSP_MSG_TYPE_IAMNAMESERVER:
		return "DPSP_MSG_TYPE_IAMNAMESERVER"
	case DPSP_MSG_TYPE_VOICE:
		return "DPSP_MSG_TYPE_VOICE"
	case DPSP_MSG_TYPE_MULTICASTDELIVERY:
		return "DPSP_MSG_TYPE_MULTICASTDELIVERY"
	case DPSP_MSG_TYPE_CREATEPLAYERVERIFY:
		return "DPSP_MSG_TYPE_CREATEPLAYERVERIFY"
	default:
		return strconv.Itoa(int(cmd))
	}
}

type DPlayPacket interface {
	Command() int
	CommandString() string
	Signature() string
	Size() int
	Version() int
	//Address() net.IP
	Port() int
	String() string
	Token() int
}

func NewDPlayPacket(data []byte) DPlayPacket {
	header := new(dpsp_MSG_HEADER)
	if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, header); err != nil {
		fmt.Println("binary.Read failed:", err)
		return nil
	}

	//Fix the Port, which for some reason is BigEndian
	header.SockAddr.Port = (header.SockAddr.Port >> 8) | (header.SockAddr.Port << 8)

	packet := &DPSP_PKT_HEADER{header.SizeAndToken, header.SockAddr, header.Signature, header.Command, header.Version}

	//This function will get big real fast...
	switch DPPacketType(packet.Command()) {
	case DPSP_MSG_TYPE_ENUMSESSIONS:
		packet := NewEnumSessionsPacket(data)
		return packet
	case DPSP_MSG_TYPE_ENUMSESSIONSREPLY:
		packet := NewEnumSessionsReplyPacket(data)
		return packet
	case DPSP_MSG_TYPE_REQUESTPLAYERID:
		return packet
	case DPSP_MSG_TYPE_REQUESTPLAYERREPLY:
		return packet
	case DPSP_MSG_TYPE_ADDFORWARDREQUEST:
		return packet
	case DPSP_MSG_TYPE_SUPERENUMPLAYERSREPLY:
		return packet
	case DPSP_MSG_TYPE_CREATEPLAYER:
		return packet
	case DPSP_MSG_TYPE_SESSIONDESCCHANGED:
		return packet
	case DPSP_MSG_TYPE_DELETEPLAYER:
		return packet
	default:
		return packet
	}

	return packet
}

//The SOCKADDR_IN structure is built as if it were on a little-endian machine and is treated as a byte array.
type SOCKADDR_IN struct {
	AddressFamily uint16 //Docs Say: MUST be 0x0002
	Port          uint16 //IP Port
	Address       uint32 //IP Address (RFC791)
	Padding       uint64 //Docs Say: MUST be 0 and MUST be ignored
}

type dpsp_MSG_HEADER struct {
	SizeAndToken uint32
	SockAddr     SOCKADDR_IN
	Signature    [4]byte
	Command      DPPacketType //NOTE: This matches my recognition of BG1 packets, NOT the official docs
	Version      uint16       // The docs show these last two reversed. My tests for, BG1 at least, do not.
}

type DPSP_PKT_HEADER struct {
	sizeAndToken uint32
	sockAddr     SOCKADDR_IN
	signature    [4]byte
	command      DPPacketType
	version      uint16
}

func (this *DPSP_PKT_HEADER) Command() int {
	return int(this.command)
}

func (this *DPSP_PKT_HEADER) CommandString() string {
	return commandToString(this.command)
}

func (this *DPSP_PKT_HEADER) Signature() string {
	return string(this.signature[:])
}

func (this *DPSP_PKT_HEADER) Size() int {
	return int((this.sizeAndToken << 12) >> 12)
}

func (this *DPSP_PKT_HEADER) Token() int {
	return int(this.sizeAndToken >> 20)
}

func (this *DPSP_PKT_HEADER) Version() int {
	return int(this.version)
}

func (this *DPSP_PKT_HEADER) Port() int {
	return int(this.sockAddr.Port)
}

func (this *DPSP_PKT_HEADER) String() string {
	ret := this.CommandString()
	ret += "\n\tPort:      " + strconv.Itoa(this.Port())
	ret += "\n\tVersion:   " + strconv.Itoa(this.Version())
	ret += "\n\tSize:      " + strconv.Itoa(this.Size())
	ret += "\n\tToken:      " + strconv.Itoa(this.Token()) + " - " + fmt.Sprintf("0x%X", this.Token())
	ret += "\n\tSignature: " + this.Signature()
	return ret
}

type dpsp_MSG_ACCESSGRANTED struct {
	dpsp_MSG_HEADER
	PublicKeySize   uint32 //DS: MUST be set to the size of the pubkey field. MUST be 24
	PublicKeyOffset uint32 //DS: MUST be offset from struct start to PublicKey (36?)
	PublicKey       []byte
}

type dpsp_MSG_ADDFORWARD struct {
	dpsp_MSG_HEADER
	IDTo           uint32 //DS: ID of the player to whom this message is being sent
	PlayerID       uint32 //DS: Identifier of the affected player
	GroupID        uint32 //DS: Identifier of the affected group.
	CreateOffset   uint32 //DS: Offset of the PlayerInfo field. MUST be set to 28
	PasswordOffset uint32 //DS: Not used. MUST be ignored
	PlayerInfo     DPLAYI_PACKEDPLAYER
}

type dpsp_MSG_FORWARDACK struct {
	dpsp_MSG_HEADER
	ID uint32 //DS: Identifier of the player for whom a dpsp_MSG_ADDFORWARD message was sent
}

type dpsp_MSG_ADDFORWARDREPLY struct {
	dpsp_MSG_HEADER
	Error uint32 //DS: Indicates the reason that the dpsp_MSG_ADDFORWARD (section 2.2.8) message failed. For a complete list of DirectPlay 4 HRESULT codes, see [MS-ERREF].
}

//======================================
//These are all used in BG1
//======================================
type dpSESSIONDESC2 struct {
	Size                uint32 //DS: MUST be the size of the struct
	Flags               uint32 //Game session flags: page 25 for bits
	InstGUID            [16]byte
	AppGUID             [16]byte
	MaxPlayers          uint32
	CurrentPlayerCount  uint32
	SessionName         uint32 //Pointer
	SessionPassword     uint32 //Pointer
	Reserved1           uint32 //DS:MUST be set to a unique value that is used to construct the player and group ID values. For more information about how this value is used to construct player and group identifiers, see section 3.2.5.4
	Reserved2           uint32 //DS: For future use
	ApplicationDefined1 uint32
	ApplicationDefined2 uint32
	ApplicationDefined3 uint32
	ApplicationDefined4 uint32
}

func (this *dpSESSIONDESC2) InstanceGUID() string {
	return fmt.Sprintf("%X%X%X%X-%X%X%X%X-%X%X%X%X-%X%X%X%X", this.InstGUID[0], this.InstGUID[1], this.InstGUID[2], this.InstGUID[3], this.InstGUID[4], this.InstGUID[5], this.InstGUID[6], this.InstGUID[7], this.InstGUID[8], this.InstGUID[9], this.InstGUID[10], this.InstGUID[11], this.InstGUID[12], this.InstGUID[13], this.InstGUID[14], this.InstGUID[15])
}

func (this *dpSESSIONDESC2) ApplicationGUID() string {
	return fmt.Sprintf("%X%X%X%X-%X%X%X%X-%X%X%X%X-%X%X%X%X", this.AppGUID[0], this.AppGUID[1], this.AppGUID[2], this.AppGUID[3], this.AppGUID[4], this.AppGUID[5], this.AppGUID[6], this.AppGUID[7], this.AppGUID[8], this.AppGUID[9], this.AppGUID[10], this.AppGUID[11], this.AppGUID[12], this.AppGUID[13], this.AppGUID[14], this.AppGUID[15])
}

func (this *dpSESSIONDESC2) FlagsCanCreateNewPlayers() bool {
	np := ((this.Flags << 31) >> 31) //NP (1 bit): Applications cannot create new players in this game session, as specified in dpsp_MSG_REQUESTPLAYERID
	if np == 1 {
		return false
	}
	return true
}

func (this *dpSESSIONDESC2) FlagsMigrateHost() bool {
	mh := ((this.Flags << 29) >> 31) //MH (1 bit): When the game host quits, the game host responsibilities migrate to another DirectPlay machine so that new players can continue to be created and nascent game instances can join the game session, as specified in section 3.1.6.2.
	if mh == 1 {
		return true
	}
	return false
}

func (this *dpSESSIONDESC2) FlagsIncludePlayerIDFields() bool {
	nm := ((this.Flags << 28) >> 31) //NM (1 bit): DirectPlay will not set the PlayerTo and PlayerFrom fields in player messages.
	if nm == 1 {
		return false
	}
	return true
}

func (this *dpSESSIONDESC2) FlagsJoinDisabled() bool {
	jd := ((this.Flags << 26) >> 31) //JD (1 bit): DirectPlay will not allow any new applications to join the game session. Applications already in the game session can still create new players
	if jd == 1 {
		return true
	}
	return false
}

func (this *dpSESSIONDESC2) FlagsUseDPPingTimer() bool {
	ka := ((this.Flags << 25) >> 31) //KA (1 bit): DirectPlay will detect when remote players exit abnormally (for example, because their computer or modem was unplugged) through the use of the Ping Timer, as described in sections 3.1.2.5 and 3.2.2.2.
	if ka == 1 {
		return true
	}
	return false
}

func (this *dpSESSIONDESC2) FlagsNoDataChangeUpdates() bool {
	nd := ((this.Flags << 24) >> 31) //ND (1 bit): DirectPlay will not send a message to all players when a player's remote data changes
	if nd == 1 {
		return true
	}
	return false
}

func (this *dpSESSIONDESC2) FlagsSecureSessionWithDPAuth() bool {
	ss := ((this.Flags << 23) >> 31) //SS (1 bit): Instructs the game session establishment logic to use user authentication as specified in sections 3.1.5.1 and 3.2.5.7
	if ss == 1 {
		return true
	}
	return false
}

func (this *dpSESSIONDESC2) FlagsPrivateSession() bool {
	p := ((this.Flags << 22) >> 31) //P (1 bit): Indicates that the game session is private and requires a password for EnumSessions as well as Open.
	if p == 1 {
		return true
	}
	return false
}

func (this *dpSESSIONDESC2) FlagsPasswordRequired() bool {
	pr := ((this.Flags << 21) >> 31) //PR (1 bit): Indicates that the game session requires a password to join.
	if pr == 1 {
		return true
	}
	return false
}

func (this *dpSESSIONDESC2) FlagsMessagesRouteViaHost() bool {
	ms := ((this.Flags << 20) >> 31) //MS (1 bit): DirectPlay will route all messages through the game host, as specified in section 3.1.5.1.
	if ms == 1 {
		return true
	}
	return false
}

func (this *dpSESSIONDESC2) FlagsCacheServerPlayerOnly() bool {
	cs := ((this.Flags << 19) >> 31) //CS (1 bit): DirectPlay will download information about the DPPLAYER_SERVERPLAYER only.
	if cs == 1 {
		return true
	}
	return false
}

func (this *dpSESSIONDESC2) FlagsReliableProtocolOnly() bool {
	rp := ((this.Flags << 18) >> 31) //RP (1 bit): Instructs the DirectPlay client to always use DirectPlay 4 Reliable Protocol [MCDPL4R]. When this bit is set, only other game sessions with the same bit set can join or be joined.
	if rp == 1 {
		return true
	}
	return false
}

func (this *dpSESSIONDESC2) FlagsNotOrderedReliablePackets() bool {
	no := ((this.Flags << 17) >> 31) //NO (1 bit): Instructs the DirectPlay client that, when using reliable delivery, preserving the order of received packets is not important. This allows messages to be indicated out of order if preceding messages have not yet arrived. If this flag is not set, DirectPlay waits for earlier messages to arrive before delivering later reliable messages.
	if no == 1 {
		return true
	}
	return false
}

func (this *dpSESSIONDESC2) FlagsOptimize4Latency() bool {
	ol := ((this.Flags << 16) >> 31) //OL (1 bit): DirectPlay will optimize communication for latency. Implementations SHOULD use the presence of the OL flag for guidance on how to send or process messages to optimize for latency rather than throughput; however, implementations can choose to ignore this flag. The presence or absence of the OL flag MUST NOT affect the sequence or binary contents of DirectPlay 4 protocol messages.<6>
	if ol == 1 {
		return true
	}
	return false
}

func (this *dpSESSIONDESC2) FlagsAcquireVoice() bool {
	av := ((this.Flags << 15) >> 31) //AV (1 bit): Allows lobby-launched games that are not voice-enabled to acquire voice capabilities.
	if av == 1 {
		return true
	}
	return false
}

func (this *dpSESSIONDESC2) FlagsNoSessionDescriptionChanges() bool {
	ns := ((this.Flags << 14) >> 31) //NS (1 bit): Suppresses transmission of game session description changes.
	if ns == 1 {
		return true
	}
	return false
}

func (this *dpSESSIONDESC2) FlagsToString() string {
	x := ((this.Flags << 30) >> 31) //X (1 bit): All bits with this label SHOULD be set to zero when sent and MUST be ignored on receipt.
	i := ((this.Flags << 27) >> 31) //I (1 bit): (Ignored). All bits with this label MUST be ignored on receipt.
	y := ((this.Flags >> 18) << 18) //Y (14 bits): All bits with this label SHOULD be set to zero when sent and MUST be ignored on receipt.
	ret := "Flags: "
	ret += "\n\t\tCanCreateNewPlayers: " + strconv.FormatBool(this.FlagsCanCreateNewPlayers())
	if x != 0 {
		ret += "\n\t\tX == " + strconv.Itoa(int(x)) + "?!?"
	}
	ret += "\n\t\tMigrateHost: " + strconv.FormatBool(this.FlagsMigrateHost())
	ret += "\n\t\tIncludePlayerIDFields: " + strconv.FormatBool(this.FlagsIncludePlayerIDFields())
	ret += "\n\t\tJoinDisabled: " + strconv.FormatBool(this.FlagsJoinDisabled())
	//if i != 0 {
	ret += "\n\t\tI == " + strconv.Itoa(int(i)) + "?!?"
	//}
	ret += "\n\t\tUseDPPingTimer: " + strconv.FormatBool(this.FlagsUseDPPingTimer())
	ret += "\n\t\tNoDataChangeUpdates: " + strconv.FormatBool(this.FlagsNoDataChangeUpdates())
	ret += "\n\t\tSecureSessionWithDPAuth: " + strconv.FormatBool(this.FlagsSecureSessionWithDPAuth())
	ret += "\n\t\tPrivateSession: " + strconv.FormatBool(this.FlagsPrivateSession())
	ret += "\n\t\tPasswordRequired: " + strconv.FormatBool(this.FlagsPasswordRequired())
	ret += "\n\t\tMessagesRouteViaHost: " + strconv.FormatBool(this.FlagsMessagesRouteViaHost())
	ret += "\n\t\tCacheServerPlayerOnly: " + strconv.FormatBool(this.FlagsCacheServerPlayerOnly())
	ret += "\n\t\tReliableProtocolOnly: " + strconv.FormatBool(this.FlagsReliableProtocolOnly())
	ret += "\n\t\tNotOrderedReliablePackets: " + strconv.FormatBool(this.FlagsNotOrderedReliablePackets())
	ret += "\n\t\tOptimize4Latency: " + strconv.FormatBool(this.FlagsOptimize4Latency())
	ret += "\n\t\tAcquireVoice: " + strconv.FormatBool(this.FlagsAcquireVoice())
	ret += "\n\t\tNoSessionDescriptionChanges: " + strconv.FormatBool(this.FlagsNoSessionDescriptionChanges())
	if y != 0 {
		ret += "\n\t\tY == " + strconv.Itoa(int(y)) + "?!?"
	}
	return ret
}

//Contains data related to players or groups
type DPLAYI_PACKEDPLAYER struct {
	Size                    uint32 //Docs Say:  MUST contain the total size of the DPLAYI_PACKEDPLAYER structure plus the values of the ShortNameLength, LongNameLength, ServiceProviderDataSize, and PlayerDataSize fields
	Flags                   uint32 //Docs Say: MUST contain 0 or more player flags: SystemPlayer, NameServer, PlayerGroup, PlayerIsLocal (all one bit. remaining bits "SHOULD" be set to zero and "MUST" be ignored)
	PlayerID                uint32
	ShortNameLength         uint32   //DS: MUST be zero if unavailable
	LongNameLength          uint32   //DS: MUST be zero if unavailable
	ServiceProviderDataSize uint32   //DS: MUST contain len(ServiceProviderData) or 0
	PlayerDataSize          uint32   //DS: MUST contain len(PlayerData) or 0
	NumberOfPlayers         uint32   //DS: MUST contain the number of entries in the PlayerIDs field. If the player represented in the DPLAYI_PACKEDPLAYER structure is not a group, this field MUST be set to zero.
	SystemPlayerID          uint32   //DS: MUST be the ID of the system player for the game session
	FixedSize               uint32   //DS: size of the fixed portion of this struct; MUST be 48
	PlayerVersion           uint32   //DS: MUST contain the version of the current player/group
	ParentID                uint32   //DS: MUST contain the identifier of the parent group. If this struct represents a player or a group that is not contained in another group, this field MUST be set to zero
	ShortName               string   //NULL Terminater Unicode
	LongName                string   //NULL Terminater Unicode
	ServiceProviderData     []byte   //May contain a Winsock thing? bottom pg18/top pg19
	PlayerData              []byte   //DS:  If PlayerDataSize is nonzero, this MUST be set to the byte array of gamespecific per-player data.
	PlayerIDs               []uint32 //Not sure on the type; DS: MUST contain an array of PlayerIDs where the array size is specified in NumberOfPlayers. If NumberOfPlayers is 0, this field MUST NOT be present.
}

//FIXME: This is really crude. We need a loader function for this
type DPLAYI_SUPERPACKEDPLAYER struct {
	Size                      uint32
	Flags                     uint32
	ID                        uint32 //PlayerID
	PlayerInfoMask            uint32
	VersionOrSystemPlayerID   uint32
	ShortName                 string
	LongName                  string
	PlayerDataLength          []byte
	PlayerData                []byte
	ServiceProviderDataLength []byte
	ServiceProviderData       []byte
	PlayerCount               []byte
	PlayerIDs                 []byte
	ParentID                  []byte
	ShortcutIDCount           []byte
	ShortcutIDs               []byte
}

type DPSECURITYDESC struct {
	Size                uint32 //DS: MUST be set to the size of the struct (Windows sets to zero says ref??)
	Flags               uint32 //DS: Game session flags. This is not used. MUST be set to zero when sent and MUST be ignored on receipt
	SSPIProvider        uint32 //DS: MUST be ignored on receipt
	CAPIProvider        uint32 //DS: MUST be ignored on receipt
	CAPIProviderType    uint32 //DS: Crypto service provider type. If the application does not specify a value, the default value of PROV_RSA_FULL is used. For more information, see Cryptographic Provider Types [MSDN-CAPI].
	EncryptionAlgorithm uint32 //DS: Encryption algorithm type. If the application does not specify a value, the default value of CALG_RC4 is used
}

type dpsp_MSG_ENUMSESSIONS struct {
	dpsp_MSG_HEADER
	ApplicationGUID [16]byte
	PasswordOffset  uint32
	Flags           uint32
	//Password        string
}

type DPSP_PKT_ENUMSESSIONS struct {
	DPSP_PKT_HEADER
	applicationGUID [16]byte
	passwordOffset  uint32
	flags           uint32
	password        string
}

func NewEnumSessionsPacket(data []byte) *DPSP_PKT_ENUMSESSIONS {
	rawpkt := new(dpsp_MSG_ENUMSESSIONS)
	if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, rawpkt); err != nil {
		fmt.Println("binary.Read failed:", err)
		return nil
	}

	//Fix the Port, which for some reason is BigEndian
	rawpkt.SockAddr.Port = (rawpkt.SockAddr.Port >> 8) | (rawpkt.SockAddr.Port << 8)

	header := DPSP_PKT_HEADER{rawpkt.SizeAndToken, rawpkt.SockAddr, rawpkt.Signature, rawpkt.Command, rawpkt.Version}
	ret := &DPSP_PKT_ENUMSESSIONS{header, rawpkt.ApplicationGUID, rawpkt.PasswordOffset, rawpkt.Flags, ""}

	return ret
}

func (this *DPSP_PKT_ENUMSESSIONS) FlagsListJoinable() bool {
	av := ((this.flags << 31) >> 31) //AV (1 bit): Enumerate game sessions that can be joined.
	if av == 1 {
		return true
	}
	return false
}

func (this *DPSP_PKT_ENUMSESSIONS) FlagsListUnjoinable() bool {
	al := ((this.flags << 30) >> 31) //AL (1 bit): Enumerate all game sessions, even if they cannot be joined.
	if al == 1 {
		return true
	}
	return false
}

func (this *DPSP_PKT_ENUMSESSIONS) FlagsListPasswordRequired() bool {
	pr := ((this.flags << 25) >> 31) //PR (1 bit): Enumerate all game sessions, even if they require a password.
	if pr == 1 {
		return true
	}
	return false
}

func (this *DPSP_PKT_ENUMSESSIONS) FlagsToString() string {
	x := ((this.flags << 26) >> 28) //X (4 bits): Not used. SHOULD be set to zero when sent and MUST be ignored on receipt.
	y := ((this.flags >> 7) << 7)   //Y (25 bits): Not used. SHOULD be set to zero when sent and MUST be ignored on receipt.

	ret := "Flags: "
	if this.FlagsListJoinable() {
		ret += "List Joinable | "
	}

	if this.FlagsListUnjoinable() {
		ret += "List Unjoinable | "
	}

	if x != 0 {
		ret += "X == " + strconv.Itoa(int(x)) + " | "
	}

	if this.FlagsListPasswordRequired() {
		ret += "List with Password Required | "
	}

	if y != 0 {
		ret += "Y == " + strconv.Itoa(int(y)) + " | "
	}
	return ret
}

func (this *DPSP_PKT_ENUMSESSIONS) ApplicationGUID() string {
	return fmt.Sprintf("%X%X%X%X-%X%X%X%X-%X%X%X%X-%X%X%X%X", this.applicationGUID[0], this.applicationGUID[1], this.applicationGUID[2], this.applicationGUID[3], this.applicationGUID[4], this.applicationGUID[5], this.applicationGUID[6], this.applicationGUID[7], this.applicationGUID[8], this.applicationGUID[9], this.applicationGUID[10], this.applicationGUID[11], this.applicationGUID[12], this.applicationGUID[13], this.applicationGUID[14], this.applicationGUID[15])
}

func (this *DPSP_PKT_ENUMSESSIONS) String() string {
	ret := this.CommandString()
	ret += "\n\tPort:      " + strconv.Itoa(this.Port())
	ret += "\n\tVersion:   " + strconv.Itoa(this.Version())
	ret += "\n\tSize:      " + strconv.Itoa(this.Size())
	ret += "\n\tToken:      " + strconv.Itoa(this.Token()) + " - " + fmt.Sprintf("0x%X", this.Token())
	ret += "\n\tSignature: " + this.Signature()
	ret += "\n\t---"
	ret += "\n\tApplication GUID: " + this.ApplicationGUID()
	if this.passwordOffset != 0 {
		ret += "\n\tPassword Offset: " + strconv.Itoa(int(this.passwordOffset))
	}
	ret += this.FlagsToString()
	if this.passwordOffset > 0 {
		ret += "\n\tPassword: TODO? - Never seen one with this"
	}
	return ret
}

type dpsp_MSG_ENUMSESSIONSREPLY struct {
	dpsp_MSG_HEADER
	SessionDescription dpSESSIONDESC2
	NameOffset         uint32 //Not sure why, but this number is 20 less than the offset from 0
	//SessionName        string
	//SessionName []byte
}

type DPSP_PKT_ENUMSESSIONSREPLY struct {
	DPSP_PKT_HEADER
	sessionDesc dpSESSIONDESC2
	nameOffset  uint32
	sessionName string
}

// UTF16BytesToString converts UTF-16 encoded bytes, in big or little endian byte order,
// to a UTF-8 encoded string.
func UTF16BytesToString(b []byte, o binary.ByteOrder) string {
	utf := make([]uint16, (len(b)+(2-1))/2)
	for i := 0; i+(2-1) < len(b); i += 2 {
		utf[i/2] = o.Uint16(b[i:])
	}
	if len(b)/2 < len(utf) {
		utf[len(utf)-1] = utf8.RuneError
	}
	return string(utf16.Decode(utf))
}

func NewEnumSessionsReplyPacket(data []byte) *DPSP_PKT_ENUMSESSIONSREPLY {
	rawpkt := new(dpsp_MSG_ENUMSESSIONSREPLY)
	if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, rawpkt); err != nil {
		fmt.Println("binary.Read failed:", err)
		return nil
	}

	//Fix the Port, which for some reason is BigEndian
	rawpkt.SockAddr.Port = (rawpkt.SockAddr.Port >> 8) | (rawpkt.SockAddr.Port << 8)

	header := DPSP_PKT_HEADER{rawpkt.SizeAndToken, rawpkt.SockAddr, rawpkt.Signature, rawpkt.Command, rawpkt.Version}
	n := int(rawpkt.NameOffset) + 20 //Not sure why, but this number is 20 less than the offset from 0
	ret := &DPSP_PKT_ENUMSESSIONSREPLY{header, rawpkt.SessionDescription, rawpkt.NameOffset, UTF16BytesToString(data[n:], binary.LittleEndian)}
	return ret
}

func (this *DPSP_PKT_ENUMSESSIONSREPLY) String() string {
	ret := this.CommandString()
	ret += "\n\tPort:      " + strconv.Itoa(this.Port())
	ret += "\n\tVersion:   " + strconv.Itoa(this.Version())
	ret += "\n\tSize:      " + strconv.Itoa(this.Size())
	ret += "\n\tToken:      " + strconv.Itoa(this.Token()) + " - " + fmt.Sprintf("0x%X", this.Token())
	ret += "\n\tSignature: " + this.Signature()
	ret += "\n\t---"
	ret += "\n\tName Offset: " + strconv.Itoa(int(this.nameOffset))
	ret += "\n\tSession Name: '" + this.sessionName + "'"
	ret += "\n\tSessionDesc:"
	ret += "\n\tSize: " + strconv.Itoa(int(this.sessionDesc.Size))
	ret += "\n\tApp GUID: " + this.sessionDesc.ApplicationGUID()
	ret += "\n\tInst GUID: " + this.sessionDesc.InstanceGUID()
	ret += "\n\t" + this.sessionDesc.FlagsToString()
	ret += "\n\tCurrent Player Count: " + strconv.Itoa(int(this.sessionDesc.CurrentPlayerCount))
	ret += "\n\tMax Players: " + strconv.Itoa(int(this.sessionDesc.MaxPlayers))
	ret += "\n\tReserved1: " + strconv.Itoa(int(this.sessionDesc.Reserved1)) + " - " + fmt.Sprintf("0x%X", int(this.sessionDesc.Reserved1))
	ret += "\n\tReserved2: " + strconv.Itoa(int(this.sessionDesc.Reserved2))
	ret += "\n\tSessionName Pointer: " + strconv.Itoa(int(this.sessionDesc.SessionName))
	ret += "\n\tSessionPassword Pointer: " + strconv.Itoa(int(this.sessionDesc.SessionPassword))
	ret += "\n\tApplicationDefined1: " + strconv.Itoa(int(this.sessionDesc.ApplicationDefined1))
	ret += "\n\tApplicationDefined2: " + strconv.Itoa(int(this.sessionDesc.ApplicationDefined2))
	ret += "\n\tApplicationDefined3: " + strconv.Itoa(int(this.sessionDesc.ApplicationDefined3))
	ret += "\n\tApplicationDefined4: " + strconv.Itoa(int(this.sessionDesc.ApplicationDefined4))
	return ret
}

type dpsp_MSG_REQUESTPLAYERID struct {
	dpsp_MSG_HEADER
	Flags uint32
}

type dpsp_MSG_REQUESTPLAYERREPLY struct {
	dpsp_MSG_HEADER
	ID                 uint32
	SecDesc            DPSECURITYDESC
	SSPIProviderOffset uint32
	CAPIProviderOffset uint32
	Result             uint32
	SSPIProvider       string
	CAPIProvider       string
}

type dpsp_MSG_ADDFORWARDREQUEST struct {
	dpsp_MSG_HEADER
	IDTo           uint32 //DS: ID of the player to whom this message is being sent
	PlayerID       uint32 //DS: MUST be the identity of the player being added.
	GroupID        uint32 //DS: SHOULD be set to zero when sent and MUST be ignored
	CreateOffset   uint32 //DS: Offset, in bytes, of the PlayerInfo field from the beginning of the Signature field in the dpsp_MSG_HEADER. SHOULD be 28
	PasswordOffset uint32
	PlayerInfo     DPLAYI_PACKEDPLAYER
	Password       string
	TickCount      uint32
}

type dpsp_MSG_SUPERENUMPLAYERSREPLY struct {
	dpsp_MSG_HEADER
	PlayerCount       uint32
	GroupCount        uint32
	PackedOffset      uint32
	ShortcutCount     uint32
	DescriptionOffset uint32
	NameOffset        uint32
	PasswordOffset    uint32
	DPSessionDesc     dpSESSIONDESC2
	SessionName       string
	Password          string
	SuperPackedPlayer DPLAYI_SUPERPACKEDPLAYER
}

//At this point, the player is asked for their username.

type dpsp_MSG_CREATEPLAYER struct {
	dpsp_MSG_HEADER
	IDTo           uint32
	PlayerID       uint32
	GroupID        uint32
	CreateOffset   uint32
	PasswordOffset uint32
	PlayerInfo     DPLAYI_PACKEDPLAYER
	Reserved1      uint16
	Reserved2      uint32
}

type dpsp_MSG_SESSIONDESCCHANGED struct {
	dpsp_MSG_HEADER
	IDTo              uint32
	SessionNameOffset uint32
	PasswordOffset    uint32
	SessionDesc       dpSESSIONDESC2
	SessionName       string
	Password          string
}

type dpsp_MSG_DELETEPLAYER struct {
	dpsp_MSG_HEADER
	IDTo           uint32
	PlayerID       uint32
	GroupID        uint32
	CreateOffset   uint32
	PasswordOffset uint32
}
