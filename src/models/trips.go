package models

import (
	"database/sql"
	"log"
	"os"
	"time"
)

func AddTrip(db *sql.DB, user_id int, start_date string, end_date string, trip_name string, location_name string) (string, int) {
	log.Println("Add trip")
	start, err := time.Parse("2006-01-02", start_date)
	if err != nil {
		log.Println(err)
		return err.Error(), 0
	}
	end, err := time.Parse("2006-01-02", end_date)
	if err != nil {
		log.Println(err)
		return err.Error(), 0
	}
	insertDynStmt := `insert into trips (user_id, start_date, end_date, trip_name, location_name) values($1,$2,$3,$4,$5) RETURNING id`
	var id int
	err = db.QueryRow(insertDynStmt, user_id, start, end, trip_name, location_name).Scan(&id)

	if err != nil {
		log.Println(err)
		return err.Error(), 0
	}

	log.Printf("Inserted trip on id = %d", id)
	return "true", id
}

type tripResponse struct {
	Error        string `json:"error"`
	Id           int    `json:"id"`
	UserID       int    `json:"userId"`
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
	TripName     string `json:"tripName"`
	LocationName string `json:"locationName"`
}

func GetTripDetailById(db *sql.DB, id string) *tripResponse {
	query := `SELECT id, user_id, start_date, end_date, trip_name, location_name FROM trips where id=$1`
	rows, err := db.Query(query, id)
	var res *tripResponse
	if err != nil {
		return res
	}

	defer rows.Close()
	for rows.Next() {
		var id int
		var userid int
		var startDate time.Time
		var endDate time.Time
		var tripName string
		var locationName string
		if err := rows.Scan(&id, &userid, &startDate, &endDate, &tripName, &locationName); err != nil {
			log.Println(err)
			return res
		}

		res = &tripResponse{
			Error:        "false",
			Id:           id,
			UserID:       userid,
			StartDate:    startDate.Format("2006-01-02"),
			EndDate:      endDate.Format("2006-01-02"),
			TripName:     tripName,
			LocationName: locationName,
		}

	}
	return res
}

func GetAllEventsFromTrip(db *sql.DB, id string) (string, []int) {
	query := `SELECT e.id as eventID FROM trips t, events e where t.id=$1 and t.id=e.trip_id`
	rows, err := db.Query(query, id)
	var res []int
	if err != nil {
		return err.Error(), res
	}

	defer rows.Close()
	for rows.Next() {
		var eventid int
		if err := rows.Scan(&eventid); err != nil {
			log.Println(err)
			return err.Error(), res
		}

		res = append(res, eventid)

	}
	return id, res
}

func SetTripDetail(db *sql.DB, tripID string, start_date string, end_date string, trip_name string, location_name string) string {
	log.Printf("update trip with ID= " + tripID)
	insertDynStmt := `UPDATE trips SET
    start_date = $1,
    end_date = $2,
	trip_name = $3,
	location_name = $4
	WHERE id = $5`
	_, err := db.Exec(insertDynStmt, start_date, end_date, trip_name, location_name, tripID)

	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return "true"
}

func DeleteTrip(db *sql.DB, tripID string) string {
	log.Printf("delete trip with ID= " + tripID)
	query := `delete from trips where id=$1`
	_, err := db.Exec(query, tripID)

	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return "true"
}

func GetTripsByUser(db *sql.DB, user_id string) (string, []int) {
	query := `SELECT id FROM trips where user_id=$1`
	rows, err := db.Query(query, user_id)
	var res []int
	if err != nil {
		return err.Error(), res
	}

	defer rows.Close()
	for rows.Next() {
		var tripid int
		if err := rows.Scan(&tripid); err != nil {
			log.Println(err)
			return err.Error(), res
		}

		res = append(res, tripid)

	}
	return user_id, res
}

func GetTripsImage(db *sql.DB, id string) []byte {
	log.Println("Get event picture with ID " + id)
	query := `SELECT picture FROM events where trip_id=$1 LIMIT 1`
	rows, err := db.Query(query, id)
	var response []byte
	if err != nil {
		log.Println(err)
		return response
	}
	found := false
	defer rows.Close()
	for rows.Next() {
		var picture []byte
		if err := rows.Scan(&picture); err != nil {
			log.Println(err)
			return response
		}

		response = picture
		found = true
	}

	if !found {
		log.Println("No picture found")
		data, err := os.ReadFile("../public/placeholder/event.jpg")
		if err != nil {
			log.Fatal(err)
		}
		response = data
	}

	return response
}
