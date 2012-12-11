package packet

import (
	"github.com/Nightgunner5/gogame/engine/message"
)

var (
	MsgDespawn = message.NewKind("Despawn")
)

type Despawn struct {
	ID uint64
}

func (Despawn) Kind() message.Kind {
	return MsgDespawn
}
