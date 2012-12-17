package server

import (
	"github.com/Nightgunner5/gogame/engine/message"
)

var (
	MsgSendLocation = message.NewKind("SendLocation")
)

type SendLocation struct {
	*Player
}

func (SendLocation) Kind() message.Kind {
	return MsgSendLocation
}
