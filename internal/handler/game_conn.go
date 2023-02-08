package handler

import (
	"context"
	"errors"
	"github.com/emilebui/GBP_BE_echo/internal/logic"
	"github.com/emilebui/GBP_BE_echo/pkg/gstatus"
	"github.com/emilebui/GBP_BE_echo/pkg/helper"
	"log"
)

func (s *WebSocketHandler) ConnectGame(gid string, cid string) error {

	rawData, err := s.redisConn.Get(context.Background(), gid).Result()
	if err != nil || len(rawData) == 0 {
		err = s.createGame(gid, cid)
		if err != nil {
			log.Println(gid, cid, err)
			return errors.New(s.responseMap["redis_error"])
		}
		helper.PublishRedis(&gstatus.ResponseMessage{
			Message: s.responseMap["waiting_player"],
			Type:    gstatus.VALID,
		}, s.redisConn, gid)

	} else {
		// Append Game if necessary
		gameState := new(logic.GameState)

		err = helper.StringToStruct(rawData, gameState)
		if err != nil {
			return errors.New(s.responseMap["redis_data_error"])
		}

		err = s.appendPlayerIfNecessary(gameState, cid)
		if err != nil {
			return err
		}
		helper.PublishRedis(&gstatus.ResponseMessage{
			Message: s.responseMap["game_ready"],
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
	_, err := s.redisConn.Set(context.Background(), gid, helper.Struct2String(gs), 0).Result()
	return err
}

func (s *WebSocketHandler) appendPlayerIfNecessary(gs *logic.GameState, cid string) error {

	gid := gs.GameID

	if gs.Player2 == "" && cid != gs.Player1 {
		gs.Player2 = cid
		gs.Status = gstatus.PLAYING
		gs.PlayerTurn = helper.RandomInList([]string{gs.Player1, gs.Player2})
		gs.Turn = 0
		s.redisConn.Set(context.Background(), gid, helper.Struct2String(gs), 0)
		return nil
	}

	if cid != gs.Player1 && cid != gs.Player2 {
		return errors.New(s.responseMap["enough_player"])
	}

	return nil
}
