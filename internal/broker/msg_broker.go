package broker

import "github.com/emilebui/GBP_BE_echo/internal/logic"

/* This package is to handle websocket messages into game logic calls*/

type MessageBroker struct {
	logic     *logic.GameLogic
	GameState *logic.GameState
	ClientID  string
}
