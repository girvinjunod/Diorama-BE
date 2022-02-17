package main

import (
	"database/sql"
	"log"
)

type userResponse struct {
	Error    string `json:"error"`
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func getUserByID(db *sql.DB, id string) *userResponse {
	query := `SELECT id, username, email FROM users where id=$1`
	rows, err := db.Query(query, id)
	var res *userResponse
	if err != nil {
		return res
	}

	defer rows.Close()
	for rows.Next() {
		var id int
		var username string
		var email string
		if err := rows.Scan(&id, &username, &email); err != nil {
			log.Println(err)
			return res
		}

		res = &userResponse{
			Error:    "false",
			Id:       id,
			Username: username,
			Email:    email,
		}
		log.Println(res.Username)
		log.Println(res.Email)
	}
	return res
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
