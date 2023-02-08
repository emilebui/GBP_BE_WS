package handler

import (
	"context"
	"errors"
	"github.com/emilebui/GBP_BE_echo/internal/logic"
	"github.com/emilebui/GBP_BE_echo/pkg/gstatus"
	"github.com/emilebui/GBP_BE_echo/pkg/helper"
)

func (s *WebSocketHandler) ConnectGame(gid string, cid string) error {

	rawData, err := s.redisConn.HGetAll(context.Background(), gid).Result()
	if err != nil {
		s.createGame(gid, cid)
	} else {
		// Append Game if necessary
		gameState := new(logic.GameState)
		err = helper.Struct2Struct(rawData, gameState)
		if err != nil {
			return errors.New("decode error from redis - connect game")
		}
		err = s.appendPlayerIfNecessary(gameState, cid)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *WebSocketHandler) createGame(gid string, cid string) {
	gs := &logic.GameState{
		GameID:  gid,
		Player1: cid,
		Status:  gstatus.WATTING,
	}
	s.redisConn.HMSet(context.Background(), gid, gs)
}

func (s *WebSocketHandler) appendPlayerIfNecessary(gs *logic.GameState, cid string) error {

	gid := gs.GameID

	if gs.Player2 == "" {
		gs.Player2 = cid
		gs.Status = gstatus.PLAYING
		gs.PlayerTurn = helper.RandomInList([]string{gs.Player1, gs.Player2})
		gs.Turn = 0
		s.redisConn.HMSet(context.Background(), gid, gs)
		return nil
	}

	if cid != gs.Player1 && cid != gs.Player2 {
		return errors.New("there is enough player already")
	}

	return nil
}
