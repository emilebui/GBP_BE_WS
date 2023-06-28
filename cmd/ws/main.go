package main

import (
	"fmt"
	"github.com/emilebui/GBP_BE_WS/internal/handler"
	"github.com/emilebui/GBP_BE_WS/internal/logic"
	"github.com/emilebui/GBP_BE_WS/pkg/conf"
	"github.com/emilebui/GBP_BE_WS/pkg/conn"
	"github.com/emilebui/GBP_BE_WS/pkg/global"
	"github.com/emilebui/GBP_BE_WS/pkg/helper"
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
		check := helper.ContainsString(CORS, r.URL.Host)
		if !check {
			log.Printf("Error Origin: Host - (%s) | URL - (%s)", r.URL.Host, r.URL)
		}
		return check
	}

	// Init text message config
	textConf := config.GetStringMapString("text_messages")
	global.InitGlobalTextConfig(textConf)

	// Init After Game Expiration
	global.InitAfterGameExp(config.GetInt("AFTER_GAME_EXP"))

	// Get Turn Format
	var tf map[string]map[int]logic.TurnInfo
	err := config.UnmarshalKey("GAME_TURN_FORMAT", &tf)
	if err != nil {
		log.Fatal(err)
	}
	logic.InitTurnFormat(tf)

	// Init Redis Connection
	redisConn := conn.GetRedisConn(config)

	// Init ws handler
	wsHandler := handler.NewWSHandler(redisConn, upgrader)

	fmt.Println("Starting Websocket server successfully!!!")
	http.HandleFunc("/play", wsHandler.Play)
	http.HandleFunc("/watch", wsHandler.Watch)
	log.Fatal(http.ListenAndServe(config.GetString("addr"), nil))
}
