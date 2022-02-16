package main

import (
	"database/sql"
	"encoding/json"
	"log"
)

func getUserByID(db *sql.DB, id string) string {
	type userResponse struct {
		Id       int    `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	query := `SELECT id, username, email FROM users where id=$1`
	rows, err := db.Query(query, id)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	var response string
	for rows.Next() {
		var id int
		var username string
		var email string
		if err := rows.Scan(&id, &username, &email); err != nil {
			log.Fatal(err)
		}

		res := &userResponse{
			Id:       id,
			Username: username,
			Email:    email,
		}

		reply, _ := json.Marshal(res)
		response = string(reply)
		log.Println(res.Username)
		log.Println(res.Email)
	}
	return response
}

func getPPByID(db *sql.DB, id string) []byte {
	query := `SELECT profile_picture FROM users where id=$1`
	rows, err := db.Query(query, id)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	var response []byte
	for rows.Next() {
		var profile_picture []byte
		if err := rows.Scan(&profile_picture); err != nil {
			log.Fatal(err)
		}

		response = profile_picture
	}
	return response
}
