package main

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type userResponse struct {
	Error    string `json:"error"`
	Id       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

func getUserById(db *sql.DB, id string) *userResponse {
	query := `SELECT id, username, name, email FROM users where id=$1`
	rows, err := db.Query(query, id)
	var res *userResponse
	if err != nil {
		return res
	}

	defer rows.Close()
	for rows.Next() {
		var id int
		var username string
		var name string
		var email string
		if err := rows.Scan(&id, &username, &name, &email); err != nil {
			log.Println(err)
			return res
		}

		res = &userResponse{
			Error:    "false",
			Id:       id,
			Username: username,
			Name:     name,
			Email:    email,
		}

	}
	return res
}

func setUserDetail(db *sql.DB, userID int, username string, name string, email string) string {
	log.Printf("update user with ID= %d", userID)
	insertDynStmt := `UPDATE users SET
    username = $1,
    name = $2,
	email = $3
	WHERE id = $4`
	_, err := db.Exec(insertDynStmt, username, name, email, userID)

	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return "true"
}

func setUserPassword(db *sql.DB, userID int, oldPassword string, newPassword string) string {
	query := `SELECT password FROM users where id=$1`
	var currPassword string
	err := db.QueryRow(query, userID).Scan(&currPassword)
	if err != nil {
		log.Println(err)
		return err.Error()
	}

	if !CheckPasswordHash(oldPassword, currPassword) {
		return "Old password doesn't match!"
	}

	if oldPassword == newPassword {
		return "New password is the same as before"
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 14)

	if err != nil {
		log.Println(err)
		return err.Error()
	}

	insertDynStmt := `UPDATE users SET
	password = $1
	WHERE id = $2`
	_, err = db.Exec(insertDynStmt, hash, userID)

	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return "true"
}

func getPPByID(db *sql.DB, id string) []byte {
	query := `SELECT profile_picture FROM users where id=$1`
	var response []byte
	rows, err := db.Query(query, id)
	if err != nil {
		log.Println(err)
		return response
	}

	defer rows.Close()
	for rows.Next() {
		var profile_picture []byte
		if err := rows.Scan(&profile_picture); err != nil {
			log.Println(err)
			return response
		}

		response = profile_picture
	}
	return response
}

func setUserPP(db *sql.DB, picture []byte, id string) string {
	log.Println("Add picture to user ID=" + id)
	insertDynStmt := `update users set profile_picture=$1 where id=$2`
	_, err := db.Exec(insertDynStmt, picture, id)

	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return "true"
}
