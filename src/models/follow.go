package models

import (
	"database/sql"
	"log"
)

func Follow(db *sql.DB, follower_id string, followed_id string) string {

	query := `SELECT *FROM following where follower_id=$1 and followed_id=$2`
	rows, err := db.Query(query, follower_id, followed_id)
	if err != nil {
		return err.Error()
	}

	exist := false
	defer rows.Close()
	for rows.Next() {
		exist = true
	}
	if exist {
		log.Println("Already followed")
		return "true"
	}
	insertDynStmt := `insert into following (follower_id, followed_id) values($1,$2)`
	_, err = db.Exec(insertDynStmt, follower_id, followed_id)

	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return "true"
}

func Unfollow(db *sql.DB, follower_id string, followed_id string) string {

	query := `SELECT *FROM following where follower_id=$1 and followed_id=$2`
	rows, err := db.Query(query, follower_id, followed_id)
	if err != nil {
		log.Println(err)
		return err.Error()
	}
	exist := false
	defer rows.Close()
	for rows.Next() {
		exist = true
	}
	if !exist {
		log.Println("Already unfollowed")
		return "true"
	}
	query = `delete from following where follower_id=$1 and followed_id=$2`
	_, err = db.Exec(query, follower_id, followed_id)

	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return "true"
}

type followResponse struct {
	Error    string `json:"error"`
	UserID   int    `json:"userId"`
	Username string `json:"username"`
}

func GetAllFollowedUsers(db *sql.DB, id string) []*followResponse {
	query := `SELECT f.followed_id as user_id, u.username from users u, following f where f.followed_id=u.id and f.follower_id=$1`
	rows, err := db.Query(query, id)
	var res []*followResponse
	if err != nil {
		log.Println(err)
		return res
	}

	defer rows.Close()
	for rows.Next() {
		var userid int
		var username string
		if err := rows.Scan(&userid, &username); err != nil {
			log.Println(err)
			return res
		}

		response := &followResponse{
			Error:    "false",
			UserID:   userid,
			Username: username,
		}

		res = append(res, response)

	}
	return res
}

func GetAllFollowers(db *sql.DB, id string) []*followResponse {
	query := `SELECT f.follower_id as user_id, u.username from users u, following f where f.follower_id=u.id and f.followed_id=$1`
	rows, err := db.Query(query, id)
	var res []*followResponse
	if err != nil {
		log.Println(err)
		return res
	}

	defer rows.Close()
	for rows.Next() {
		var userid int
		var username string
		if err := rows.Scan(&userid, &username); err != nil {
			log.Println(err)
			return res
		}

		response := &followResponse{
			Error:    "false",
			UserID:   userid,
			Username: username,
		}

		res = append(res, response)

	}
	return res
}

func CheckIfFollowed(db *sql.DB, follower_id string, followed_id string) (string, bool) {
	query := `SELECT *FROM following where follower_id=$1 and followed_id=$2`
	rows, err := db.Query(query, follower_id, followed_id)
	if err != nil {
		log.Println(err)
		return err.Error(), false
	}
	exist := false
	defer rows.Close()
	for rows.Next() {
		exist = true
	}
	return "", exist
}
