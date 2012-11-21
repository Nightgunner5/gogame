package network

type PacketID uint64

const (
	AttackerID PacketID = iota
	VictimID
	Amount
)
