package main

import (
	"database/sql"
	"log"
	"time"
)

func addEvent(db *sql.DB, trip_id string, user_id string, caption string, event_date string, post_time string, picture []byte) string {
	log.Println("Add new event")
	event, _ := time.Parse("2006-01-02", event_date)
	postTime, _ := time.Parse("2006-01-02 15:04:05", post_time)
	insertDynStmt := `insert into events (trip_id, user_id, caption, event_date, post_time, picture) values($1,$2,$3,$4,$5,$6)`
	_, err := db.Exec(insertDynStmt, trip_id, user_id, caption, event, postTime, picture)

	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return "true"
}

func setEventPicture(db *sql.DB, picture []byte, eventID string) string {
	log.Println("Add picture to event ID=" + eventID)
	insertDynStmt := `update events set picture=$1 where id=$2`
	_, err := db.Exec(insertDynStmt, picture, eventID)

	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return "true"
}

func getEventPictureByID(db *sql.DB, id string) []byte {
	log.Println("Get event picture with ID " + id)
	query := `SELECT picture FROM events where id=$1`
	rows, err := db.Query(query, id)
	var response []byte
	if err != nil {
		log.Println(err)
		return response
	}

	defer rows.Close()
	for rows.Next() {
		var picture []byte
		if err := rows.Scan(&picture); err != nil {
			log.Println(err)
			return response
		}

		response = picture
	}
	return response
}

type timelineResponse struct {
	UserFeedID int    `json:"userID"`
	Username   string `json:"username"`
	EventID    int    `json:"eventID"`
}

func getTimeline(db *sql.DB, userID string, limit int) []*timelineResponse {
	log.Println("Get Timeline")
	query := `SELECT u.username as Username, f.followed_id as UserFeedID, e.id as EventID FROM users u, following f, events e where f.follower_id = $1 and u.id=f.followed_id and f.followed_id=e.user_id limit $2`
	rows, err := db.Query(query, userID, limit)
	var response []*timelineResponse
	if err != nil {
		log.Println(err)
		return response
	}

	defer rows.Close()
	var res *timelineResponse
	for rows.Next() {
		var username string
		var userfeedid int
		var eventid int
		if err := rows.Scan(&username, &userfeedid, &eventid); err != nil {
			log.Println(err)
			return response
		}

		res = &timelineResponse{
			UserFeedID: userfeedid,
			Username:   username,
			EventID:    eventid,
		}
		response = append(response, res)

	}

	return response
}

// data, err := os.ReadFile("public/profile-picture/elephant-seal.jpg")
// if err != nil {
// 	log.Fatal(err)
// }

// insertDynStmt := `insert into events (trip_id, user_id, caption, event_date, post_time, picture) values($1,$2,$3,$4,$5,$6)`
// _, err = db.Exec(insertDynStmt, 1, 1, "Melihat gajah laut", time.Now(), time.Now(), data)

// if err != nil {
// 	log.Fatal(err)
// }

// return c.SendString("hi")
