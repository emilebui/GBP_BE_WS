package test

import (
	"fmt"
	"github.com/emilebui/GBP_BE_WS/internal/logic"
	"github.com/emilebui/GBP_BE_WS/pkg/gstatus"
	"testing"
	"time"
)

func TestMashUp(t *testing.T) {

	// Init Test Context
	wsHandler, r := initTestContext(t)

	resetFakeRedis()
	time.Sleep(500 * time.Millisecond)

	println()
	printTestLog("____________TEST PLAYING MASHUP GAME______________")
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
	time.Sleep(1 * time.Second)
	conn2 := moveSet.playerMap[2]
	conn1 := moveSet.playerMap[1]

	// Test Move in Wrong Turn
	makeMove(t, conn2, "BAN", 69)

	// Make Wrong Move
	makeMove(t, conn1, "PICK", 69)

	// Wrong Type
	makeMove(t, conn1, "CHAT", 123123)
	makeMove(t, conn1, "BAN", "asd")

	// Correct Move
	makeMove(t, conn1, "BAN", 69)
	makeMove(t, conn2, "BAN", 96)

	// Wrong Move, BAN/PICK Already selected
	makeMove(t, conn1, "BAN", 69)

	i := 3

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

	// Make move when the game is already ended
	makeMove(t, conn1, "PICK", 99)
	makeMove(t, conn2, "PICK", 99)

	_ = p1Conn.Close()
	_ = p2Conn.Close()
	time.Sleep(1 * time.Second)
	p1Conn = generatePlayerConn(t, wsHandler, P1CID)
	p2Conn = generatePlayerConn(t, wsHandler, P2CID)

	println()
	printTestLog("____________END TEST______________")
	println()
	p1Conn.Close()
	p2Conn.Close()
}
