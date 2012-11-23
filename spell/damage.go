package spell

import "github.com/Nightgunner5/gogame/entity"

func DamageSpell(delay, damage float64, caster, target entity.Entity, interruptable bool, tag string) Spell {
	if caster == nil || target == nil || damage == 0 || delay < 0 {
		return nil
	}

	if _, ok := target.(entity.Healther); !ok {
		return nil
	}

	return &BasicSpell{
		CastTime:      delay,
		Interruptable: interruptable,
		Caster_:       caster.ID(),
		Target_:       target.ID(),
		Action: func(target, caster entity.Entity) {
			target.(entity.Healther).TakeDamage(damage, caster)
		},
		Tag_: tag,
	}
}

func HealingSpell(delay, healing float64, caster, target entity.Entity, interruptable bool, tag string) Spell {
	return DamageSpell(delay, -healing, caster, target, interruptable, tag)
}

func DamageOverTimeSpell(duration, damagePerSecond float64, caster, target entity.Entity, interruptable bool, tag string) Spell {
	if caster == nil || target == nil || damagePerSecond == 0 || duration < 0 {
		return nil
	}

	if _, ok := target.(entity.Healther); !ok {
		return nil
	}

	return &ChanneledSpell{
		CastTime:      duration,
		Interruptable: interruptable,
		Caster_:       caster.ID(),
		Target_:       target.ID(),
		Action: func(target, caster entity.Entity, Δtime float64) {
			target.(entity.Healther).TakeDamage(damagePerSecond*Δtime, caster)
		},
		Tag_: tag,
	}
}

func HealingOverTimeSpell(duration, healingPerSecond float64, caster, target entity.Entity, interruptable bool, tag string) Spell {
	return DamageOverTimeSpell(duration, -healingPerSecond, caster, target, interruptable, tag)
}
