package effect

import (
	"fmt"
	"github.com/Nightgunner5/gogame/entity"
)

// Absorbs a total of [amount] damage over any number of TakeDamage calls.
func AbsorbDamage(amount float64) Effect {
	if amount < 0 {
		amount = 0
	}
	return (*absorbDamage)(&amount)
}

// Absorbs a total of [amount] healing over any number of TakeDamage calls.
func AbsorbHealing(amount float64) Effect {
	if amount < 0 {
		amount = 0
	}
	amount = -amount
	return (*absorbDamage)(&amount)
}

type absorbDamage float64

func (absorb *absorbDamage) OnTakeDamage(amount *float64, attacker, victim entity.Entity) {
	if *absorb > 0 {
		if *amount <= 0 {
			return
		}

		if float64(*absorb) < *amount {
			*amount -= float64(*absorb)
			*absorb = 0
		} else {
			*absorb -= absorbDamage(*amount)
			*amount = 0
		}
	} else if *absorb < 0 {
		if *amount >= 0 {
			return
		}

		if float64(*absorb) < -*amount {
			*amount += float64(*absorb)
			*absorb = 0
		} else {
			*absorb += absorbDamage(*amount)
			*amount = 0
		}
	}
}

func (absorb absorbDamage) String() string {
	if absorb > 0 {
		return fmt.Sprintf("Absorbs %d damage.", int(absorb))
	}
	if absorb < 0 {
		return fmt.Sprintf("Absorbs %d healing.", int(-absorb))
	}
	return ""
}

func (*absorbDamage) effect() {}
