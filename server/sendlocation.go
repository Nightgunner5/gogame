package server

import (
	"github.com/Nightgunner5/gogame/engine/message"
	"github.com/Nightgunner5/gogame/shared/packet"
)

var (
	MsgSendLocation = message.NewKind("SendLocation")
)

type SendLocation chan<- packet.Packet

func (SendLocation) Kind() message.Kind {
	return MsgSendLocation
}
