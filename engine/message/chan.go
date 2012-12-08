package message

type (
	Sender   chan<- Message
	Reciever <-chan Message
)
