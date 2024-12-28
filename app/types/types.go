package types

import "time"

type UserApp interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByUsername(username string) (*User, error)
	GetUserByID(id string) (*User, error)
	CreateUser(User) error
}

type ChessApp interface {
	CreateGame(initialTime, timeControl int, color string) error
}

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

type LoginUserPayload struct {
	Username string `json:"username" validate:"required,min=4,max=32"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,min=4,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type NewGamePayload struct {
	GameMode    string `json:"game_mode" validate:"required"`
	Color       string `json:"color" validate:"required"`
	InitialTime int    `json:"initial_time" validate:"required"`
	TimeControl int    `json:"time_control" validate:"required"`
}

type ErrorMessage struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}