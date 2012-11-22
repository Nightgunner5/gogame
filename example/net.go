package main

import (
	"github.com/Nightgunner5/gogame/entity"
	"github.com/Nightgunner5/gogame/network"
	"github.com/Nightgunner5/gogame/spell"
	"math/rand"
	"net"
	"sync"
	"time"
)

const (
	Handshake = network.FirstUnusedPacketID + iota
	CastSpell
	KeepAlive
)

type netMagician struct {
	id       entity.EntityID
	lastSeen time.Time
}

var (
	magicians    = make(map[string]*netMagician)
	magicianLock sync.Mutex
)

func init() {
	go func() {
		for {
			time.Sleep(time.Minute)

			magicianLock.Lock()
			for addr, m := range magicians {
				if time.Since(m.lastSeen) > time.Minute {
					entity.Despawn(entity.Get(m.id))
					delete(magicians, addr)
				}
			}
			magicianLock.Unlock()
		}
	}()

	network.RegisterHandler(Handshake, func(packet network.Packet, send chan<- network.Packet, addr net.Addr) {
		magicianLock.Lock()
		defer magicianLock.Unlock()

		if _, exists := magicians[addr.String()]; exists {
			return
		}

		x := rand.Float64()*20 - 10
		y := rand.Float64()*20 - 10

		magician := NewMagician(x, y, 0)

		magicians[addr.String()] = &netMagician{
			id:       magician.ID(),
			lastSeen: time.Now(),
		}

		send <- network.NewPacket(Handshake).Set(network.EntityID, magician.ID())
	})

	network.RegisterHandler(CastSpell, func(packet network.Packet, send chan<- network.Packet, addr net.Addr) {
		magicianLock.Lock()
		if _, ok := magicians[addr.String()]; !ok {
			magicianLock.Unlock()
			return
		}
		magicians[addr.String()].lastSeen = time.Now()
		var magician Magician
		if m, ok := entity.Get(magicians[addr.String()].id).(Magician); ok {
			magician = m
		} else {
			magicianLock.Unlock()
			return
		}
		magicianLock.Unlock()

		if name, ok := packet.Get(CastSpell).(string); ok {
			switch name {
			case "imp":
				const (
					cost = 50
					time = 2
				)
				if magician.UseResource(cost) {
					magician.Cast(&spell.BasicSpell{
						CastTime: time,
						Caster_:  magician.ID(),
						Target_:  magician.ID(),
						Action:   summonImp,
					})
				}
			case "shield":
				const (
					cost = 20
					time = 0.3
				)
				if magician.UseResource(cost) {
					magician.Cast(&spell.BasicSpell{
						CastTime: time,
						Caster_:  magician.ID(),
						Target_:  magician.ID(),
						Action:   summonShield,
					})
				}
			}
		}
	})

	network.RegisterHandler(KeepAlive, func(packet network.Packet, send chan<- network.Packet, addr net.Addr) {
		magicianLock.Lock()
		if _, ok := magicians[addr.String()]; !ok {
			magicianLock.Unlock()
			return
		}
		magicians[addr.String()].lastSeen = time.Now()
		magicianLock.Unlock()
		send <- network.NewPacket(KeepAlive)
	})
}
