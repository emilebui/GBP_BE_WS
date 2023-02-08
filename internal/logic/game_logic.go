package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/emilebui/GBP_BE_echo/pkg/global"
	"github.com/emilebui/GBP_BE_echo/pkg/gstatus"
	"github.com/emilebui/GBP_BE_echo/pkg/helper"
	"github.com/redis/go-redis/v9"
)

type GameLogic struct {
	r *redis.Client
}

func (g *GameLogic) Pick(c int) {
	fmt.Printf("Blah %d\n", c)
}

func CheckIfMoveValid(gs *GameState, cid string) error {

	// Check if 2 player are connected
	if gs.Turn == 0 || gs.Status != gstatus.PLAYING {
		return errors.New(global.TextConfig["player_not_connected"])
	}

	// Check if not your turn --> Ignore
	playerTurn := GetPlayerTurn(gs)
	if playerTurn != cid {
		return errors.New(global.TextConfig["not_your_turn"])
	}

	return nil
}

func Ban(r *redis.Client, gs *GameState, mr *MoveRequest, cid string) error {
	return BanPickLogic(r, gs, mr, cid, false)
}

func Pick(r *redis.Client, gs *GameState, mr *MoveRequest, cid string) error {
	return BanPickLogic(r, gs, mr, cid, true)
}

func BanPickLogic(r *redis.Client, gs *GameState, mr *MoveRequest, cid string, pick bool) error {

	if CheckIfPickTurn(gs.Turn) != pick {
		return errors.New(global.TextConfig["wrong_move"])
	}

	hid, ok := helper.Interface2Int(mr.Data)
	if !ok {
		return errors.New(global.TextConfig["invalid_data_type"])
	}

	appendBPListPlayer(cid, gs, pick, hid)
	gs.Turn = gs.Turn + 1

	err := r.Set(context.Background(), gs.GameID, helper.Struct2String(gs), 0).Err()
	if err != nil {
		return errors.New(global.TextConfig["redis_error"])
	}

	result, err := r.Get(context.Background(), gs.GameID).Result()

	helper.PublishRedis(&gstatus.ResponseMessage{
		Message: global.TextConfig["game_state_update"],
		Type:    gstatus.INFO,
		Data:    result,
	}, r, gs.GameID)

	return nil
}

func Chat(r *redis.Client, gs *GameState, mr *MoveRequest, cid string) error {
	return nil
}

func appendBPListPlayer(cid string, gs *GameState, pick bool, hid int) {
	if cid == gs.Player1 {
		if pick {
			gs.Board.P1Pick = append(gs.Board.P1Pick, hid)
		} else {
			gs.Board.P1Ban = append(gs.Board.P1Ban, hid)
		}
	} else {
		if pick {
			gs.Board.P2Pick = append(gs.Board.P2Pick, hid)
		} else {
			gs.Board.P2Ban = append(gs.Board.P2Ban, hid)
		}
	}
}
