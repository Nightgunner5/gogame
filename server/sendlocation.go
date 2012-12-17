package server

import (
	"github.com/Nightgunner5/gogame/engine/message"
	"github.com/Nightgunner5/gogame/shared/packet"
)

var (
	MsgSendLocation = message.NewKind("SendLocation")
	MsgMoveIntoView = message.NewKind("MoveIntoView")
)

type SendLocation struct {
	*Player
}

type MoveIntoView struct {
	*packet.Location
}

func (SendLocation) Kind() message.Kind {
	return MsgSendLocation
}

func (MoveIntoView) Kind() message.Kind {
	return MsgMoveIntoView
}
