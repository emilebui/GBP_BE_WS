package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/emilebui/GBP_BE_echo/pkg/helper"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
)

type WebSocketHandler struct {
	redisConn   *redis.Client
	upgrader    websocket.Upgrader
	responseMap map[string]string
}

func NewWSHandler(redisConn *redis.Client, upgrader websocket.Upgrader, responseMap map[string]string) *WebSocketHandler {
	return &WebSocketHandler{
		redisConn:   redisConn,
		upgrader:    upgrader,
		responseMap: responseMap,
	}
}

func (s *WebSocketHandler) Play(w http.ResponseWriter, r *http.Request) {

	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade Error:", err)
		return
	}

	gid, cid, err := s.getParams(r)
	if err != nil {
		helper.WSError(c, err, "Param Error")
		return
	}

	err = s.ConnectGame(gid, cid)
	if err != nil {
		helper.WSError(c, err, "Connect Game Error")
		return
	}

	log.Printf("Client %s connected successfully to the game %s\n", cid, gid)
	defer c.Close()

	go s.handleRedisMessage(c, gid)
	s.handleWSMessage(c)
}

func (s *WebSocketHandler) getParams(r *http.Request) (gid string, cid string, err error) {
	params := r.URL.Query()
	gameID := params.Get("gid")
	if gameID == "" {
		return "", "", errors.New(fmt.Sprintf(s.responseMap["params_required"], "gid"))
	}
	clientID := params.Get("cid")
	if clientID == "" {
		return "", "", errors.New(fmt.Sprintf(s.responseMap["params_required"], "cid"))
	}
	return gameID, clientID, nil
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
			log.Println("WriteMessage Error:", err)
			break
		}
	}
}

func (s *WebSocketHandler) handleWSMessage(c *websocket.Conn) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("ReadMessage Error:", err)
			break
		}
		fmt.Printf("Got message: %s\n", string(message))

		err = s.redisConn.Publish(context.Background(), "livechat", message).Err()
		if err != nil {
			log.Println("Publish Redis Error:", err)
			break
		}
	}
}
