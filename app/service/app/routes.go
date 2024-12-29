package app

import (
	"ChessApp/service/auth"
	"ChessApp/types"
	"ChessApp/utils"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handler struct {
	app     types.ChessApp
	userApp types.UserApp
}

func NewHandler(app types.ChessApp, userApp types.UserApp) *Handler {
	return &Handler{
		app:     app,
		userApp: userApp,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/create", h.createGame).Methods(http.MethodPost)
	router.HandleFunc("/game/{id}/join", h.handleJoin).Methods(http.MethodPost)
	router.HandleFunc("/game/{id}", h.handleGame).Methods(http.MethodGet)

}

func (h *Handler) createGame(w http.ResponseWriter, r *http.Request) {

	var payload types.NewGamePayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// Validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	h.app.CreateGame(payload.InitialTime, payload.TimeControl, payload.Color)

}

func (h *Handler) handleJoin(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	gameID := vars["id"]

	game, exists := GameStore[gameID]
	if !exists {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("game does not exist"))
		return
	}

	username := auth.GetUsernameFromJWT(r, h.userApp)

	if username == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user has no JWT"))
		return
	}

	if err := JoinGame(game, username); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

}

func (h *Handler) handleGame(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	gameID := vars["id"]
	username := auth.GetUsernameFromJWT(r, h.userApp)

	game, exists := GameStore[gameID]
	if !exists {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("game does not exist"))
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade to WebSocket", http.StatusInternalServerError)
		return
	}
	defer conn.Close()
	fmt.Printf("WebSocket connection established for game: %s\n", gameID)

	game.mu.Lock()
	game.Connections = append(game.Connections, conn)
	game.mu.Unlock()

	defer func() {
		game.mu.Lock()
		for i, connections := range game.Connections {
			if connections == conn {
				game.Connections = append(game.Connections[:i], game.Connections[i+1:]...)
				fmt.Printf("WebSocket disconnected for game: %s", gameID)
				break
			}
		}
		game.mu.Unlock()

		conn.Close()
	}()

	for {

		// Read Incoming Message
		var message map[string]interface{}
		if err := conn.ReadJSON(&message); err != nil {
			utils.SendJSON(conn, false, http.StatusBadRequest, fmt.Sprintf("error reading message %v", err))
			break
		}

		// Handle message by type
		msgType, ok := message["type"].(string)
		if !ok {
			utils.SendJSON(conn, false, http.StatusBadRequest, "missing or invalid message type")
			continue
		}

		switch msgType {

		case "move":

			move, ok := message["move"].(string)
			if !ok {
				utils.SendJSON(conn, false, http.StatusBadRequest, "no move in request")
				continue
			}

			if err := h.handleMove(game, username, move); err != nil {
				utils.SendJSON(conn, false, http.StatusBadRequest, err.Error())
				continue
			}

		default:
			utils.SendJSON(conn, false, http.StatusBadRequest, "missing type ('move')")
		}
	}
}

func (h *Handler) handleMove(game *ChessGame, username, move string) error {

	if !game.GameStarted {
		return fmt.Errorf("game has not started")
	}

	if !(username == game.PlayerWhite || username == game.PlayerBlack) {
		return fmt.Errorf("user is not part of the game")
	}

	fmt.Println(game.CurrentTurn, username, game.PlayerWhite, game.PlayerBlack)

	if game.CurrentTurn == "w" {
		if !(username == game.PlayerWhite){
			return fmt.Errorf("not blacks turn")
		}
	}

	if game.CurrentTurn == "b" {
		if !(username == game.PlayerBlack){
			return fmt.Errorf("not blacks turn")
		}
	}

	
	newFen, err := MakeMove(game, move)
	if err != nil {
		return fmt.Errorf("error making move %v", err)
	}

	broadcastFen(game, newFen)

	return nil
}

func broadcastFen(game *ChessGame, fen string) {
	game.mu.Lock()
	defer game.mu.Unlock()

	for _, conn := range game.Connections {
		utils.SendJSON(conn, true, http.StatusOK, fen)
	}
}
