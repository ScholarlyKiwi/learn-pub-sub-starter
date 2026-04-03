package pubsub

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T) AckType,
) error {

	channel, _, err := DeclareAndBind(conn, exchange, queueName, key, QueueType_Durable)
	if err != nil {
		return fmt.Errorf("SubscribeJSON error: %v", err)
	}

	deliveryChan, err := channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("SubscribeJSON dilvery channelerror: %v", err)
	}
	go func() {
		for msg := range deliveryChan {
			var data T
			err := json.Unmarshal(msg.Body, &data)
			if err != nil {
				log.Printf("SubscribeJSON Unmarshalerror: %v", err)
				continue
			}
			ackresult := handler(data)

			switch ackresult {
			case AckType_Ack:
				log.Println("Message acknowledged")
				err = msg.Ack(false)
				if err != nil {
					log.Printf("SubscribeJSON Ack error: %v", err)
				}
			case AckType_Nack_Discard:
				log.Println("Message rejected")
				err = msg.Nack(false, false)
				if err != nil {
					log.Printf("SubscribeJSON Nackerror: %v", err)
				}
			case AckType_Nack_Requeue:
				log.Println("Message requeued")
				err = msg.Nack(false, true)
				if err != nil {
					log.Printf("SubscribeJSON Requeue error: %v", err)
				}
			}
		}
	}()

	return nil
}
