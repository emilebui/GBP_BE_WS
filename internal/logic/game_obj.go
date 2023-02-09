package logic

import (
	"context"
	"errors"
	"github.com/emilebui/GBP_BE_echo/pkg/global"
	"github.com/emilebui/GBP_BE_echo/pkg/helper"
	"github.com/redis/go-redis/v9"
	"log"
	"math/rand"
	"time"
)

type GameState struct {
	GameID            string          `json:"game_id"`
	Player1           Player          `json:"player_1"`
	Player2           Player          `json:"player_2"`
	Turn              int             `json:"turn"`
	Status            int             `json:"status"`
	PlayerTurnMap     map[int]string  `json:"player_turn_map"`
	PlayerTurn        string          `json:"player_turn"`
	ConnectionTracker map[string]bool `json:"connection_tracker"`
	Pick              bool            `json:"pick"`
	BPMap             map[int]bool    `json:"bp_map"`
	Board             GameBoard       `json:"board"`
}

type Player struct {
	CID      string `json:"cid"`
	Nickname string `json:"nickname"`
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

func ShufflePlayer(p1 string, p2 string) map[int]string {
	rand.Seed(time.Now().Unix())
	check := rand.Intn(101)
	if check > 50 {
		return map[int]string{1: p2, 2: p1}
	}

	return map[int]string{1: p1, 2: p2}
}

type MoveRequest struct {
	Call string      `json:"call"`
	Data interface{} `json:"data"`
}
