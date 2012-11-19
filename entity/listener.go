package entity

type DamageListener interface {
	OnTakeDamage(amount *float64, attacker, victim Entity)
}
