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
	fmt.Println("Starting Peril client...")

	const connectString = "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connectString)
	if err != nil {
		log.Fatalf("Error opening connection: %v\n", err)
	}
	defer connection.Close()

	fmt.Println("Peril Client successfully start!")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("could not get username: %v", err)
	}

	perilState := gamelogic.NewGameState(username)

	pause_queue := fmt.Sprintf("%v.%v", routing.PauseKey, username)
	err = pubsub.SubscribeJSON(connection, routing.ExchangePerilDirect, pause_queue, routing.PauseKey, pubsub.QueueType_Transient, handlerPause(perilState))
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}

	armyMove_queue := fmt.Sprintf("%v.%v", routing.ArmyMovesPrefix, username)
	armyMove_key := fmt.Sprintf("%v.*", routing.ArmyMovesPrefix)
	err = pubsub.SubscribeJSON(connection, routing.ExchangePerilTopic, armyMove_queue, armyMove_key, pubsub.QueueType_Transient, handlerMove(perilState))
	if err != nil {
		log.Fatalf("could not subscribe to army moves: %v", err)
	}
	channel, err := connection.Channel()
	if err != nil {
		log.Fatalf("could not create channel: %v", err)
	}
	defer channel.Close()

	for {

		commands := gamelogic.GetInput()
		if len(commands) > 0 {
			switch commands[0] {
			case "spawn":
				err = perilState.CommandSpawn(commands)
				if err != nil {
					fmt.Println(err)
					continue
				}
			case "move":
				moveResponse, err := perilState.CommandMove(commands)
				err = pubsub.PublishJSON(channel, routing.ExchangePerilTopic, armyMove_queue, moveResponse)
				if err != nil {
					log.Println("move published")
				}
				fmt.Println(moveResponse)
				if err != nil {
					fmt.Println(err)
					continue
				}
			case "status":
				perilState.CommandStatus()
			case "help":
				gamelogic.PrintClientHelp()
			case "spam":
				fmt.Println("Spamming not allowed yet!")
			case "quit":
				gamelogic.PrintQuit()
				return
			default:
				fmt.Printf("Unknown command %v", commands[0])
			}
		}
	}
}
