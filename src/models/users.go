package models

import (
	"database/sql"
	"diorama/v2/auth"
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

func GetUserById(db *sql.DB, id string) *userResponse {
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

func SetUserDetail(db *sql.DB, userID string, username string, name string, email string) string {
	log.Printf("update user with ID= " + userID)
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

func SetUserPassword(db *sql.DB, userID int, oldPassword string, newPassword string) string {
	query := `SELECT password FROM users where id=$1`
	var currPassword string
	err := db.QueryRow(query, userID).Scan(&currPassword)
	if err != nil {
		log.Println(err)
		return err.Error()
	}

	if !auth.CheckPasswordHash(oldPassword, currPassword) {
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

func GetPPByID(db *sql.DB, id string) []byte {
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

func SetUserPP(db *sql.DB, picture []byte, id string) string {
	log.Println("Update picture at user ID=" + id)
	insertDynStmt := `update users set profile_picture=$1 where id=$2`
	_, err := db.Exec(insertDynStmt, picture, id)

	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return "true"
}

func DeleteUser(db *sql.DB, userID string) string {
	log.Printf("delete user with ID= " + userID)
	query := `delete from users where id=$1`
	_, err := db.Exec(query, userID)

	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return "true"
}

type searchResponse struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

func SearchUser(db *sql.DB, query string) []*searchResponse {
	log.Println("Search user")
	sqlquery := `SELECT id, username FROM users where username LIKE $1`
	rows, err := db.Query(sqlquery, query+"%")
	var response []*searchResponse
	if err != nil {
		log.Println(err)
		return response
	}

	defer rows.Close()
	for rows.Next() {
		var id int
		var username string
		if err := rows.Scan(&id, &username); err != nil {
			log.Println(err)
			return response
		}

		res := &searchResponse{
			Id:       id,
			Username: username,
		}
		response = append(response, res)

	}
	return response
}
