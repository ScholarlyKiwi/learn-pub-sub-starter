package main

import (
	"fmt"

	"github.com/ScholarlyKiwi/learn-pub-sub-starter/internal/gamelogic"
	"github.com/ScholarlyKiwi/learn-pub-sub-starter/internal/pubsub"
)

func handlerWar(gs *gamelogic.GameState) func(gamelogic.RecognitionOfWar) pubsub.AckType {
	return func(row gamelogic.RecognitionOfWar) pubsub.AckType {
		defer fmt.Print("> ")
		switch outcome, _, _ := gs.HandleWar(row); outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.AckType_Nack_Requeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.AckType_Nack_Discard
		case gamelogic.WarOutcomeOpponentWon:
			return pubsub.AckType_Ack
		case gamelogic.WarOutcomeYouWon:
			return pubsub.AckType_Ack
		case gamelogic.WarOutcomeDraw:
			return pubsub.AckType_Ack
		default:
			return pubsub.AckType_Nack_Discard
		}
	}
}
