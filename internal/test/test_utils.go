package test

import (
	"fmt"
	"github.com/emilebui/GBP_BE_WS/internal/handler"
	"github.com/emilebui/GBP_BE_WS/internal/logic"
	"github.com/emilebui/GBP_BE_WS/pkg/conf"
	"github.com/emilebui/GBP_BE_WS/pkg/global"
	"github.com/emilebui/GBP_BE_WS/pkg/helper"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

const (
	P1CID   = "player1"
	P2CID   = "player2"
	GIDTest = "blah"
)

func initTestContext(t *testing.T) (*handler.WebSocketHandler, *redis.Client) {

	config := conf.Get("../../config.yaml")
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// Init text message config
	textConf := config.GetStringMapString("text_messages")
	global.InitGlobalTextConfig(textConf)

	// Init After Game Expiration
	global.InitAfterGameExp(config.GetInt("AFTER_GAME_EXP"))

	var tf map[string]map[int]logic.TurnInfo
	err := config.UnmarshalKey("GAME_TURN_FORMAT", &tf)
	if err != nil {
		t.Error("Failed to load config")
	}
	logic.InitTurnFormat(tf)

	r := getFakeRedis(t)

	wsHandler := handler.NewWSHandler(r, upgrader)

	return wsHandler, r
}

func printTestLog(msg string) {
	println()
	println(msg)
	println()
}

func generatePlayerConn(t *testing.T, wsHandler *handler.WebSocketHandler, cid string) *websocket.Conn {

	srv := httptest.NewServer(http.HandlerFunc(wsHandler.Play))
	u, _ := url.Parse(srv.URL)
	u.Scheme = "ws"
	u.RawQuery = fmt.Sprintf("gid=%s&cid=%s", GIDTest, cid)

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("cannot make websocket connection: %v", err)
	}

	return conn
}

func readMessage(t *testing.T, wsConn *websocket.Conn) *WSResponseMsg {
	_, p, err := wsConn.ReadMessage()
	if err != nil {
		t.Fatalf("cannot read message: %v", err)
	}

	result := &WSResponseMsg{}
	err = helper.BytesToStruct(p, result)
	if err != nil {
		t.Fatalf("cannot parse ws message: %v", err)
	}

	return result
}

type WSResponseMsg struct {
	Message string `json:"message"`
	Type    int    `json:"type"`
	Data    string `json:"data"`
}

func createMoveSet(t *testing.T, p1Conn *websocket.Conn, p2Conn *websocket.Conn, gs *logic.GameState) *MoveSet {
	if gs.Player1.CID == P1CID {
		return newMoveSet(t, p1Conn, p2Conn)
	} else {
		return newMoveSet(t, p2Conn, p1Conn)
	}
}

func makeMove(t *testing.T, wsConn *websocket.Conn, call string, data interface{}) {
	mr := &logic.MoveRequest{
		Call: call,
		Data: data,
	}

	err := wsConn.WriteMessage(1, []byte(helper.Struct2String(mr)))
	if err != nil {
		t.Fatalf("failed to write message: %v", err)
	}
	time.Sleep(250 * time.Millisecond)
}
