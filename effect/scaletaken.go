package effect

import (
	"fmt"
	"github.com/Nightgunner5/gogame/entity"
	"math"
)

// Scales all damage taken by [scale]. A scale of 1 does nothing.
func ScaleDamageTaken(scale float64) Effect {
	if scale <= 0 {
		scale = math.Nextafter(0, 1)
	}
	return scaleDamageTaken(scale)
}

// Scales all healing taken by [scale]. A scale of 1 does nothing.
func ScaleHealingTaken(scale float64) Effect {
	if scale <= 0 {
		scale = math.Nextafter(0, 1)
	}
	return scaleDamageTaken(-scale)
}

type scaleDamageTaken float64

func (scale scaleDamageTaken) OnTakeDamage(amount *float64, attacker, victim entity.Entity) {
	if scale > 0 && *amount > 0 {
		*amount *= float64(scale)
	} else if scale < 0 && *amount < 0 {
		*amount *= float64(-scale)
	}
}

func (scale scaleDamageTaken) String() string {
	if scale > 0 && scale != 1 {
		if scale > 1 {
			return fmt.Sprintf("Increases damage taken by %d%%.", int(scale*100-100))
		}
		return fmt.Sprintf("Decreases damage taken by %d%%.", int(100-scale*100))
	}
	if scale < 0 && scale != -1 {
		if scale < -1 {
			return fmt.Sprintf("Increases healing taken by %d%%.", int(-scale*100-100))
		}
		return fmt.Sprintf("Decreases healing taken by %d%%.", int(100+scale*100))
	}
	return ""
}

func (scaleDamageTaken) effect() {}
