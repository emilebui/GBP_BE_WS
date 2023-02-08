package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/emilebui/GBP_BE_echo/internal/broker"
	"github.com/emilebui/GBP_BE_echo/pkg/global"
	"github.com/emilebui/GBP_BE_echo/pkg/gstatus"
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

	// Create done channel
	done := make(chan bool)

	go s.handleRedisMessage(done, c, gid)
	s.handleWSMessage(done, c, gid, cid)

}

func (s *WebSocketHandler) getParams(r *http.Request) (gid string, cid string, err error) {
	params := r.URL.Query()
	gameID := params.Get("gid")
	if gameID == "" {
		return "", "", errors.New(fmt.Sprintf(global.TextConfig["params_required"], "gid"))
	}
	clientID := params.Get("cid")
	if clientID == "" {
		return "", "", errors.New(fmt.Sprintf(global.TextConfig["params_required"], "cid"))
	}
	return gameID, clientID, nil
}

func (s *WebSocketHandler) handleRedisMessage(done <-chan bool, c *websocket.Conn, gid string) {
	subscriber := s.redisConn.Subscribe(context.Background(), gid)

	for {
		select {
		case <-done:
			println("Unsub and closed go routines!!")
			_ = subscriber.Unsubscribe(context.Background(), gid)
			return
		default:
			msg, err := subscriber.ReceiveMessage(context.Background())
			if err != nil {
				panic(err)
			}

			err = c.WriteMessage(1, []byte(msg.Payload))
			if err != nil {
				log.Println("WriteMessage Error:", err)
			}
		}
	}
}

func (s *WebSocketHandler) handleWSMessage(done chan<- bool, c *websocket.Conn, gid string, cid string) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("ReadMessage Error:", err)
			helper.PublishRedis(&gstatus.ResponseMessage{
				Message: fmt.Sprintf("Player %s has disconnected", cid),
				Type:    gstatus.LOG,
			}, s.redisConn, gid)
			done <- true
			break
		}
		log.Printf("Got message: %s - GID: %s - CID: %s\n", string(message), gid, cid)

		err = broker.ProcessMessage(s.redisConn, message, gid, cid)
		if err != nil {
			helper.WSError(c, err, fmt.Sprintf("Handle Message Error - GID: %s - CID: %s", gid, cid))
		}
	}
}
