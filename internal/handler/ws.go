package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/emilebui/GBP_BE_echo/internal/broker"
	"github.com/emilebui/GBP_BE_echo/pkg/global"
	"github.com/emilebui/GBP_BE_echo/pkg/helper"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
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

	gid, cid, nickname, err := s.getParams(r)
	if err != nil {
		helper.WSError(c, err, "Param Error")
		return
	}

	err = s.ConnectGame(gid, cid, nickname)
	if err != nil {
		helper.WSError(c, err, "Connect Game Error")
		return
	}

	log.Printf("Client %s connected successfully to the game %s\n", cid, gid)
	defer c.Close()

	go s.handleRedisMessage(c, gid)
	s.handleWSMessage(c, gid, cid, nickname)

}

func (s *WebSocketHandler) getParams(r *http.Request) (gid string, cid string, nickname string, err error) {
	params := r.URL.Query()
	gameID := params.Get("gid")
	if gameID == "" {
		return "", "", "", errors.New(fmt.Sprintf(global.TextConfig["params_required"], "gid"))
	}
	clientID := params.Get("cid")
	if clientID == "" {
		return "", "", "", errors.New(fmt.Sprintf(global.TextConfig["params_required"], "cid"))
	}
	nn := params.Get("nickname")
	if nn == "" {
		nn = clientID
	}
	return gameID, clientID, nn, nil
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

func (s *WebSocketHandler) handleWSMessage(c *websocket.Conn, gid string, cid string, nickname string) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, 1005) {
				s.disconnect(gid, cid, nickname)
			} else {
				log.Println("ReadMessage Error:", err)
			}
			break
		}
		log.Printf("Got message: %s - GID: %s - CID: %s\n", string(message), gid, cid)

		err = broker.ProcessMessage(s.redisConn, message, gid, cid)
		if err != nil {
			helper.WSError(c, err, fmt.Sprintf("Handle Message Error - GID: %s - CID: %s", gid, cid))
		}
	}
}
