package server

import (
	"github.com/Nightgunner5/gogame/shared/packet"
	"log"
)

func Dispatch(player *Player, msg packet.Packet) bool {
	switch {
	case msg.Location != nil:
		dx, dy := msg.Location.Coord.X, msg.Location.Coord.Y
		if (dx == 0 && (dy == -1 || dy == 1)) || (dy == 0 && (dx == -1 || dx == 1)) {
			player.Send <- *msg.Location
		}
	case msg.Chat != nil:
		msg.Chat.User = player.Name
		SendToAll <- packet.Packet{
			Chat: msg.Chat,
		}
	default:
		log.Printf("Client %q sent unknown packet %#v", player.id, msg)
		return false
	}
	return true
}
