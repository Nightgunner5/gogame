package actor

import (
	"github.com/Nightgunner5/gogame/engine/message"
)

type Subscriber struct {
	kind  message.Kind
	queue message.Sender
}

// Constructs a Subscriber and returns it and the reciever end of its channel.
// maxWaiting is the size of the buffer for the channel. If a send to the
// channel would block, the Message is simply dropped.
func Subscribe(kind message.Kind, maxWaiting int) (Subscriber, message.Receiver) {
	queue := make(chan message.Message, maxWaiting)

	return Subscriber{
		kind:  kind,
		queue: queue,
	}, queue
}

func (s Subscriber) offer(msg message.Message) {
	if msg.Kind() == s.kind {
		select {
		case s.queue <- msg:
		default:
		}
	}
}
