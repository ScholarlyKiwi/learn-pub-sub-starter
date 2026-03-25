package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/ScholarlyKiwi/learn-pub-sub-starter/internal/pubsub"
	"github.com/ScholarlyKiwi/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	const connectString = "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connectString)
	if err != nil {
		fmt.Printf("Error opening connection: %v\n", err)
	}
	defer connection.Close()

	fmt.Println("Peril Server successfully start!")

	perilChan, err := connection.Channel()

	err = pubsub.PublishJSON(perilChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
		IsPaused: true,
	})

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	fmt.Println("Peril Server shutting done.")
}
