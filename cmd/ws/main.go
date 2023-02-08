package main

import (
	"fmt"
	"github.com/emilebui/GBP_BE_echo/internal/handler"
	"github.com/emilebui/GBP_BE_echo/internal/logic"
	"github.com/emilebui/GBP_BE_echo/pkg/conf"
	"github.com/emilebui/GBP_BE_echo/pkg/conn"
	"github.com/emilebui/GBP_BE_echo/pkg/global"
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

	// Init text message config
	textConf := config.GetStringMapString("text_messages")
	global.InitGlobalTextConfig(textConf)

	// Get Turn Format
	var tf map[int]logic.TurnInfo
	err := config.UnmarshalKey("GAME_TURN_FORMAT", &tf)
	if err != nil {
		log.Fatal(err)
	}
	logic.InitTurnFormat(tf)

	// Init Redis Connection
	redisConn := conn.GetRedisConn(config)

	// Init game logic

	// Init ws handler
	wsHandler := handler.NewWSHandler(redisConn, upgrader)

	fmt.Println("Starting Websocket server successfully!!!")
	http.HandleFunc("/play", wsHandler.Play)
	log.Fatal(http.ListenAndServe(config.GetString("addr"), nil))
}
