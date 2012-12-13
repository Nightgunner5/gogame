// +build noclient,noserver,monkey

package main

import (
	"github.com/Nightgunner5/fatchan"
	"github.com/Nightgunner5/gogame/shared/layout"
	"github.com/Nightgunner5/gogame/shared/packet"
	"io"
	"log"
	"math/rand"
	"time"
)

const (
	DefaultAddr   = "127.0.0.1:7031"
	DefaultServer = false
)

func client(username string, server io.ReadWriteCloser) {
	log.Print("Connected to server")
	defer log.Print("Server disconnected")

	xport := fatchan.New(server, nil)
	login := make(chan Handshake)
	xport.FromChan(login)
	defer close(login)

	me := Handshake{
		User: "monkey",
		Send: make(chan packet.Packet),
		Recv: make(chan packet.Packet),
	}
	login <- me

	defer close(me.Send)

	for {
		select {
		case _, ok := <-me.Recv:
			if !ok {
				return
			}
		default:
			me.Send <- packet.Packet{
				Location: &packet.Location{
					Coord: layout.Coord{rand.Intn(3) - 1, rand.Intn(3) - 1},
				},
			}
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
	}
}
