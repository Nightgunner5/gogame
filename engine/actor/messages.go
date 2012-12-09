package actor

import (
	"github.com/Nightgunner5/gogame/engine/message"
)

var (
	MsgAddSubscriber    = message.NewKind("AddSubscriber")
	MsgRemoveSubscriber = message.NewKind("RemoveSubscriber")
)

// Add a Subscriber to an Actor. Subscribers are described in the documentation
// of Subscribe.
type AddSubscriber struct {
	Subscriber
}

// Returns MsgAddSubscriber.
func (AddSubscriber) Kind() message.Kind {
	return MsgAddSubscriber
}

// Remove a Subscriber from an Actor. Subscribers are described in the
// documentation of Subscribe.
type RemoveSubscriber struct {
	Subscriber
}

// Returns MsgRemoveSubscriber.
func (RemoveSubscriber) Kind() message.Kind {
	return MsgRemoveSubscriber
}
