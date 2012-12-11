package packet

import (
	"github.com/Nightgunner5/gogame/engine/message"
)

var (
	MsgHandshake = message.NewKind("Handshake")
)

type Handshake struct {
	ID uint64
}

func (Handshake) Kind() message.Kind {
	return MsgHandshake
}
