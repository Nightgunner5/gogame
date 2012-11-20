package effect

import (
	"fmt"
	"github.com/Nightgunner5/gogame/entity"
	"math"
)

type scaleDamageDone struct {
	Scale float64
}

func ScaleDamageDone(scale float64) Effect {
	if scale <= 0 {
		scale = math.Nextafter(0, 1)
	}
	return &scaleDamageDone{scale}
}

func ScaleHealingDone(scale float64) Effect {
	if scale <= 0 {
		scale = math.Nextafter(0, 1)
	}
	return &scaleDamageDone{-scale}
}

func (s *scaleDamageDone) OnDoDamage(amount *float64, attacker, victim entity.Entity) {
	if s.Scale > 0 && *amount > 0 {
		*amount *= s.Scale
	} else if s.Scale < 0 && *amount < 0 {
		*amount *= -s.Scale
	}
}

func (s *scaleDamageDone) String() string {
	if s.Scale > 0 && s.Scale != 1 {
		if s.Scale > 1 {
			return fmt.Sprintf("Increases damage done by %d%%.", int(s.Scale*100-100))
		}
		return fmt.Sprintf("Decreases damage done by %d%%.", int(100-s.Scale*100))
	}
	if s.Scale < 0 && s.Scale != -1 {
		if s.Scale < -1 {
			return fmt.Sprintf("Increases healing done by %d%%.", int(-s.Scale*100-100))
		}
		return fmt.Sprintf("Decreases healing done by %d%%.", int(100+s.Scale*100))
	}
	return ""
}

func (s *scaleDamageDone) effect() {}
