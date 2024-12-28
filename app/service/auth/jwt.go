package auth

import (
	"ChessApp/types"
	"ChessApp/config"
	"time"
	"net/http"
	"log"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(secret []byte, userID, username string) (string, error) {

	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    userID,
		"username":  username,
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {

		return "", err
	}

	return tokenString, nil

}

func GetUsernameFromJWT(r *http.Request, app types.UserApp) string{

	tokenString := getTokenFromRequest(r)

	token, err := validateToken(tokenString)
	if err != nil {
		log.Printf("failed to validate token: %v", err)
		return ""
	}

	if !token.Valid {
		log.Println("invalid token")
		return ""
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := claims["userID"].(string)

	user, err := app.GetUserByID(userID)
	if err != nil {
		log.Printf("failed to fetch user: %v", err)
		return ""
	}

	return user.Username

}

func validateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected singing method: %v", token.Header["alg"])
		}
		return []byte(config.Envs.JWTSecret), nil
	})
}

func getTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")

	if tokenAuth == "" {
		return ""
	}

	return tokenAuth

}