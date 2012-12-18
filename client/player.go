package client

import (
	"github.com/Nightgunner5/gogame/engine/actor"
	"github.com/Nightgunner5/gogame/engine/message"
	"github.com/Nightgunner5/gogame/shared/layout"
	"github.com/Nightgunner5/gogame/shared/packet"
	"github.com/Nightgunner5/gogame/shared/res"
	"image"
	"sync/atomic"
	"time"
)

type Player struct {
	actor.Actor
	x, y          int
	isLocalPlayer bool
	paint         *PaintContext
}

func (p *Player) Initialize() (message.Receiver, func(message.Message)) {
	msgIn, broadcast := p.Actor.Initialize()

	messages := make(chan message.Message)

	go p.dispatch(msgIn, messages)

	return messages, broadcast
}

func (p *Player) dispatch(msgIn message.Receiver, messages message.Sender) {
	for msg := range msgIn {
		switch m := msg.(type) {
		case SetLocation:
			if p.x != m.X || p.y != m.Y {
				Invalidate(p.screenRect())
				p.x, p.y = m.X, m.Y

				paintLock.RLock()
				if p.paint.Sprite == HumanSuit || p.paint.Sprite == HumanSuitHelm {
					if layout.Get(p.x, p.y).Space() {
						p.paint.Sprite = HumanSuitHelm
					} else {
						p.paint.Sprite = HumanSuit
					}
				}
				if p.paint.Changed.IsZero() {
					p.paint.From.X, p.paint.From.Y = p.x, p.y
				} else {
					p.paint.From = p.paint.To
				}
				p.paint.To.X, p.paint.To.Y = p.x, p.y
				p.paint.Changed = time.Now()
				paintLock.RUnlock()

				if p.isLocalPlayer {
					atomic.StoreInt64(&topLeftX, ViewportWidth/2-int64(p.x))
					atomic.StoreInt64(&topLeftY, ViewportHeight/2-int64(p.y))
				} else {
					Invalidate(p.screenRect())
				}
			}

		case packet.Despawn:
			paintLock.Lock()
			delete(paintContexts, &p.Actor)
			paintLock.Unlock()
			Invalidate(p.screenRect())

		default:
			messages <- m
		}
	}
	close(messages)
}

func (p *Player) screenRect() image.Rectangle {
	if p.isLocalPlayer {
		return image.Rect(0, 0, ViewportWidth<<res.TileSize, ViewportHeight<<res.TileSize)
	}
	x, y := GetTopLeft()
	return image.Rect((p.x+x)<<res.TileSize, (p.y+y)<<res.TileSize,
		(p.x+1+x)<<res.TileSize, (p.y+1+y)<<res.TileSize)
}

var thePlayer = NewPlayer(true, false)

func NewPlayer(isLocalPlayer, monkey bool) *Player {
	player := new(Player)
	player.isLocalPlayer = isLocalPlayer
	player.paint = new(PaintContext)
	if monkey {
		player.paint.Sprite = Monkey
	} else {
		player.paint.Sprite = HumanSuit
	}
	if isLocalPlayer {
		player.paint.Changed = time.Now()
	}
	paintLock.Lock()
	paintContexts[&player.Actor] = player.paint
	paintLock.Unlock()
	actor.Init("client:player", &player.Actor, player)
	return player
}

var (
	MsgSetLocation = message.NewKind("SetLocation")
)

type SetLocation struct {
	layout.Coord
}

func (SetLocation) Kind() message.Kind {
	return MsgSetLocation
}
