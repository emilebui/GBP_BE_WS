package logic

import (
	"context"
	"errors"
	"github.com/emilebui/GBP_BE_WS/pkg/global"
	"github.com/emilebui/GBP_BE_WS/pkg/helper"
	"github.com/redis/go-redis/v9"
	"log"
)

type GameState struct {
	GameID            string          `json:"game_id"`
	Player1           Player          `json:"player_1"`
	Player2           Player          `json:"player_2"`
	Turn              int             `json:"turn"`
	Status            int             `json:"status"`
	PlayerTurn        string          `json:"player_turn"`
	ConnectionTracker map[string]bool `json:"connection_tracker"`
	Pick              bool            `json:"pick"`
	BPMap             map[int]bool    `json:"bp_map"`
	Board             GameBoard       `json:"board"`
}

type Player struct {
	CID      string `json:"cid"`
	Nickname string `json:"nickname"`
	Avatar   int    `json:"avatar"`
}

type GameBoard struct {
	P1Ban  []int `json:"p_1_ban"`
	P2Ban  []int `json:"p_2_ban"`
	P1Pick []int `json:"p_1_pick"`
	P2Pick []int `json:"p_2_pick"`
}

func GetGameState(gid string, r *redis.Client) (*GameState, error) {
	rawData, err := r.Get(context.Background(), gid).Result()

	if err != nil && !errors.Is(err, redis.Nil) {
		log.Println(err)
		return nil, errors.New(global.TextConfig["redis_data_error"])
	}

	if len(rawData) == 0 {
		return nil, nil
	}

	gameState := new(GameState)

	err = helper.StringToStruct(rawData, gameState)
	if err != nil {
		return nil, errors.New(global.TextConfig["redis_data_error"])
	}

	return gameState, nil
}

type MoveRequest struct {
	Call string      `json:"call"`
	Data interface{} `json:"data"`
}

type ChatInfo struct {
	Message  string `json:"message"`
	CID      string `json:"cid"`
	Nickname string `json:"nickname"`
	JoinChat bool   `json:"join_chat"`
}
