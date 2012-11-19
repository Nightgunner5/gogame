package effect

type Effect interface {
	effect()

	String() string
}
