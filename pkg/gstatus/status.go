package gstatus

// Game status
const (
	WATTING int = 0
	PLAYING int = 1
	ENDED   int = 2
)

// Message Types
const (
	VALID int = 0
	ERROR int = 1
	INFO  int = 2
	LOG   int = 3
)

type ResponseMessage struct {
	Message string      `json:"message"`
	Type    int         `json:"type"`
	Data    string      `json:"data"`
	Info    interface{} `json:"info,omitempty"`
}
