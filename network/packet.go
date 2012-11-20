package network

import (
	"encoding/json"
	"io"
	"log"
)

type Packet struct {
	ID      uint64
	Payload map[string]interface{}
}

func DecodeStream(r io.ReadCloser) <-chan Packet {
	ch := make(chan Packet)
	decoder := json.NewDecoder(r)

	go func() {
		for {
			var p Packet
			if err := decoder.Decode(&p); err != nil {
				if err != io.EOF {
					log.Print("Error decoding packet: ", err)
				}
				r.Close()
				close(ch)
				return
			}
			ch <- p
		}
	}()

	return ch
}

func EncodeStream(w io.WriteCloser) chan<- Packet {
	ch := make(chan Packet)
	encoder := json.NewEncoder(w)

	go func() {
		for p := range ch {
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
