package auth

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

func CreateJWTToken(user_id int) (string, int64, error) {
	//belum disimpen di env
	exp := time.Now().Add(time.Minute * 30).Unix()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user_id
	claims["exp"] = exp

	err := godotenv.Load("../.env")
	if err != nil {
		return "", 0, err
	}

	secret_key := os.Getenv("SECRET_KEY")
	log.Println(secret_key)

	t, err := token.SignedString([]byte(secret_key))
	if err != nil {
		return "", 0, err
	}

	return t, exp, nil
}
