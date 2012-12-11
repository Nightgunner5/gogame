package packet

type Packet struct {
	*Chat
	*Handshake
	*Location
}
