package effect

import (
	"fmt"
	"github.com/Nightgunner5/gogame/entity"
)

type AbsorbDamage struct {
	Amount float64
}

func (a *AbsorbDamage) OnTakeDamage(amount *float64, attacker, victim entity.Entity) {
	if a.Amount > 0 {
		if *amount <= 0 {
			return
		}

		if a.Amount < *amount {
			*amount -= a.Amount
			a.Amount = 0
		} else {
			a.Amount -= *amount
			*amount = 0
		}
	} else if a.Amount < 0 {
		if *amount >= 0 {
			return
		}

		if a.Amount < -*amount {
			*amount += a.Amount
			a.Amount = 0
		} else {
			a.Amount += *amount
			*amount = 0
		}
	}
}

func (a *AbsorbDamage) String() string {
	if a.Amount > 0 {
		return fmt.Sprintf("Absorbs %d damage.", int(a.Amount))
	}
	if a.Amount < 0 {
		return fmt.Sprintf("Absorbs %d healing.", int(-a.Amount))
	}
	return ""
}

func (a *AbsorbDamage) effect() {}

type ReduceDamage struct {
	Fraction float64
}

func (r *ReduceDamage) assureValidFraction() {
	if r.Fraction > 1 {
		r.Fraction = 1
	}
	if r.Fraction < -1 {
		r.Fraction = -1
	}
}

func (r *ReduceDamage) OnTakeDamage(amount *float64, attacker, victim entity.Entity) {
	r.assureValidFraction()
	if r.Fraction > 0 && *amount > 0 {
		*amount *= r.Fraction
	} else if r.Fraction < 0 && *amount < 0 {
		*amount *= -r.Fraction
	}
}

func (r *ReduceDamage) String() string {
	r.assureValidFraction()
	if r.Fraction > 0 {
		return fmt.Sprintf("Reduces damage by %d%%.", int(r.Fraction*100))
	}
	if r.Fraction < 0 {
		return fmt.Sprintf("Reduces healing by %d%%.", int(-r.Fraction*100))
	}
	return ""
}

func (r *ReduceDamage) effect() {}
