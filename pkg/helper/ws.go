package helper

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

func WSError(c *websocket.Conn, err error, v ...string) {
	errMsg := fmt.Sprintf("%v %v", v, err)

	log.Println(errMsg)
	_ = c.WriteMessage(1, []byte(errMsg))
}
