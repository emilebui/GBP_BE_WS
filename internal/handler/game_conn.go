package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/emilebui/GBP_BE_WS/internal/logic"
	"github.com/emilebui/GBP_BE_WS/pkg/global"
	"github.com/emilebui/GBP_BE_WS/pkg/gstatus"
	"github.com/emilebui/GBP_BE_WS/pkg/helper"
	"time"
)

func (s *WebSocketHandler) ConnectGame(gid string, player *logic.Player, setting *logic.GameSetting) error {

	gameState, err := logic.GetGameState(gid, s.redisConn)

	if err != nil {
		return errors.New(global.TextConfig["redis_error"])
	}

	if gameState == nil {
		err = s.createGame(gid, player, setting)
		if err != nil {
			return errors.New(global.TextConfig["redis_error"])
		}
		helper.PublishRedis(&gstatus.ResponseMessage{
			Message: global.TextConfig["waiting_player"],
			Type:    gstatus.VALID,
		}, s.redisConn, gid)
	} else {
		err = s.processingConnection(gameState, player)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *WebSocketHandler) createGame(gid string, player *logic.Player, setting *logic.GameSetting) error {
	gs := &logic.GameState{
		GameID:  gid,
		Player1: *player,
		Player2: logic.Player{
			CID:      "",
			Nickname: "",
		},
		Status:            gstatus.WATTING,
		ConnectionTracker: map[string]bool{player.CID: true},
		Settings:          *setting,
	}
	return s.redisConn.Set(context.Background(), gid, helper.Struct2String(gs), 0).Err()
}

func (s *WebSocketHandler) appendPlayer(gs *logic.GameState, player *logic.Player) error {
	gs.Player2 = *player
	gs.Status = gstatus.PLAYING
	logic.ShufflePlayer(gs)
	gs.Turn = 1
	gs.Pick = logic.CheckIfPickTurn(gs)
	gs.PlayerTurn = logic.GetPlayerTurn(gs)
	gs.BPMapP1 = map[int]bool{0: false}
	gs.BPMapP2 = map[int]bool{0: false}
	gs.ConnectionTracker[player.CID] = true
	return s.redisConn.Set(context.Background(), gs.GameID, helper.Struct2String(gs), 0).Err()
}

func (s *WebSocketHandler) processingConnection(gs *logic.GameState, player *logic.Player) error {

	// If player 2 join the game --> init new game, game_status == PLAYING
	if gs.Player2.CID == "" && player.CID != gs.Player1.CID {
		err := s.appendPlayer(gs, player)
		if err != nil {
			return err
		}
		helper.PublishRedis(&gstatus.ResponseMessage{
			Message: global.TextConfig["game_ready"],
			Type:    gstatus.INFO,
			Data:    helper.Struct2String(gs),
		}, s.redisConn, gs.GameID)
		return nil
	}

	return s.checkIfReconnect(gs, player)
}

func (s *WebSocketHandler) checkIfReconnect(gs *logic.GameState, player *logic.Player) error {

	val, ok := gs.ConnectionTracker[player.CID]
	if !ok {
		return errors.New(global.TextConfig["enough_player"])
	}

	if val {
		return errors.New(global.TextConfig["already_join"])
	}

	gs.ConnectionTracker[player.CID] = true
	s.determineGameStatusAfterConnecting(gs)

	err := s.redisConn.Set(context.Background(), gs.GameID, helper.Struct2String(gs), 0).Err()
	if err != nil {
		return err
	}

	truePlayer := gs.Player1
	if player.CID != gs.Player1.CID {
		truePlayer = gs.Player2
	}

	helper.PublishRedis(&gstatus.ResponseMessage{
		Message: fmt.Sprintf(global.TextConfig["game_resume"], truePlayer.Nickname),
		Type:    gstatus.RECONN,
		Data:    helper.Struct2String(gs),
	}, s.redisConn, gs.GameID)

	s.announceChat(gs.GameID, &truePlayer, fmt.Sprintf(global.TextConfig["game_resume"], truePlayer.Nickname))

	return nil
}

func (s *WebSocketHandler) determineGameStatusAfterConnecting(gs *logic.GameState) {

	// If player 1 disconnected and reconnect --> gstatus should be WAITING after p1 reconnect
	if len(gs.ConnectionTracker) == 1 {
		gs.Status = gstatus.WATTING
	} else {
		check := helper.CheckIfAllBoolValueInMapIs(gs.ConnectionTracker, true)

		// If Both player are connected then gstatus should be PLAYING
		if check && gs.Status != gstatus.ENDED {
			gs.Status = gstatus.PLAYING
		}
	}

}

func (s *WebSocketHandler) disconnect(gid string, player *logic.Player) {

	gs, _ := logic.GetGameState(gid, s.redisConn)
	s.disconnectLogic(gs, player.CID)

	helper.PublishRedis(&gstatus.ResponseMessage{
		Message: fmt.Sprintf(global.TextConfig["player_disconnected"], player.Nickname),
		Type:    gstatus.DISCON,
		Data:    helper.Struct2String(gs),
	}, s.redisConn, gid)

	truePlayer := gs.Player1
	if player.CID != gs.Player1.CID {
		truePlayer = gs.Player2
	}
	s.announceChat(gid, &truePlayer, fmt.Sprintf(global.TextConfig["un_watch"], truePlayer.Nickname))
}

func (s *WebSocketHandler) determineIfGameIsDeadOrNot(gs *logic.GameState) {
	if helper.CheckIfAllBoolValueInMapIs(gs.ConnectionTracker, false) {
		s.redisConn.Set(context.Background(), gs.GameID, helper.Struct2String(gs), time.Duration(global.AfterGameExp)*time.Second)
	} else {
		s.redisConn.Set(context.Background(), gs.GameID, helper.Struct2String(gs), 0)
	}
}

func (s *WebSocketHandler) disconnectLogic(gs *logic.GameState, cid string) {
	gs.ConnectionTracker[cid] = false
	if gs.Status != gstatus.ENDED {
		gs.Status = gstatus.HALT
	}
	s.determineIfGameIsDeadOrNot(gs)
}

func (s *WebSocketHandler) announceChat(gid string, player *logic.Player, message string) {
	chatInfo := &logic.ChatInfo{
		Message:  message,
		CID:      player.CID,
		Nickname: player.Nickname,
		JoinChat: true,
	}

	helper.PublishRedis(&gstatus.ResponseMessage{
		Message: message,
		Type:    gstatus.LOG,
		Data:    helper.Struct2String(chatInfo),
	}, s.redisConn, gid)
}
