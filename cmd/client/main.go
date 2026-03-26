package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/ScholarlyKiwi/learn-pub-sub-starter/internal/gamelogic"
	"github.com/ScholarlyKiwi/learn-pub-sub-starter/internal/pubsub"
	"github.com/ScholarlyKiwi/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")

	const connectString = "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connectString)
	if err != nil {
		fmt.Printf("Error opening connection: %v\n", err)
	}
	defer connection.Close()

	fmt.Println("Peril Client successfully start!")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Println(err)
	}
	queueName := fmt.Sprintf("%v.%v", routing.PauseKey, username)
	pauseChannel, pauseQueue, err := pubsub.DeclareAndBind(connection, "peril_direct", queueName, routing.PauseKey, pubsub.QueueType_Transient)
	if err != nil {
		fmt.Println(err)
	}
	defer pauseChannel.Close()

	err = fmt.Errorf("%v", pauseQueue.Name)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	fmt.Println("Peril client exit.")
}
