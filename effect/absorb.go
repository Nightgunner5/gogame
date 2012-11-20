package effect

import (
	"fmt"
	"github.com/Nightgunner5/gogame/entity"
	"math"
)

func AbsorbDamage(amount float64) Effect {
	if amount < 0 {
		amount = 0
	}
	return &absorbDamage{amount}
}

func AbsorbHealing(amount float64) Effect {
	if amount < 0 {
		amount = 0
	}
	return &absorbDamage{-amount}
}

type absorbDamage struct {
	Amount float64
}

func (a *absorbDamage) OnTakeDamage(amount *float64, attacker, victim entity.Entity) {
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

func (a *absorbDamage) String() string {
	if a.Amount > 0 {
		return fmt.Sprintf("Absorbs %d damage.", int(a.Amount))
	}
	if a.Amount < 0 {
		return fmt.Sprintf("Absorbs %d healing.", int(-a.Amount))
	}
	return ""
}

func (a *absorbDamage) effect() {}

func ScaleDamageTaken(scale float64) Effect {
	if scale <= 0 {
		scale = math.Nextafter(0, 1)
	}
	return &scaleDamageTaken{scale}
}

func ScaleHealingTaken(scale float64) Effect {
	if scale <= 0 {
		scale = math.Nextafter(0, 1)
	}
	return &scaleDamageTaken{-scale}
}

type scaleDamageTaken struct {
	Scale float64
}

func (s *scaleDamageTaken) OnTakeDamage(amount *float64, attacker, victim entity.Entity) {
	if s.Scale > 0 && *amount > 0 {
		*amount *= s.Scale
	} else if s.Scale < 0 && *amount < 0 {
		*amount *= -s.Scale
	}
}

func (s *scaleDamageTaken) String() string {
	if s.Scale > 0 && s.Scale != 1 {
		if s.Scale > 1 {
			return fmt.Sprintf("Increases damage taken by %d%%.", int(s.Scale*100-100))
		}
		return fmt.Sprintf("Decreases damage taken by %d%%.", int(100-s.Scale*100))
	}
	if s.Scale < 0 && s.Scale != -1 {
		if s.Scale < -1 {
			return fmt.Sprintf("Increases healing taken by %d%%.", int(-s.Scale*100-100))
		}
		return fmt.Sprintf("Decreases healing taken by %d%%.", int(100+s.Scale*100))
	}
	return ""
}

func (s *scaleDamageTaken) effect() {}
