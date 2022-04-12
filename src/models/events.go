package models

import (
	"database/sql"
	"log"
	"time"
)

func AddEvent(db *sql.DB, trip_id string, user_id string, caption string, event_date string, picture []byte) (string, int) {
	log.Println("Add new event")
	event, err := time.Parse("2006-01-02", event_date)
	if err != nil {
		log.Println(err)
		return err.Error(), 0
	}

	insertDynStmt := `insert into events (trip_id, user_id, caption, event_date, post_time, picture) values($1,$2,$3,$4,$5,$6) RETURNING id`
	var id int
	err = db.QueryRow(insertDynStmt, trip_id, user_id, caption, event, time.Now(), picture).Scan(&id)

	if err != nil {
		log.Println(err)
		return err.Error(), 0
	}

	log.Printf("Inserted events on id = %d", id)
	return "true", id
}

type eventResponse struct {
	Error     string `json:"error"`
	Id        int    `json:"id"`
	TripID    int    `json:"tripId"`
	UserID    int    `json:"userId"`
	Caption   string `json:"caption"`
	EventDate string `json:"eventDate"`
	PostTime  string `json:"postTime"`
}

func GetEventDetailByID(db *sql.DB, id string) *eventResponse {
	query := `SELECT id, trip_id, user_id, caption, event_date, post_time FROM events where id=$1`
	rows, err := db.Query(query, id)
	var res *eventResponse
	if err != nil {
		return res
	}

	defer rows.Close()
	for rows.Next() {
		var id int
		var tripid int
		var userid int
		var caption string
		var eventDate time.Time
		var postTime time.Time
		if err := rows.Scan(&id, &tripid, &userid, &caption, &eventDate, &postTime); err != nil {
			log.Println(err)
			return res
		}

		res = &eventResponse{
			Error:     "false",
			Id:        id,
			TripID:    tripid,
			UserID:    userid,
			Caption:   caption,
			EventDate: eventDate.Format("2006-01-02"),
			PostTime:  postTime.Format("2006-01-02 15:04:05"),
		}

	}
	return res
}

func SetEventDetail(db *sql.DB, eventID string, caption string, event_date string) string {
	log.Printf("update event with ID= " + eventID)
	insertDynStmt := `UPDATE events SET
    caption = $1,
    event_date = $2
	WHERE id = $3`
	_, err := db.Exec(insertDynStmt, caption, event_date, eventID)

	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return "true"
}

func SetEventPicture(db *sql.DB, picture []byte, eventID string) string {
	log.Println("Change picture at event ID=" + eventID)
	insertDynStmt := `update events set picture=$1 where id=$2`
	_, err := db.Exec(insertDynStmt, picture, eventID)

	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return "true"
}

func GetEventPictureByID(db *sql.DB, id string) []byte {
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

func DeleteEvent(db *sql.DB, eventID string) string {
	log.Printf("delete event with ID= " + eventID)
	query := `delete from events where id=$1`
	_, err := db.Exec(query, eventID)

	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return "true"
}

type timelineResponse struct {
	UserFeedID int    `json:"userID"`
	Username   string `json:"username"`
	EventID    int    `json:"eventID"`
	Caption    string `json:"caption"`
	TripName   string `json:"tripname"`
}

func GetTimeline(db *sql.DB, userID string) []*timelineResponse {
	log.Println("Get Timeline")
	query := `SELECT u.username as Username, f.followed_id as UserFeedID, e.id as EventID, e.caption as Caption, t.trip_name FROM users u, following f, events e, trips t where f.follower_id = $1 and u.id=f.followed_id and f.followed_id=e.user_id and e.trip_id = t.id ORDER BY e.post_time desc`
	rows, err := db.Query(query, userID)
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
		var caption string
		var trip_name string
		if err := rows.Scan(&username, &userfeedid, &eventid, &caption, &trip_name); err != nil {
			log.Println(err)
			return response
		}

		res = &timelineResponse{
			UserFeedID: userfeedid,
			Username:   username,
			EventID:    eventid,
			Caption:    caption,
			TripName:   trip_name,
		}
		response = append(response, res)

	}

	return response
}
