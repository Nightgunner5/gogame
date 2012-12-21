package packet

type Packet struct {
	*Handshake
	*Location
	*Interact
	*Despawn
	*MapChange
	*MapOverride
}
