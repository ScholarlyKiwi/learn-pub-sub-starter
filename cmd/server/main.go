package main

import (
	"fmt"
	"log"

	"github.com/ScholarlyKiwi/learn-pub-sub-starter/internal/gamelogic"
	"github.com/ScholarlyKiwi/learn-pub-sub-starter/internal/pubsub"
	"github.com/ScholarlyKiwi/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	const connectString = "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connectString)
	if err != nil {
		log.Fatalf("Error opening connection: %v\n", err)
	}
	defer connection.Close()

	fmt.Println("Peril Server successfully started!")

	_, queue, err := pubsub.DeclareAndBind(
		connection,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".*",
		pubsub.QueueType_Durable,
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	gamelogic.PrintServerHelp()

	perilChan, err := connection.Channel()
	if err != nil {
		log.Fatalf("Error pausing: %v\n", err)
		return
	}
	defer perilChan.Close()

commandLoop:
	for {
		commands := gamelogic.GetInput()
		if len(commands) > 0 {
			command := commands[0]
			switch command {
			case "pause":
				fmt.Println("Publishing paused game state")
				err = pubsub.PublishJSON(perilChan,
					routing.ExchangePerilDirect,
					routing.PauseKey,
					routing.PlayingState{
						IsPaused: true,
					})
				if err != nil {
					log.Printf("Could not publish pause: %v\n", err)
					return
				}
			case "resume":
				fmt.Println("Sending resume message")
				err = pubsub.PublishJSON(perilChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
					IsPaused: false,
				})
				if err != nil {
					log.Printf("Could not publish resume: %v\n", err)
				}
			case "quit":
				log.Println("Goodbye!")
				break commandLoop
			default:
				fmt.Printf("Command %v not understand.\n", command)
			}
		}
	}
}
