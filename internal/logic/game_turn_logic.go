package logic

import (
	"math/rand"
	"time"
)

type TurnInfo struct {
	Player int  `mapstructure:"player"`
	Pick   bool `mapstructure:"pick"`
}

type FormatInfo struct {
	NumAbyssTurn int `mapstructure:"abyss"`
	NumSOTurn    int `mapstructure:"onslaught"`
}

var TurnFormat map[string]map[int]TurnInfo
var FormatInfoMap map[string]FormatInfo

func InitTurnFormat(tf map[string]map[int]TurnInfo, fi map[string]FormatInfo) {
	TurnFormat = tf
	FormatInfoMap = fi
}

func GetPlayerTurn(gs *GameState) string {

	turnFormat := gs.Settings.NumBan

	turnInfo := TurnFormat[turnFormat][gs.Turn]

	if turnInfo.Player != 1 {
		return gs.Player2.CID
	}

	return gs.Player1.CID
}

func CheckIfPickTurn(gs *GameState) bool {
	turnFormat := gs.Settings.NumBan

	return TurnFormat[turnFormat][gs.Turn].Pick
}

func CheckIfTurnValid(gs *GameState) bool {
	if gs.Turn > GetMaxTurn(gs) {
		return false
	}

	turnFormat := gs.Settings.NumBan

	_, ok := TurnFormat[turnFormat][gs.Turn]
	return ok
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

func GetMaxTurn(gs *GameState) int {
	if gs.Settings.SOMode {
		return FormatInfoMap[gs.Settings.NumBan].NumSOTurn
	}

	return FormatInfoMap[gs.Settings.NumBan].NumAbyssTurn
}
