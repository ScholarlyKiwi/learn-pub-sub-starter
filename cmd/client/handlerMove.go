package main

import (
	"fmt"

	"github.com/ScholarlyKiwi/learn-pub-sub-starter/internal/gamelogic"
	"github.com/ScholarlyKiwi/learn-pub-sub-starter/internal/pubsub"
	"github.com/ScholarlyKiwi/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerMove(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.ArmyMove) pubsub.AckType {
	return func(mv gamelogic.ArmyMove) pubsub.AckType {
		defer fmt.Print("> ")
		outcome := gs.HandleMove(mv)
		switch outcome {
		case gamelogic.MoveOutComeSafe:
			return pubsub.AckType_Ack

		case gamelogic.MoveOutcomeMakeWar:
			war_routing_key := fmt.Sprintf("%v.%v", routing.WarRecognitionsPrefix, gs.GetUsername())
			err := pubsub.PublishJSON(ch, routing.ExchangePerilTopic, war_routing_key, gamelogic.RecognitionOfWar{
				Attacker: mv.Player,
				Defender: gs.GetPlayerSnap(),
			})

			if err != nil {
				fmt.Printf("Error publishing war recognition: %v\n", err)
				return pubsub.AckType_Nack_Requeue
			}
			return pubsub.AckType_Ack

		case gamelogic.MoveOutcomeSamePlayer:
			return pubsub.AckType_Ack
		}
		fmt.Printf("Unknown move outcome: %v\n", outcome)
		return pubsub.AckType_Nack_Discard
	}
}
