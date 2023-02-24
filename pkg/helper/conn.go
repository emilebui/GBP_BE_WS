package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/emilebui/GBP_BE_WS/pkg/gstatus"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"log"
)

func ResponseWS(res *gstatus.ResponseMessage, c *websocket.Conn) {
	bytes, _ := json.Marshal(res)
	_ = c.WriteMessage(1, bytes)
}

func PublishRedis(res *gstatus.ResponseMessage, r *redis.Client, gid string) {
	dataStr := Struct2String(res)
	r.Publish(context.Background(), gid, dataStr)
}

func WSError(c *websocket.Conn, err error, et int, v ...string) {
	errMsg := fmt.Sprintf("%v %v", v, err)
	log.Println(errMsg)
	ResponseWS(&gstatus.ResponseMessage{
		Message: err.Error(),
		Type:    et,
		Info:    v,
	}, c)
}
