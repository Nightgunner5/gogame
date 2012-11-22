package entity

type (
	Positioner interface {
		Position() (x, y, z float64)
		Move(dx, dy, dz float64)

		positionArray() []float64
	}

	Healther interface {
		Health() float64

		TakeDamage(amount float64, attacker Entity)
	}

	Resourcer interface {
		Resource() float64

		UseResource(amount float64) bool
	}

	Targeter interface {
		Target() Entity

		SetTarget(Entity)
	}

	Thinker interface {
		// This method will be called on all Spawned Entities
		// approximately once per 100ms with the amount of time
		// since the previous call times the value of TimeScale.
		Think(Δtime float64)
	}
)
