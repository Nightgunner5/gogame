package actor

import (
	"github.com/Nightgunner5/gogame/engine/message"
)

type Subscriber struct {
	kind  message.Kind
	queue message.Sender
}

func Subscribe(kind message.Kind, maxWaiting int) (Subscriber, message.Reciever) {
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
