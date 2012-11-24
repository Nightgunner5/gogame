package effect

import (
	"github.com/Nightgunner5/gogame/entity"
	"github.com/Nightgunner5/gogame/network"
)

func broadcastPacket(ent entity.EntityID, effects []*Effect) {
	network.Broadcast(toPacket(ent, effects), false)
}

func toPacket(ent entity.EntityID, effects []*Effect) network.Packet {
	encoded := make([]map[string]interface{}, len(effects))

	for i, e := range effects {
		encoded[i] = make(map[string]interface{})
		encoded[i]["desc"] = e.String()
	}

	return network.NewPacket(network.EntityEffects).
		Set(network.EntityID, ent).
		Set(network.EntityEffects, encoded)
}
