# Genshin Ban Pick - Websocket Server (GBP_WS)

- This is the websocket backend server for Genshin Ban Pick
- Tech stack:
  - gorilla/websocket
  - redis
  - spf13/viper

## Installation Instructions

### Requirements
- git
- golang 1.19

### Initial setup
- Install dependencies:
    - Run: `go mod download`

### Run application

- **RUN**:
    - `go run cmd/ws/main.go`

- **Note**:
    - You may need to run redis database and config them first to run the app
    - Check `docker-compose.yaml` and `config.yaml`

### Docker

- Build: `docker build -t gbp_ws:latest .`
- To build and run the entire stack (gbp_ws and redis):
    - Run: `docker-compose up -d`

### Testing

- Run: `go test ./internal/test/... -cover -coverpkg ./...`

### Project Structure

- `cmd`: This contains all the run commands for the app
    -  `ws/main.go`: run the websocket server
- `pkg`: Contains all the pkg, helper modules to build the app
    - `conf`: Contains helper function for app configurations
    - `conn`: Contains helper functions to connect to redis
    - `global`: Contains singleton global configuration
    - `gstatus`: Contains enums, global structs
    - `helper`: Contains general helper functions
- `internal`: Contain endpoint and logic for each API. The structure of this module is as follows:
    - `handler`: Handle the all the endpoints of the application (take request, get input and pass it to logic layer)
    - `logic`: Logic Layer, all the logic of the service are handled here
    - `broker`: All the websocket messages are translated to input for logic functions