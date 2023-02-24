package test

import (
	"github.com/emilebui/GBP_BE_WS/internal/logic"
	"github.com/emilebui/GBP_BE_WS/pkg/helper"
	"github.com/gorilla/websocket"
	"testing"
)

type MoveSet struct {
	playerMap map[int]*websocket.Conn
	t         *testing.T
	hid       int
}

func newMoveSet(t *testing.T, p1Conn *websocket.Conn, p2Conn *websocket.Conn) *MoveSet {
	pMap := map[int]*websocket.Conn{
		1: p1Conn,
		2: p2Conn,
	}
	return &MoveSet{
		playerMap: pMap,
		t:         t,
		hid:       1,
	}
}

func (m *MoveSet) MakeTestMove(moveIndex int) {
	turnInfo, ok := logic.TurnFormat[moveIndex]
	if !ok {
		m.t.Fatalf("invalid turn index!!")
	}
	m.makeCall(turnInfo)
}

func (m *MoveSet) makeCall(info logic.TurnInfo) {
	playerConn := m.playerMap[info.Player]

	call := "BAN"
	if info.Pick {
		call = "PICK"
	}

	mr := &logic.MoveRequest{
		Call: call,
		Data: m.hid,
	}

	err := playerConn.WriteMessage(1, []byte(helper.Struct2String(mr)))
	if err != nil {
		m.t.Fatalf("failed to write message: %v", err)
	}
	m.hid = m.hid + 1
}
