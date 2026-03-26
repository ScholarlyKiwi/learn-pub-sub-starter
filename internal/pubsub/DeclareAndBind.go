package pubsub

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // SimpleQueueType is an "enum" type I made to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {

	var channel *amqp.Channel
	var queue amqp.Queue

	channel, err := conn.Channel()
	if err != nil {
		return channel, queue, fmt.Errorf("DeclareAndBind channel error: %v", err)
	}

	durable := queueType == QueueType_Durable
	transient := queueType == QueueType_Transient // Set to autoDelete and exclusive

	queue, err = channel.QueueDeclare(
		queueName,
		durable,
		transient,
		transient,
		false,
		nil)

	err = channel.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return channel, queue, fmt.Errorf("DeclareAndBind queuebind error: %v", err)
	}

	return channel, queue, nil
}
