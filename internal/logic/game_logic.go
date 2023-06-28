package logic

import (
	"context"
	"errors"
	"github.com/emilebui/GBP_BE_WS/pkg/global"
	"github.com/emilebui/GBP_BE_WS/pkg/gstatus"
	"github.com/emilebui/GBP_BE_WS/pkg/helper"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"time"
)

type GameLogic struct {
	r      *redis.Client
	c      *websocket.Conn
	GID    string
	Player *Player
}

func NewGameLogic(r *redis.Client, c *websocket.Conn, gid string, player *Player) *GameLogic {
	return &GameLogic{
		r, c, gid, player,
	}
}

func (g *GameLogic) Ban(mr *MoveRequest) error {

	// Get GameState
	gs, err := GetGameState(g.GID, g.r)
	if err != nil {
		return err
	}

	// Check valid move
	err = g.checkIfMoveValid(gs)
	if err != nil {
		return err
	}

	return g.banPickLogic(gs, mr, false)
}

func (g *GameLogic) Pick(mr *MoveRequest) error {
	// Get GameState
	gs, err := GetGameState(g.GID, g.r)
	if err != nil {
		return err
	}

	// Check valid move
	err = g.checkIfMoveValid(gs)
	if err != nil {
		return err
	}

	return g.banPickLogic(gs, mr, true)
}

func (g *GameLogic) GetGameState(_ *MoveRequest) error {

	res := &gstatus.ResponseMessage{}

	gs, err := GetGameState(g.GID, g.r)
	if err != nil {
		return err
	}

	if gs.Status == gstatus.WATTING {
		res.Message = global.TextConfig["waiting_player"]
		res.Type = gstatus.VALID
	} else {
		res.Message = global.TextConfig["game_state_return"]
		res.Type = gstatus.INFO
		res.Data = helper.Struct2String(gs)
	}

	helper.ResponseWS(res, g.c)
	return nil
}

func (g *GameLogic) Chat(mr *MoveRequest) error {

	msg, ok := mr.Data.(string)
	if !ok {
		return errors.New(global.TextConfig["invalid_data_type"])
	}

	// Get GameState
	gs, err := GetGameState(g.GID, g.r)
	if err != nil {
		return err
	}

	nickname := g.Player.Nickname
	if g.Player.CID == gs.Player1.CID {
		nickname = gs.Player1.Nickname
	}

	if g.Player.CID == gs.Player2.CID {
		nickname = gs.Player2.Nickname
	}

	chatInfo := &ChatInfo{
		Message:  msg,
		CID:      g.Player.CID,
		Nickname: nickname,
		JoinChat: false,
	}

	helper.PublishRedis(&gstatus.ResponseMessage{
		Message: global.TextConfig["game_state_update"],
		Type:    gstatus.LOG,
		Data:    helper.Struct2String(chatInfo),
	}, g.r, g.GID)

	return nil
}

func (g *GameLogic) checkIfMoveValid(gs *GameState) error {

	// Check if 2 player are connected
	if gs.Turn == 0 || gs.Status == gstatus.WATTING {
		return errors.New(global.TextConfig["player_not_connected"])
	}

	if gs.Status == gstatus.ENDED {
		return errors.New(global.TextConfig["already_ended"])
	}

	// Check if not your turn --> Ignore
	playerTurn := GetPlayerTurn(gs)
	if playerTurn != g.Player.CID {
		return errors.New(global.TextConfig["not_your_turn"])
	}

	return nil
}

func (g *GameLogic) banPickLogic(gs *GameState, mr *MoveRequest, pick bool) error {

	if CheckIfPickTurn(gs) != pick {
		return errors.New(global.TextConfig["wrong_move"])
	}

	hid, ok := helper.Interface2Int(mr.Data)
	if !ok {
		return errors.New(global.TextConfig["invalid_data_type"])
	}

	if g.checkAlreadyChosen(gs, hid) {
		return errors.New(global.TextConfig["already_chosen"])
	}

	return g.bpCore(gs, hid, pick)
}

func (g *GameLogic) bpCore(gs *GameState, hid int, pick bool) error {
	g.appendBPListPlayer(gs, pick, hid)
	g.changeBPList(gs, hid, pick)
	gs.Turn = gs.Turn + 1

	ended := false
	if !CheckIfTurnValid(gs) {
		gs.Status = gstatus.ENDED
		ended = true
	} else {
		gs.Pick = CheckIfPickTurn(gs)
	}

	gs.PlayerTurn = GetPlayerTurn(gs)

	exp := 0 * time.Second
	if ended {
		exp = time.Duration(global.AfterGameExp) * time.Second
	}

	err := g.r.Set(context.Background(), gs.GameID, helper.Struct2String(gs), exp).Err()
	if err != nil {
		return errors.New(global.TextConfig["redis_error"])
	}

	result, err := g.r.Get(context.Background(), gs.GameID).Result()
	if err != nil {
		return errors.New(global.TextConfig["redis_error"])
	}

	helper.PublishRedis(&gstatus.ResponseMessage{
		Message: global.TextConfig["game_state_update"],
		Type:    gstatus.INFO,
		Data:    result,
	}, g.r, gs.GameID)

	if ended {
		helper.PublishRedis(&gstatus.ResponseMessage{
			Message: global.TextConfig["game_ended"],
			Type:    gstatus.INFO,
			Data:    result,
		}, g.r, gs.GameID)
	}
	return nil
}

func (g *GameLogic) appendBPListPlayer(gs *GameState, pick bool, hid int) {
	if g.Player.CID == gs.Player1.CID {
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

func (g *GameLogic) changeBPList(gs *GameState, hid int, pick bool) {

	if pick && gs.Settings.Casual {
		if g.Player.CID == gs.Player1.CID {
			gs.BPMapP1[hid] = true
		} else {
			gs.BPMapP2[hid] = true
		}
	} else {
		gs.BPMapP1[hid] = true
		gs.BPMapP2[hid] = true
	}
}

func (g *GameLogic) checkAlreadyChosen(gs *GameState, hid int) bool {

	if g.Player.CID == gs.Player1.CID {
		if _, ok := gs.BPMapP1[hid]; ok {
			return true
		}
	} else {
		if _, ok := gs.BPMapP2[hid]; ok {
			return true
		}
	}

	return false
}
