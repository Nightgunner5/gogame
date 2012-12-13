package packet

type Packet struct {
	*Handshake
	*Location
	*Despawn
	*MapChange
	*MapOverride
}
