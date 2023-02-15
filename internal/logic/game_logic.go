package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/emilebui/GBP_BE_echo/pkg/global"
	"github.com/emilebui/GBP_BE_echo/pkg/gstatus"
	"github.com/emilebui/GBP_BE_echo/pkg/helper"
	"github.com/redis/go-redis/v9"
	"time"
)

type GameLogic struct {
	r *redis.Client
}

func (g *GameLogic) Pick(c int) {
	fmt.Printf("Blah %d\n", c)
}

func CheckIfMoveValid(gs *GameState, cid string) error {

	// Check if 2 player are connected
	if gs.Turn == 0 || gs.Status == gstatus.WATTING {
		return errors.New(global.TextConfig["player_not_connected"])
	}

	if gs.Status == gstatus.ENDED {
		return errors.New(global.TextConfig["already_ended"])
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

	if _, ok = gs.BPMap[hid]; ok {
		return errors.New(global.TextConfig["already_chosen"])
	}

	return BPCore(gs, hid, cid, pick, r)
}

func BPCore(gs *GameState, hid int, cid string, pick bool, r *redis.Client) error {
	appendBPListPlayer(cid, gs, pick, hid)
	gs.BPMap[hid] = true
	gs.Turn = gs.Turn + 1

	ended := false
	if _, ok := TurnFormat[gs.Turn]; !ok {
		gs.Status = gstatus.ENDED
		ended = true
	} else {
		gs.Pick = CheckIfPickTurn(gs.Turn)
	}

	gs.PlayerTurn = GetPlayerTurn(gs)

	exp := 0 * time.Second
	if ended {
		exp = time.Duration(global.AfterGameExp) * time.Second
	}

	err := r.Set(context.Background(), gs.GameID, helper.Struct2String(gs), exp).Err()
	if err != nil {
		return errors.New(global.TextConfig["redis_error"])
	}

	result, err := r.Get(context.Background(), gs.GameID).Result()
	if err != nil {
		return errors.New(global.TextConfig["redis_error"])
	}

	helper.PublishRedis(&gstatus.ResponseMessage{
		Message: global.TextConfig["game_state_update"],
		Type:    gstatus.INFO,
		Data:    result,
	}, r, gs.GameID)

	if ended {
		helper.PublishRedis(&gstatus.ResponseMessage{
			Message: global.TextConfig["game_ended"],
			Type:    gstatus.INFO,
			Data:    result,
		}, r, gs.GameID)
	}
	return nil
}

func Chat(r *redis.Client, gs *GameState, mr *MoveRequest, cid string) error {
	return nil
}

func appendBPListPlayer(cid string, gs *GameState, pick bool, hid int) {
	if cid == gs.Player1.CID {
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
