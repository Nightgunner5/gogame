package network

import (
	"net"
	"testing"
)

const TestPacketID uint64 = 1337

func TestSendRecieve(t *testing.T) {
	a, b := net.Pipe()
	reader, writer := DecodeStream(a), EncodeStream(b)

	send := Packet{
		ID: TestPacketID,
		Payload: map[string]interface{}{
			"string": "abc",
			"number": 123.0, // Recieved numbers are always float64
		},
	}

	writer <- send
	recv := <-reader
	close(writer)

	if send.ID != recv.ID {
		t.Errorf("Sent ID (%d) != Recieved ID (%d)", send.ID, recv.ID)
	}

	if len(send.Payload) == len(recv.Payload) {
		for k, v := range send.Payload {
			if recv.Payload[k] != v {
				t.Errorf("Sent Payload [%q] %T(%#v) != Recieved Payload [%q] %T(%#v)", k, v, v, k, recv.Payload[k], recv.Payload[k])
			}
		}
	} else {
		t.Errorf("Sent Payload (length %d) != Recieved Payload (length %d)", len(send.Payload), len(recv.Payload))
	}
}
