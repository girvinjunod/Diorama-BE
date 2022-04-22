package auth

import (
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

func CreateJWTToken(username string) (string, error) {

	// Create the Claims
	claims := jwt.MapClaims{
		"user_id": username,
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	err := godotenv.Load("../.env")
	if err != nil {
		return "", err
	}
	secret_key := os.Getenv("SECRET_KEY")
	// log.Println(secret_key)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(secret_key))

	if err != nil {
		return "", err
	}

	return t, nil
}
