package effect

type Primitive interface {
	isSameType(Primitive) bool
	// Only to be called if isSameType returns true for the same argument.
	combine(Primitive) Primitive

	effectThink(float64)
	String() string
}
