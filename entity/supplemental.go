package entity

type Positioner interface {
	Position() (x, y, z float64)
}

type Healther interface {
	Health() float64
}

type TakeDamager interface {
	TakeDamage(amount float64, attacker Entity)
}

type Thinker interface {
	Think(delta float64)
}
