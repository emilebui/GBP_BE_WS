package main

import (
	"fmt"
	"github.com/emilebui/GBP_BE_echo/internal/handler"
	"github.com/emilebui/GBP_BE_echo/pkg/conf"
	"github.com/emilebui/GBP_BE_echo/pkg/conn"
	"github.com/emilebui/GBP_BE_echo/pkg/helper"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{}

func main() {
	fmt.Println("Loading Websocket server...")

	// Init config
	config := conf.Get("config.yaml")
	CORS := config.GetStringSlice("CORS")
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return helper.ContainsString(CORS, r.URL.Host)
	}

	// Init Redis Connection
	redisConn := conn.GetRedisConn(config)

	// Init game logic

	// Init ws handler
	wsHandler := handler.NewWSHandler(redisConn, upgrader)

	http.HandleFunc("/play", wsHandler.Play)
	log.Fatal(http.ListenAndServe(config.GetString("addr"), nil))
}
