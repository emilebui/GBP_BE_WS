package test

import (
	"github.com/emilebui/GBP_BE_WS/internal/logic"
	"github.com/emilebui/GBP_BE_WS/pkg/gstatus"
	"testing"
	"time"
)

func TestGameConnection(t *testing.T) {

	// Init Test Context
	wsHandler, r := initTestContext(t)

	resetFakeRedis()
	time.Sleep(500 * time.Millisecond)

	println()
	printTestLog("____________TEST GAME CONNECTION______________")
	println()
	printTestLog("Init Test Context Successfully...")

	// Open test connection
	p1Conn := generatePlayerConn(t, wsHandler, P1CID)

	printTestLog("Created Test Websocket for Player 1 Successfully...")

	time.Sleep(1 * time.Second)
	gs, err := logic.GetGameState(GIDTest, r)
	if err != nil {
		t.Fatalf("Redis error: %v", err)
	}

	if gs.Status != gstatus.WATTING {
		t.Fatalf("The game should be in waiting mode!!")
	}

	time.Sleep(1 * time.Second)
	makeMove(t, p1Conn, "BAN", 69)
	readMessage(t, p1Conn)

	p2Conn := generatePlayerConn(t, wsHandler, P2CID)

	printTestLog("Created Test Websocket for Player 2 Successfully...")

	time.Sleep(1 * time.Second)
	gs, err = logic.GetGameState(GIDTest, r)
	if err != nil {
		t.Fatalf("Redis error: %v", err)
	}

	err = p2Conn.Close()
	if err != nil {
		printTestLog("Failed to close p2 connection")
	}

	time.Sleep(1 * time.Second)
	gs, err = logic.GetGameState(GIDTest, r)
	if err != nil {
		t.Fatalf("Redis error: %v", err)
	}

	if gs.Status != gstatus.HALT {
		t.Fatalf("This game should be halted!!!")
	}

	// Unable to submit move
	makeMove(t, p1Conn, "BAN", 69)

	p2Conn = generatePlayerConn(t, wsHandler, P2CID)
	time.Sleep(1 * time.Second)

	gs, err = logic.GetGameState(GIDTest, r)
	if err != nil {
		t.Fatalf("Redis error: %v", err)
	}

	if gs.Status != gstatus.PLAYING {
		t.Fatalf("This game should be playing since p2 reconnected!!!")
	}

	println()
	printTestLog("____________END TEST______________")
	println()
	p1Conn.Close()
	p2Conn.Close()
}
