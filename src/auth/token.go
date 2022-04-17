package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

func CreateJWTToken(username string) (string, int64, error) {
	//belum disimpen di env
	exp := time.Now().Add(time.Minute * 30).Unix()

	// Create the Claims
	claims := jwt.MapClaims{
		"user_id": username,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	err := godotenv.Load("../.env")
	if err != nil {
		return "", 0, err
	}
	secret_key := os.Getenv("SECRET_KEY")
	// log.Println(secret_key)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(secret_key))

	if err != nil {
		return "", 0, err
	}

	return t, exp, nil
}
