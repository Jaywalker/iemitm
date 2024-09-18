package interprocess

type PacketData struct {
	Source, Dest, Port string
	Size               int
	Data               []byte
}

type RespPacketData struct {
	Dest    string
	Size    int
	Forward bool
	Data    []byte
}
