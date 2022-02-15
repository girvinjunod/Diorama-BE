package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

var (
	host     = goDotEnvVariable("HOST")
	port     = 5432
	user     = goDotEnvVariable("USER")
	password = goDotEnvVariable("PASSWORD")
	dbname   = goDotEnvVariable("DBNAME")
)

func main() {
	log.Println("Starting server on " + host)
	app := fiber.New()

	app.Static("/", "./public")

	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// check db
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Succesfully connected to database")

	app.Get("/", func(c *fiber.Ctx) error {
		log.Println("Hello")
		return c.SendString("Hello, World 👋!")
	})

	app.Get("/getUserByID/:id", func(c *fiber.Ctx) error {
		type userResponse struct {
			Id       int    `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
		}
		id := c.Params("id")
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
		return c.SendString(response)
	},
	)

	app.Get("/getPPByID/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
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
		return c.Send(response)
	},
	)

	//TODO
	app.Post("/addTrip", func(c *fiber.Ctx) error {
		insertDynStmt := `insert into trips (user_id, start_date, end_date, trip_name, location_name) values($1,$2,$3,$4,$5)`
		_, err = db.Exec(insertDynStmt, "1", time.Now(), time.Now().AddDate(0, 0, 10), "Jalan-jalan ke Bandung", "ITB")

		if err != nil {
			log.Fatal(err)
		}

		return c.SendString("hi")
	})

	//TODO
	app.Post("/addEvent", func(c *fiber.Ctx) error {
		data, err := os.ReadFile("public/profile-picture/elephant-seal.jpg")
		if err != nil {
			log.Fatal(err)
		}

		insertDynStmt := `insert into events (trip_id, user_id, caption, event_date, post_time, picture) values($1,$2,$3,$4,$5,$6)`
		_, err = db.Exec(insertDynStmt, 1, 1, "Melihat gajah laut", time.Now(), time.Now(), data)

		if err != nil {
			log.Fatal(err)
		}

		return c.SendString("hi")
	})

	app.Listen(":3000")
}
