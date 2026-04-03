package pubsub

type AckType int

const (
	AckType_Ack AckType = iota
	AckType_Nack_Discard
	AckType_Nack_Requeue
)
