package main

import (
	"database/sql"
	"log"
	"time"
)

func addTrip(db *sql.DB, user_id int, start_date string, end_date string, trip_name string, location_name string) string {
	log.Println("Add trip")
	start, _ := time.Parse("2006-01-02", start_date)
	end, _ := time.Parse("2006-01-02", end_date)
	insertDynStmt := `insert into trips (user_id, start_date, end_date, trip_name, location_name) values($1,$2,$3,$4,$5)`
	_, err := db.Exec(insertDynStmt, user_id, start, end, trip_name, location_name)

	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return "true"
}
