package network

type PacketID uint64

const (
	AttackerID PacketID = iota
	VictimID
	Amount
	HealthChange

	debugEcho

	EntityID
	ParentID
	Tag
	EntitySpawned
	EntityDespawned
	EntityPosition

	// Insert new default packet IDs here

	FirstUnusedPacketID
)
