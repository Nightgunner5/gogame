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
	EntitySpawned
	EntityDespawned
)
