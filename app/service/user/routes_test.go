package user

import (
	"ChessApp/service/auth"
	"ChessApp/types"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestUserRegister(t *testing.T) {
	userApp := &mockUserAppRegister{}
	handler := NewHandler(userApp)

	register_payloads := []struct {
		Name     string
		Payload  types.RegisterUserPayload
		Expected int
	}{
		{
			Name: "Valid Register Payload",
			Payload: types.RegisterUserPayload{
				Username: "testuser",
				Email:    "testuser@example.com",
				Password: "strongpassword",
			},
			Expected: http.StatusCreated,
		},
		{
			Name: "Valid Shortest Name Register Payload",
			Payload: types.RegisterUserPayload{
				Username: "test",
				Email:    "testuser@example.com",
				Password: "strongpassword",
			},
			Expected: http.StatusCreated,
		},
		{
			Name: "Valid Longest Name Register Payload",
			Payload: types.RegisterUserPayload{
				Username: strings.Repeat("a", 20),
				Email:    "testuser@example.com",
				Password: "strongpassword",
			},
			Expected: http.StatusCreated,
		},
		{
			Name: "Valid Shortest Password Register Payload",
			Payload: types.RegisterUserPayload{
				Username: "testuser",
				Email:    "testuser@example.com",
				Password: "strongpa",
			},
			Expected: http.StatusCreated,
		},
		{
			Name: "Valid Longest Password Register Payload",
			Payload: types.RegisterUserPayload{
				Username: "testuser",
				Email:    "testuser@example.com",
				Password: strings.Repeat("a", 64),
			},
			Expected: http.StatusCreated,
		},
		{
			Name: "Empty Username Register Payload",
			Payload: types.RegisterUserPayload{
				Username: "",
				Email:    "testuser@example.com",
				Password: "strongpassword",
			},
			Expected: http.StatusBadRequest,
		},
		{
			Name: "Missing Username Register Payload",
			Payload: types.RegisterUserPayload{
				Email:    "testuser@example.com",
				Password: "strongpassword",
			},
			Expected: http.StatusBadRequest,
		},
		{
			Name: "Too short Username Register Payload",
			Payload: types.RegisterUserPayload{
				Username: "tes",
				Email:    "testuser@example.com",
				Password: "strongpassword",
			},
			Expected: http.StatusBadRequest,
		},
		{
			Name: "Too long Username Register Payload",
			Payload: types.RegisterUserPayload{
				Username: strings.Repeat("a", 33),
				Email:    "testuser@example.com",
				Password: "strongpassword",
			},
			Expected: http.StatusBadRequest,
		},
		{
			Name: "Empty Email Register Payload",
			Payload: types.RegisterUserPayload{
				Username: "testuser",
				Email:    "",
				Password: "strongpassword",
			},
			Expected: http.StatusBadRequest,
		},
		{
			Name: "Missing Email Register Payload",
			Payload: types.RegisterUserPayload{
				Username: "testuser",
				Password: "strongpassword",
			},
			Expected: http.StatusBadRequest,
		},
		{
			Name: "Invalid Email Register Payload",
			Payload: types.RegisterUserPayload{
				Username: "testuser",
				Email:    "testuserexample.com",
				Password: "strongpassword",
			},
			Expected: http.StatusBadRequest,
		},
		{
			Name: "Empty Password Register Payload",
			Payload: types.RegisterUserPayload{
				Username: "testuser",
				Email:    "testuser@example.com",
				Password: "",
			},
			Expected: http.StatusBadRequest,
		},
		{
			Name: "Missing Password Register Payload",
			Payload: types.RegisterUserPayload{
				Username: "testuser",
				Email:    "testuser@example.com",
			},
			Expected: http.StatusBadRequest,
		},
		{
			Name: "Too short Password Register Payload",
			Payload: types.RegisterUserPayload{
				Username: "testuser",
				Email:    "testuser@example.com",
				Password: "str",
			},
			Expected: http.StatusBadRequest,
		},
		{
			Name: "Too long Password Register Payload",
			Payload: types.RegisterUserPayload{
				Username: "testuser",
				Email:    "testuser@example.com",
				Password: strings.Repeat("a", 65),
			},
			Expected: http.StatusBadRequest,
		},
	}


	for _, tc := range register_payloads {
		t.Run(tc.Name, func(t *testing.T) {

	
			marshalled, _ := json.Marshal(tc.Payload)
	
			req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
			if err != nil {
				t.Fatal(err)
			}
	
			rr := httptest.NewRecorder()
			router := mux.NewRouter()
	
			router.HandleFunc("/register", handler.handleRegister)
			router.ServeHTTP(rr, req)
	
			if rr.Code != tc.Expected {
				t.Errorf("expected status code %d, got %d", tc.Expected, rr.Code)
			}
		})
	}
	
}

func TestUserLogin(t *testing.T) {
	userApp := &mockUserAppLogin{}
	handler := NewHandler(userApp)

	login_payloads := []struct {
		Name     string
		Payload  types.LoginUserPayload
		Expected int
	}{
		{
			Name: "Valid Login Payload",
			Payload: types.LoginUserPayload{
				Username: "testuser",
				Password: "strongpassword",
			},
			Expected: http.StatusOK,
		},
		{
			Name: "Empty Username Login Payload",
			Payload: types.LoginUserPayload{
				Username: "",
				Password: "strongpassword",
			},
			Expected: http.StatusBadRequest,
		},
		{
			Name: "Missing Username Login Payload",
			Payload: types.LoginUserPayload{
				Username: "",
				Password: "strongpassword",
			},
			Expected: http.StatusBadRequest,
		},
		{
			Name: "Too Short Username Login Payload",
			Payload: types.LoginUserPayload{
				Username: "tes",
				Password: "strongpassword",
			},
			Expected: http.StatusBadRequest,
		},
		{
			Name: "Too Long Login Payload",
			Payload: types.LoginUserPayload{
				Username: strings.Repeat("a", 30),
				Password: "strongpassword",
			},
			Expected: http.StatusBadRequest,
		},
		{
			Name: "Empty Password Login Payload",
			Payload: types.LoginUserPayload{
				Username: "testuser",
				Password: "",
			},
			Expected: http.StatusBadRequest,
		},
		{
			Name: "Missing Password Login Payload",
			Payload: types.LoginUserPayload{
				Username: "testuser",
			},
			Expected: http.StatusBadRequest,
		},
		{
			Name: "Too Short Password Login Payload",
			Payload: types.LoginUserPayload{
				Username: "testuser",
				Password: "strong",
			},
			Expected: http.StatusBadRequest,
		},
		{
			Name: "Too Long Password Payload",
			Payload: types.LoginUserPayload{
				Username: "testuser",
				Password: strings.Repeat("a", 70),
			},
			Expected: http.StatusBadRequest,
		},
	}


	for _, tc := range login_payloads {
		t.Run(tc.Name, func(t *testing.T) {

	
			marshalled, _ := json.Marshal(tc.Payload)
	
			req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(marshalled))
			if err != nil {
				t.Fatal(err)
			}
	
			rr := httptest.NewRecorder()
			router := mux.NewRouter()
	
			router.HandleFunc("/login", handler.handleLogin)
			router.ServeHTTP(rr, req)
	
			if rr.Code != tc.Expected {
				t.Errorf("expected status code %d, got %d", tc.Expected, rr.Code)
			}
		})
	}
	
}

type mockUserAppRegister struct{}

func (m *mockUserAppRegister) GetUserByEmail(email string) (*types.User, error) {

	return nil, fmt.Errorf("user not found")
}
func (m *mockUserAppRegister) GetUserByUsername(username string) (*types.User, error) {
	return nil, fmt.Errorf("user not found")
}
func (m *mockUserAppRegister) GetUserByID(id string) (*types.User, error) {
	return nil, nil
}
func (m *mockUserAppRegister) CreateUser(user types.User) error {
	return nil
}

type mockUserAppLogin struct{}

func (m *mockUserAppLogin) GetUserByEmail(email string) (*types.User, error) {
	return nil, nil
}

func (m *mockUserAppLogin) GetUserByUsername(username string) (*types.User, error) {

	if len(username) < 4 || len(username) > 20 {
		return nil, fmt.Errorf("error")
	}

	password, err := auth.HashPassword("strongpassword")
	if err != nil {
		return nil, fmt.Errorf("error")
	}
	return &types.User{
		ID: "",
		Username: username,
		Email: "",
		Password: password,
	}, nil
}

func (m *mockUserAppLogin) GetUserByID(id string) (*types.User, error) {
	return nil, nil
}

func (m *mockUserAppLogin) CreateUser(user types.User) error {
	return nil
}