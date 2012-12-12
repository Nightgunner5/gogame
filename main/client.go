// +build !noclient

package main

import (
	"bufio"
	clientpkg "github.com/Nightgunner5/gogame/client"
	"github.com/Nightgunner5/gogame/shared/packet"
	"github.com/kylelemons/fatchan"
	"io"
	"log"
	"os"
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

	in := bufio.NewReader(os.Stdin)

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
