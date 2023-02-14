package broker

import (
	"errors"
	"github.com/emilebui/GBP_BE_echo/internal/logic"
	"github.com/emilebui/GBP_BE_echo/pkg/global"
	"github.com/emilebui/GBP_BE_echo/pkg/gstatus"
	"github.com/emilebui/GBP_BE_echo/pkg/helper"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

/* This package is to handle websocket messages into game logic calls*/

var funcMap = map[string]func(r *redis.Client, gs *logic.GameState, mr *logic.MoveRequest, cid string) error{
	"BAN":       logic.Ban,
	"PICK":      logic.Pick,
	"CHAT":      logic.Chat,
	"GET_STATE": nil,
}

func ProcessMessage(r *redis.Client, msg []byte, gid string, cid string, c *websocket.Conn) error {

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

	// Get GameState
	gs, err := logic.GetGameState(gid, r)
	if err != nil {
		return err
	}

	if mr.Call == "GET_STATE" {
		returnCurrentGameState(gs, c)
		return nil
	}

	// Check valid move
	err = logic.CheckIfMoveValid(gs, cid)
	if err != nil {
		return err
	}

	// Move
	err = funcMap[mr.Call](r, gs, mr, cid)
	return err
}

func returnCurrentGameState(gs *logic.GameState, c *websocket.Conn) {

	res := &gstatus.ResponseMessage{}

	if gs.Status == gstatus.WATTING {
		res.Message = global.TextConfig["waiting_player"]
		res.Type = gstatus.VALID
	} else {
		res.Message = global.TextConfig["game_state_return"]
		res.Type = gstatus.INFO
		res.Data = helper.Struct2String(gs)
	}

	helper.ResponseWS(res, c)
}

func ParseMessage(msg []byte) (*logic.MoveRequest, error) {
	mr := new(logic.MoveRequest)
	err := helper.BytesToStruct(msg, mr)
	return mr, err
}
