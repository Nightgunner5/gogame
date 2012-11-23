package network

type PacketID uint64

const (
	debugEcho PacketID = iota

	EntityID
	OtherEntID
	EntityTag

	Amount
	ChangeHealth
	ChangeResource

	EntitySpawned
	EntityDespawned
	EntityPosition

	// Insert new default packet IDs here

	FirstUnusedPacketID
)
