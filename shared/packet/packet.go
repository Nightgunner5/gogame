package packet

import "reflect"

type Packet struct {
	*HandshakeServer
	*PlayerLocation
	*MoveRequest
}

var Type = reflect.TypeOf(Packet{})
