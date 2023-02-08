package handler

import (
	"context"
	"errors"
	"github.com/emilebui/GBP_BE_echo/internal/logic"
	"github.com/emilebui/GBP_BE_echo/pkg/global"
	"github.com/emilebui/GBP_BE_echo/pkg/gstatus"
	"github.com/emilebui/GBP_BE_echo/pkg/helper"
	"log"
)

func (s *WebSocketHandler) ConnectGame(gid string, cid string) error {

	gameState, err := logic.GetGameState(gid, s.redisConn)

	if err != nil {
		log.Println(gid, cid, err)
		return errors.New(global.TextConfig["redis_error"])
	}

	if gameState == nil {
		err = s.createGame(gid, cid)
		if err != nil {
			log.Println(gid, cid, err)
			return errors.New(global.TextConfig["redis_error"])
		}
		helper.PublishRedis(&gstatus.ResponseMessage{
			Message: global.TextConfig["waiting_player"],
			Type:    gstatus.VALID,
		}, s.redisConn, gid)
	} else {
		err = s.appendPlayerIfNecessary(gameState, cid)
		if err != nil {
			return err
		}
		helper.PublishRedis(&gstatus.ResponseMessage{
			Message: global.TextConfig["game_ready"],
			Type:    gstatus.INFO,
			Data:    helper.Struct2String(gameState),
		}, s.redisConn, gid)
	}

	return nil
}

func (s *WebSocketHandler) createGame(gid string, cid string) error {
	gs := &logic.GameState{
		GameID:  gid,
		Player1: cid,
		Status:  gstatus.WATTING,
	}
	return s.redisConn.Set(context.Background(), gid, helper.Struct2String(gs), 0).Err()
}

func (s *WebSocketHandler) appendPlayerIfNecessary(gs *logic.GameState, cid string) error {

	gid := gs.GameID

	if gs.Player2 == "" && cid != gs.Player1 {
		gs.Player2 = cid
		gs.Status = gstatus.PLAYING
		gs.PlayerTurnMap = logic.ShufflePlayer(gs.Player1, gs.Player2)
		gs.Turn = 1
		gs.PlayerTurn = logic.GetPlayerTurn(gs)
		return s.redisConn.Set(context.Background(), gid, helper.Struct2String(gs), 0).Err()
	}

	if cid != gs.Player1 && cid != gs.Player2 {
		return errors.New(global.TextConfig["enough_player"])
	}

	return nil
}
