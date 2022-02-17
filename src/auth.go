package main

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func register(db *sql.DB, username string, email string, name string, password string) string {
	log.Println("register")

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Println(err)
		return err.Error()
	}

	insertDynStmt := `insert into users (username, email, name, password) values($1,$2,$3,$4)`
	_, err = db.Exec(insertDynStmt, username, email, name, string(hash))

	if err != nil {
		log.Println(err)
		return err.Error()
	}
	return "true"
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
