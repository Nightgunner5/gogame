package message

type (
	Sender   chan<- Message
	Receiver <-chan Message
)
