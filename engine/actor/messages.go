package actor

import (
	"github.com/Nightgunner5/gogame/engine/message"
)

var (
	MsgAddSubscriber    = message.NewKind("AddSubscriber")
	MsgRemoveSubscriber = message.NewKind("RemoveSubscriber")
)

type AddSubscriber struct {
	Subscriber
}

func (AddSubscriber) Kind() message.Kind {
	return MsgAddSubscriber
}

type RemoveSubscriber struct {
	Subscriber
}

func (RemoveSubscriber) Kind() message.Kind {
	return MsgRemoveSubscriber
}
