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

	in := bufio.NewReader(os.Stdin)

	go func() {
		defer close(me.Send)
		for {
			line, err := in.ReadString('\n')
			if err == io.EOF {
				return
			}

			if err != nil {
				log.Fatalf("readline(): %s", err)
			}

			me.Send <- packet.Packet{
				Chat: &packet.Chat{
					Message: line,
				},
			}
		}
	}()

	go func() {
		for msg := range me.Recv {
			clientpkg.Handle(msg)
		}
		clientpkg.Disconnected()
	}()

	clientpkg.Main()
}
