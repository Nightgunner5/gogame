package network

type PacketID uint64

const (
	debugEcho PacketID = iota

	EntityID
	OtherEntID
	Tag

	Amount
	ChangeHealth
	ChangeResource

	EntitySpawned
	EntityDespawned
	EntityPosition

	CastSpell
	TimeLeft
	TotalTime

	// Insert new default packet IDs here

	FirstUnusedPacketID
)
