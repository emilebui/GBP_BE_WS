package broker

import (
	"errors"
	"github.com/emilebui/GBP_BE_WS/internal/logic"
	"github.com/emilebui/GBP_BE_WS/pkg/global"
	"github.com/emilebui/GBP_BE_WS/pkg/helper"
)

/* This package is to handle websocket messages into game logic calls*/

// Logic to call map

func Ban(gl *logic.GameLogic, mr *logic.MoveRequest) error {
	return gl.Ban(mr)
}

func Pick(gl *logic.GameLogic, mr *logic.MoveRequest) error {
	return gl.Pick(mr)
}

func Chat(gl *logic.GameLogic, mr *logic.MoveRequest) error {
	return gl.Chat(mr)
}

func GetGameState(gl *logic.GameLogic, mr *logic.MoveRequest) error {
	return gl.GetGameState(mr)
}

var funcMap = map[string]func(gl *logic.GameLogic, mr *logic.MoveRequest) error{
	"BAN":       Ban,
	"PICK":      Pick,
	"CHAT":      Chat,
	"GET_STATE": GetGameState,
}

func ProcessMessage(msg []byte, gl *logic.GameLogic) error {

	// Parse message
	mr, err := ParseMessage(msg)
	if err != nil {
		return errors.New(global.TextConfig["data_input_error"])
	}

	// Check if MR is valid
	_, ok := funcMap[mr.Call]
	if !ok {
		return errors.New(global.TextConfig["data_input_error"])
	}

	// Move
	err = funcMap[mr.Call](gl, mr)
	return err
}

func ParseMessage(msg []byte) (*logic.MoveRequest, error) {
	mr := new(logic.MoveRequest)
	err := helper.BytesToStruct(msg, mr)
	return mr, err
}
