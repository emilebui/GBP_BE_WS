package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/emilebui/GBP_BE_echo/internal/broker"
	"github.com/emilebui/GBP_BE_echo/internal/logic"
	"github.com/emilebui/GBP_BE_echo/pkg/global"
	"github.com/emilebui/GBP_BE_echo/pkg/gstatus"
	"github.com/emilebui/GBP_BE_echo/pkg/helper"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"strconv"
)

type WebSocketHandler struct {
	redisConn *redis.Client
	upgrader  websocket.Upgrader
}

func NewWSHandler(redisConn *redis.Client, upgrader websocket.Upgrader) *WebSocketHandler {
	return &WebSocketHandler{
		redisConn: redisConn,
		upgrader:  upgrader,
	}
}

func (s *WebSocketHandler) Play(w http.ResponseWriter, r *http.Request) {

	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade Error:", err)
		return
	}

	gid, player, err := s.getParams(r)
	if err != nil {
		helper.WSError(c, err, gstatus.JOIN_GAME_ERROR, "Param Error")
		return
	}

	err = s.ConnectGame(gid, player)
	if err != nil {
		helper.WSError(c, err, gstatus.JOIN_GAME_ERROR, "Connect Game Error")
		return
	}

	log.Printf("Client %s connected successfully to the game %s\n", player.CID, gid)
	defer c.Close()

	go s.handleRedisMessage(c, gid)
	s.handleWSMessage(c, gid, player)

}

func (s *WebSocketHandler) getParams(r *http.Request) (gid string, p *logic.Player, err error) {
	params := r.URL.Query()
	gameID := params.Get("gid")
	if gameID == "" {
		return "", nil, errors.New(fmt.Sprintf(global.TextConfig["params_required"], "gid"))
	}
	clientID := params.Get("cid")
	if clientID == "" || clientID == "undefined" {
		return "", nil, errors.New(fmt.Sprintf(global.TextConfig["params_required"], "cid"))
	}
	nn := params.Get("nickname")
	if nn == "" {
		nn = clientID
	}
	avaStr := params.Get("avatar")
	ava := 0
	if avaStr != "" {
		ava, err = strconv.Atoi(avaStr)
		if err != nil {
			ava = 0
		}
	}

	player := &logic.Player{
		CID:      clientID,
		Nickname: nn,
		Avatar:   ava,
	}

	return gameID, player, nil
}

func (s *WebSocketHandler) handleRedisMessage(c *websocket.Conn, gid string) {
	subscriber := s.redisConn.Subscribe(context.Background(), gid)

	for {
		msg, err := subscriber.ReceiveMessage(context.Background())
		if err != nil {
			panic(err)
		}

		err = c.WriteMessage(1, []byte(msg.Payload))
		if err != nil {
			if errors.Is(err, websocket.ErrCloseSent) {
				println("Unsub and closed go routines!!")
				_ = subscriber.Unsubscribe(context.Background(), gid)
			} else {
				log.Println("WriteMessage Error:", err)
			}
			return
		}
	}
}

func (s *WebSocketHandler) handleWSMessage(c *websocket.Conn, gid string, player *logic.Player) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, 1005, 1001) {
				s.disconnect(gid, player)
			} else {
				log.Println("ReadMessage Error:", err)
			}
			break
		}
		log.Printf("Got message: %s - GID: %s - CID: %s\n", string(message), gid, player.CID)

		err = broker.ProcessMessage(s.redisConn, message, gid, player.CID, c)
		if err != nil {
			helper.WSError(c, err, gstatus.ERROR, fmt.Sprintf("Handle Message Error - GID: %s - CID: %s", gid, player.CID))
		}
	}
}
