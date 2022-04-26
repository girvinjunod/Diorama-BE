package models

import (
	"database/sql"
	"log"
	"time"
)

func AddComment(db *sql.DB, event_id int, user_id int, text string) (string, int) {
	log.Println("Add comment")

	insertDynStmt := `insert into comments (event_id, user_id, text, comment_time) values($1,$2,$3,$4) RETURNING id`
	var id int
	err := db.QueryRow(insertDynStmt, event_id, user_id, text, time.Now()).Scan(&id)

	if err != nil {
		log.Println(err)
		return err.Error(), 0
	}

	log.Printf("Inserted comment with id = %d", id)
	return "true", id
}

type commentResponse struct {
	Error       string `json:"error"`
	Id          int    `json:"id"`
	EventID     int    `json:"eventId"`
	UserID      int    `json:"userId"`
	Text        string `json:"text"`
	CommentTime string `json:"commentTime"`
}

func GetCommentDetailById(db *sql.DB, id string) *commentResponse {
	query := `SELECT id, event_id, user_id, text,comment_time FROM comments where id=$1`
	rows, err := db.Query(query, id)
	var res *commentResponse
	if err != nil {
		return res
	}

	defer rows.Close()
	for rows.Next() {
		var id int
		var eventid int
		var userid int
		var text string
		var comment_time time.Time
		if err := rows.Scan(&id, &eventid, &userid, &text, &comment_time); err != nil {
			log.Println(err)
			return res
		}

		res = &commentResponse{
			Error:       "false",
			Id:          id,
			EventID:     eventid,
			UserID:      userid,
			Text:        text,
			CommentTime: comment_time.Format("2006-01-02 15:04:05"),
		}

	}
	return res
}

type commentResponse2 struct {
	Error       string `json:"error"`
	Id          int    `json:"id"`
	Text        string `json:"text"`
	UserID      int    `json:"userID"`
	Username    string `json:"username"`
	CommentTime string `json:"commentTime"`
}

func GetAllCommentsFromEvent(db *sql.DB, eventID string) (string, []*commentResponse2) {
	log.Println("Get all comments from event")
	query := `SELECT c.id as id, c.text as text, c.user_id as user_id, u.username, c.comment_time as comment_time FROM events e, comments c, users u 
	where e.id=$1 and e.id=c.event_id and c.user_id=u.id`
	rows, err := db.Query(query, eventID)
	var res []*commentResponse2
	if err != nil {
		return err.Error(), res
	}

	defer rows.Close()
	for rows.Next() {
		var id int
		var text string
		var userID int
		var username string
		var comment_time time.Time
		if err := rows.Scan(&id, &text, &userID, &username, &comment_time); err != nil {
			log.Println(err)
			return err.Error(), res
		}

		response := &commentResponse2{
			Error:       "false",
			Id:          id,
			Text:        text,
			UserID:      userID,
			Username:    username,
			CommentTime: comment_time.Format("2006-01-02 15:04:05"),
		}

		res = append(res, response)

	}
	return eventID, res
}

func SetCommentDetail(db *sql.DB, commentID string, text string) string {
	log.Printf("update comment with ID= " + commentID)
	insertDynStmt := `UPDATE comments SET
    text = $1
	WHERE id = $2`
	_, err := db.Exec(insertDynStmt, text, commentID)

	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return "true"
}

func DeleteComment(db *sql.DB, commentID string) string {
	log.Printf("delete comment with ID= " + commentID)
	query := `delete from comments where id=$1`
	_, err := db.Exec(query, commentID)

	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return "true"
}
