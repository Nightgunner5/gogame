package network

import (
	"encoding/json"
	"io"
	"log"
)

func DecodeStream(r io.ReadCloser) <-chan Packet {
	ch := make(chan Packet)
	decoder := json.NewDecoder(r)

	go func() {
		var p packetTransmit
		for {
			if err := decoder.Decode(&p); err != nil {
				if err != io.EOF {
					log.Print("Error decoding packet: ", err)
				}
				r.Close()
				close(ch)
				return
			}
			ch <- p.toPacket()
		}
	}()

	return ch
}

func EncodeStream(w io.WriteCloser) chan<- Packet {
	ch := make(chan Packet)
	encoder := json.NewEncoder(w)

	go func() {
		var p packetTransmit
		for packet := range ch {
			p.fromPacket(packet)
			if err := encoder.Encode(p); err != nil {
				log.Print("Error encoding packet: ", err)
				w.Close()
				return
			}
		}
		w.Close()
	}()

	return ch
}
