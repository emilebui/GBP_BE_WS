package logic

import "github.com/emilebui/GBP_BE_echo/pkg/helper"

type GameState struct {
	GameID     string    `json:"game_id"`
	Player1    string    `json:"player_1"`
	Player2    string    `json:"player_2"`
	Turn       int       `json:"turn"`
	Status     int       `json:"status"`
	PlayerTurn string    `json:"player_turn"`
	Pick       bool      `json:"pick"`
	Board      GameBoard `json:"board"`
}

type GameBoard struct {
	P1Ban  []string `json:"p_1_ban"`
	P2Ban  []string `json:"p_2_ban"`
	P1Pick []string `json:"p_1_pick"`
	P2Pick []string `json:"p_2_pick"`
}

func Bytes2GameState(b []byte) (*GameState, error) {
	g := new(GameState)
	return g, helper.BytesToStruct(b, g)
}
