package message

type Message interface {
	// The Kind of message, which should be a value obtained from NewKind
	// at initialization. The concrete type of a Message and the return
	// value of its Kind method must be a one-to-one ratio.
	Kind() Kind
}
