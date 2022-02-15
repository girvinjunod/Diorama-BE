package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "diorama"
	password = "diorama"
	dbname   = "diorama"
)

func main() {
	log.Println("Starting Server")
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

	fmt.Println("Succesfully connected to database")

	// insertDynStmt := `insert into users (username, email, password, profile_picture) values($1,$2,$3,$4)`
	// _, err = db.Exec(insertDynStmt, "scrooge", "scrooge@gmail.com", "scrooge", "profile-picture/rich.jpg")
	// CheckError(err)

	app.Get("/", func(c *fiber.Ctx) error {
		log.Println("Hello")
		return c.SendString("Hello, World 👋!")
	})

	app.Get("/getUserByID/:id", func(c *fiber.Ctx) error {
		type userResponse struct {
			Id              int    `json:"id"`
			Username        string `json:"username"`
			Email           string `json:"email"`
			Profile_picture string `json:"profile_picture"`
		}
		id := c.Params("id")
		query := `SELECT * FROM users where id=$1`
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
			var password string
			var profile_picture string
			if err := rows.Scan(&id, &username, &email, &password, &profile_picture); err != nil {
				log.Fatal(err)
			}

			res := &userResponse{
				Id:              id,
				Username:        username,
				Email:           email,
				Profile_picture: profile_picture,
			}

			reply, _ := json.Marshal(res)
			response = string(reply)

			log.Println(response)
		}
		return c.SendString(response)
	},
	)

	app.Listen(":3000")
}
