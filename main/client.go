// +build !noclient

package main

import (
	"github.com/Nightgunner5/fatchan"
	clientpkg "github.com/Nightgunner5/gogame/client"
	"github.com/Nightgunner5/gogame/shared/packet"
	"io"
	"log"
)

const (
	DefaultAddr   = "nightgunner5.is-a-geek.net:7031"
	DefaultServer = false
)

func client(username string, server io.ReadWriteCloser) {
	log.Printf("Connected to server")
	defer log.Printf("Server disconnected")

	xport := fatchan.New(server, nil)
	login := make(chan Handshake)
	xport.FromChan(login)
	defer close(login)

	me := Handshake{
		User: username,
		Send: make(chan packet.Packet),
		Recv: make(chan packet.Packet),
	}
	login <- me

	defer close(me.Send)

	clientpkg.Network = me.Send

	go dispatch(me.Recv)

	clientpkg.Main()
}

func dispatch(recv <-chan packet.Packet) {
	for msg := range recv {
		clientpkg.Handle(msg)
	}
	clientpkg.Disconnected()
}
