package main

import (
	"fmt"
	"log"

	"github.com/ScholarlyKiwi/learn-pub-sub-starter/internal/gamelogic"
	"github.com/ScholarlyKiwi/learn-pub-sub-starter/internal/pubsub"
)

func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) pubsub.AckType {
	return func(mv gamelogic.ArmyMove) pubsub.AckType {
		defer fmt.Print("> ")
		outcome := gs.HandleMove(mv)
		switch outcome {
		case gamelogic.MoveOutComeSafe:
		case gamelogic.MoveOutcomeMakeWar:
			log.Printf("Move successful: %v\n", outcome)
			return pubsub.AckType_Ack
		case gamelogic.MoveOutcomeSamePlayer:
			log.Printf("Unable to move, same player: %v\n,", outcome)
			return pubsub.AckType_Nack_Discard
		}
		log.Printf("Unknown move outcome: %v\n", outcome)
		return pubsub.AckType_Nack_Discard
	}
}
