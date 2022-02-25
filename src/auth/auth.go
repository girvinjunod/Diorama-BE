package auth

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func Register(db *sql.DB, username string, email string, name string, password string) string {
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

func Login(db *sql.DB, username string, password string) (string, int) {
	log.Println("login")

	query := `SELECT id, password FROM users WHERE username = $1`
	rows, err := db.Query(query, username)

	if err != nil {
		log.Println(err)
		return err.Error(), -999
	}

	count := 0
	var pass_ string
	var user_id int

	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&user_id, &pass_); err != nil {
			log.Println(err)
			return err.Error(), -999
		}
		count++
	}

	if count > 1 {
		log.Println(err)
		return err.Error(), -999
	}
	if CheckPasswordHash(password, pass_) == true {
		return "true", user_id
	}
	log.Println(pass_)
	log.Println(password)
	return "Invalid login credential", -999
}
