package logic

import (
	"math/rand"
	"time"
)

type TurnInfo struct {
	Player int  `mapstructure:"player"`
	Pick   bool `mapstructure:"pick"`
}

var TurnFormat map[int]TurnInfo

func InitTurnFormat(tf map[int]TurnInfo) {
	TurnFormat = tf
}

func GetPlayerTurn(gs *GameState) string {
	turnInfo := TurnFormat[gs.Turn]

	if turnInfo.Player != 1 {
		return gs.Player2.CID
	}

	return gs.Player1.CID
}

func CheckIfPickTurn(turn int) bool {
	return TurnFormat[turn].Pick
}

func ShufflePlayer(gs *GameState) {
	rand.Seed(time.Now().Unix())
	check := rand.Intn(101)
	if check > 50 {
		SwapPlayer(gs)
	}
}

func SwapPlayer(gs *GameState) {
	temp := gs.Player1
	gs.Player1 = gs.Player2
	gs.Player2 = temp
}
