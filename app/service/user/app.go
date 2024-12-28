package user

import (
	"ChessApp/types"
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type App struct {
	db *sql.DB
}

func NewApp(db *sql.DB) *App {
	return &App{db: db}
}

func (a *App) GetUserByEmail(email string) (*types.User, error) {

	rows, err := a.db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}

	u := new(types.User)
	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == "" {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil

}

func (a *App) GetUserByUsername(username string) (*types.User, error) {

	rows, err := a.db.Query("SELECT * FROM users WHERE username = ?", username)
	if err != nil {
		return nil, err
	}

	u := new(types.User)

	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == "" {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil

}

func scanRowIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *App) GetUserByID(id string) (*types.User, error) {

	rows, err := a.db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	u := new(types.User)

	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == "" {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil

}

func (a *App) CreateUser(user types.User) error {

	id := uuid.New()

	_, err := a.db.Exec("INSERT INTO users (id, username, email, password) VALUES (?, ?, ?, ?)", id, user.Username, user.Email, user.Password)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
