package effect

import (
	"fmt"
	"github.com/Nightgunner5/gogame/entity"
	"math"
)

type scaleDamageDone float64

func ScaleDamageDone(scale float64) Effect {
	if scale <= 0 {
		scale = math.Nextafter(0, 1)
	}
	return scaleDamageDone(scale)
}

func ScaleHealingDone(scale float64) Effect {
	if scale <= 0 {
		scale = math.Nextafter(0, 1)
	}
	return scaleDamageDone(-scale)
}

func (scale scaleDamageDone) OnDoDamage(amount *float64, attacker, victim entity.Entity) {
	if scale > 0 && *amount > 0 {
		*amount *= float64(scale)
	} else if scale < 0 && *amount < 0 {
		*amount *= float64(-scale)
	}
}

func (scale scaleDamageDone) String() string {
	if scale > 0 && scale != 1 {
		if scale > 1 {
			return fmt.Sprintf("Increases damage done by %d%%.", int(scale*100-100))
		}
		return fmt.Sprintf("Decreases damage done by %d%%.", int(100-scale*100))
	}
	if scale < 0 && scale != -1 {
		if scale < -1 {
			return fmt.Sprintf("Increases healing done by %d%%.", int(-scale*100-100))
		}
		return fmt.Sprintf("Decreases healing done by %d%%.", int(100+scale*100))
	}
	return ""
}

func (scaleDamageDone) effect() {}
