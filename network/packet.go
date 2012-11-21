package network

type Packet struct {
	id      PacketID
	payload map[PacketID]interface{}
}

func NewPacket(id PacketID) Packet {
	return Packet{
		id:      id,
		payload: make(map[PacketID]interface{}),
	}
}

func (p Packet) ID() PacketID {
	return p.id
}

func (p Packet) Set(key PacketID, value interface{}) Packet {
	p.payload[key] = value
	return p
}

func (p Packet) Get(key PacketID) interface{} {
	return p.payload[key]
}

func (p Packet) Each(f func(PacketID, interface{})) Packet {
	for k, v := range p.payload {
		f(k, v)
	}
	return p
}

func (p Packet) FieldCount() int {
	return len(p.payload)
}
