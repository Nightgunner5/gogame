package network

const base32Alphabet = "0123456789abcdefghijklmnopqrstuv"

func toBase32(n PacketID) (s string) {
	for n > 0 {
		s = base32Alphabet[n&0x1F:][:1] + s
		n >>= 5
	}
	if s == "" {
		s = "0"
	}
	return
}

func fromBase32(s string) (n PacketID) {
	for i := range s {
		b := s[i]
		n <<= 5
		if b >= '0' && b <= '9' {
			n |= PacketID(b - '0')
		} else if b >= 'a' && b <= 'v' {
			n |= PacketID(b - 'a' + 10)
		}
	}
	return
}

type packetTransmit struct {
	ID      string                 `json:"i"`
	Payload map[string]interface{} `json:"p"`
}

func (p *packetTransmit) fromPacket(packet Packet) {
	p.ID = toBase32(packet.ID())
	p.Payload = make(map[string]interface{})
	packet.Each(func(k PacketID, v interface{}) {
		p.Payload[toBase32(k)] = v
	})
}

func (p *packetTransmit) toPacket() Packet {
	packet := NewPacket(fromBase32(p.ID))
	for k, v := range p.Payload {
		packet.Set(fromBase32(k), v)
	}
	return packet
}
