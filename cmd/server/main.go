package main

import (
	"fmt"

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
		fmt.Printf("Error opening connection: %v\n", err)
	}
	defer connection.Close()

	fmt.Println("Peril Server successfully started!")

	gamelogic.PrintServerHelp()

	perilChan, err := connection.Channel()
	if err != nil {
		fmt.Printf("Error pausing: %v\n", err)
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
				fmt.Println("Sending pause message")
				err = pubsub.PublishJSON(perilChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
					IsPaused: true,
				})
				if err != nil {
					fmt.Printf("Error pausing game: %v\n", err)
					return
				}
			case "resume":
				fmt.Println("Sending resume message")
				err = pubsub.PublishJSON(perilChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
					IsPaused: false,
				})
				if err != nil {
					fmt.Printf("Error pausing game: %v\n", err)
					return
				}
			case "quit":
				fmt.Println("Goodbye!")
				break commandLoop
			default:
				fmt.Printf("Command %v not understand.", command)
			}
		}
	}

	fmt.Println("Peril Server shutting done.")
}
