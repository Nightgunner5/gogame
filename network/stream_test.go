package network

import (
	"net"
	"testing"
)

const (
	TestPacketID PacketID = iota
	TestString
	TestNumber
)

func TestSendRecieve(t *testing.T) {
	a, b := net.Pipe()
	reader, writer := DecodeStream(a), EncodeStream(b)

	send := NewPacket(TestPacketID).
		Set(TestString, "abc").
		Set(TestNumber, 123.0)

	writer <- send
	recv := <-reader
	close(writer)

	if send.ID() != recv.ID() {
		t.Errorf("Sent ID (%d) != Recieved ID (%d)", send.ID(), recv.ID())
	}

	if send.FieldCount() == recv.FieldCount() {
		send.Each(func(k PacketID, v interface{}) {
			if recv.Get(k) != v {
				t.Errorf("Sent Payload [%q] %T(%#v) != Recieved Payload [%q] %T(%#v)", k, v, v, k, recv.Get(k), recv.Get(k))
			}
		})
	} else {
		t.Errorf("Sent Payload (length %d) != Recieved Payload (length %d)", send.FieldCount(), recv.FieldCount())
	}
}
