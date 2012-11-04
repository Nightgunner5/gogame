package entity

type Positioner interface {
	Entity

	Position() (x, y, z float64)
}

type Healther interface {
	Entity

	Health() float64
}
