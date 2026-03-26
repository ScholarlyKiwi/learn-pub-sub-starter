package pubsub

type SimpleQueueType int

const (
	QueueType_Durable SimpleQueueType = iota
	QueueType_Transient
)
