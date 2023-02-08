package logic

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
	return gs.PlayerTurnMap[turnInfo.Player]
}

func CheckIfPickTurn(turn int) bool {
	return TurnFormat[turn].Pick
}
