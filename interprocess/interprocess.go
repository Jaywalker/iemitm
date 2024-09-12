package interprocess

type PacketData struct {
	Source, Dest, Port string
	Size               int
	Data               []byte
}
