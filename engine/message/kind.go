package message

type Kind *string

// Create a new Kind, which is used to identify Message objects. See
// Message.Kind for more information.
func NewKind(name string) Kind {
	return &name
}
