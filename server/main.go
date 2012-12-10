package main

import (
	"github.com/Nightgunner5/gogame/shared/packet"
	"github.com/Nightgunner5/netchan"
	"net"
)

const (
	MaxQueue = 8
)

var (
	sendToAll    chan<- packet.Packet
	connected    chan<- chan<- packet.Packet
	disconnected chan<- chan<- packet.Packet
)

func init() {
	sendToAll_ := make(chan packet.Packet)
	sendToAll = sendToAll_

	connected_ := make(chan chan<- packet.Packet)
	connected = connected_

	disconnected_ := make(chan chan<- packet.Packet)
	disconnected = disconnected_

	go func() {
		connections := make(map[chan<- packet.Packet]bool)
		for {
			select {
			case p := <-sendToAll_:
				for c := range connections {
					select {
					case c <- p:
					default:
					}
				}
			case c := <-disconnected_:
				delete(connections, c)
			case c := <-connected_:
				connections[c] = true
			}
		}
	}()
}

func main() {
	ln, err := net.Listen("tcp", ":7031")
	if err != nil {
		panic(err)
	}

	netchan.Listen(ln, func(addr net.Addr, c *netchan.Chan) {
		recv := c.ChanRecv().(<-chan packet.Packet)
		send := c.ChanSend().(chan<- packet.Packet)

		NewPlayer(addr, recv, send)

		world.onConnect <- send
		connected <- send
	}, packet.Type, MaxQueue)
}
