package message

type Kind *string

func NewKind(name string) Kind {
	return &name
}
