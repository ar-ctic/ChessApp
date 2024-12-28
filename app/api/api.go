package api

import (
	"ChessApp/service/user"
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)
	


type APIServer struct {
	addr string
	db *sql.DB 
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer {
		addr: addr,
		db: db,
	}
}

func (s* APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userApp := user.NewApp(s.db)
	userHandler := user.NewHandler(userApp)
	userHandler.RegisterRoutes(subrouter)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
} 