package utils

import (
	"ChessApp/types"
	"encoding/json"
	"fmt"
	"net/http"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
)

var Validate = validator.New()

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})

}

func SendJSON(conn *websocket.Conn, success bool, status int, message string) {
	errorMessage := types.ErrorMessage{
		Success:   success,
		Message: message,
		Status:  status,
	}
	jsonData, err := json.Marshal(errorMessage)
	if err != nil {
		log.Println("Error marshalling error message:", err)
		return
	}
	fmt.Printf("Data to send: %s\n", string(jsonData))
	err = conn.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		log.Println("Error sending WebSocket message:", err)
	}
}

