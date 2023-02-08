package logic

import "fmt"

type GameLogic struct {
	State *GameState
}

func (g *GameLogic) Pick(c int) {
	fmt.Printf("Blah %d\n", c)
}
