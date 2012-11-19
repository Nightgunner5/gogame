package main

import (
	"net/http"
	"github.com/Nightgunner5/gogame/entity"
	"github.com/Nightgunner5/gogame/spell"
	"sync"
	"time"
	"math/rand"
)

var (
	magicianForIP = make(map[string]entity.EntityID)
	magicianForIPLock sync.Mutex
)

func getMagician(ip string) Magician {
	magicianForIPLock.Lock()
	defer magicianForIPLock.Unlock()

	if id, ok := magicianForIP[ip]; ok {
		if ent := entity.Get(id); ent != nil {
			return ent.(Magician)
		}
	}

	x := rand.Float64() * 20 - 10
	y := rand.Float64() * 20 - 10

	m := NewMagician(x, y, 0)
	magicianForIP[ip] = m.ID()
	return m
}

func init() {
	go func() {
		for {
			time.Sleep(time.Minute)

			magicianForIPLock.Lock()
			for ip, id := range magicianForIP {
				if entity.Get(id) == nil {
					delete(magicianForIP, ip)
				}
			}
			magicianForIPLock.Unlock()
		}
	}()

	const (
		summonCost     = 50
		summonCastTime = 0.5
		shieldCost     = 20
		shieldCastTime = 0.2
	)

	http.HandleFunc("/cast/shield", func(w http.ResponseWriter, r *http.Request) {
		m := getMagician(r.Header.Get("X-Forwarded-For"))
		if m == nil {
			return
		}

		if len(m.Effects()) == 0 && m.UseResource(shieldCost) {
			m.Cast(&spell.BasicSpell{
				CastTime: shieldCastTime,
				Caster_:  m.ID(),
				Target_:  m.ID(),
				Action:   summonShield,

				Tag: "summonshield",
			})
		}
	})

	http.HandleFunc("/cast/imp", func(w http.ResponseWriter, r *http.Request) {
		m := getMagician(r.Header.Get("X-Forwarded-For"))
		if m == nil {
			return
		}

		if m.UseResource(summonCost) {
			m.Cast(&spell.BasicSpell{
				CastTime: summonCastTime,
				Caster_:  m.ID(),
				Target_:  m.ID(),
				Action:   summonImp,

				Tag: "summonimp",
			})
		}
	})
}
