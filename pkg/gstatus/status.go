package gstatus

// Game status
const (
	WATTING int = 0
	PLAYING int = 1
	ENDED   int = 2
	HALT    int = 3
)

// Message Types
const (
	VALID           int = 4
	ERROR           int = 2
	INFO            int = 0
	LOG             int = 1
	RECONN          int = 3
	DISCON          int = 5
	JOIN_GAME_ERROR int = 6
)

type ResponseMessage struct {
	Message string      `json:"message"`
	Type    int         `json:"type"`
	Data    string      `json:"data"`
	Info    interface{} `json:"info,omitempty"`
}
