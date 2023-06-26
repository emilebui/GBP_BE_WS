package test

import (
	"fmt"
	"github.com/emilebui/GBP_BE_WS/internal/logic"
	"github.com/emilebui/GBP_BE_WS/pkg/gstatus"
	"testing"
	"time"
)

func TestPlayEntireGame(t *testing.T) {

	// Init Test Context
	wsHandler, r := initTestContext(t)

	resetFakeRedis()
	time.Sleep(500 * time.Millisecond)

	println()
	printTestLog("____________TEST PLAYING ENTIRE GAME______________")
	println()
	printTestLog("Init Test Context Successfully...")

	// Open test connection
	p1Conn := generatePlayerConn(t, wsHandler, P1CID)

	printTestLog("Created Test Websocket for Player 1 Successfully...")

	p2Conn := generatePlayerConn(t, wsHandler, P2CID)

	printTestLog("Created Test Websocket for Player 2 Successfully...")

	time.Sleep(1 * time.Second)
	gs, err := logic.GetGameState(GIDTest, r)
	if err != nil {
		t.Fatalf("Redis error: %v", err)
	}
	moveSet := createMoveSet(t, p1Conn, p2Conn, gs)
	i := 1

	for i <= len(logic.GetTurnFormat(gs)) {
		printTestLog(fmt.Sprintf("Make Move Turn %d", i))
		moveSet.MakeTestMove(gs, i)
		time.Sleep(250 * time.Millisecond)
		i++
	}

	gs, err = logic.GetGameState(GIDTest, r)
	if err != nil {
		t.Fatalf("Redis error: %v", err)
	}

	if gs.Status != gstatus.ENDED {
		t.Fatalf("This game should be ended!!!")
	}

	println()
	printTestLog("____________END TEST______________")
	println()
	p1Conn.Close()
	p2Conn.Close()
}
