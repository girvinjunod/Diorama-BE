package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

var (
	host     = goDotEnvVariable("PQ_HOST")
	port     = 5432
	user     = goDotEnvVariable("PQ_USER")
	password = goDotEnvVariable("PQ_PASSWORD")
	dbname   = goDotEnvVariable("PQ_DBNAME")
)

func main() {
	log.Println("Starting server on " + host)
	app := fiber.New()

	app.Static("/public", "../public")

	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	// log.Println(psqlconn)
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
		return successMsg(c, "Hello World!")
	})

	// Auth API

	app.Post("/register", func(c *fiber.Ctx) error {
		type User struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Name     string `json:"name"`
			Password string `json:"password"`
		}
		p := new(User)
		if err := c.BodyParser(p); err != nil {
			return errorMsg(c, err.Error())
		}
		res := register(db, p.Username, p.Email, p.Name, p.Password)
		if res == "true" {
			return successMsg(c, "Successfully registered user")
		} else {
			return errorMsg(c, res)
		}
	})

	// User API

	app.Get("/getUserByID/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		response := getUserById(db, id)

		if response != nil {
			return c.Status(fiber.StatusOK).JSON(response)

		} else {
			return errorMsg(c, "User not found")
		}
	},
	)

	app.Put("/setUserDetail/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		type User struct {
			Username string `json:"Username"`
			Name     string `json:"Name"`
			Email    string `json:"Email"`
		}
		p := new(User)
		if err := c.BodyParser(p); err != nil {
			return errorMsg(c, err.Error())
		}

		res := setUserDetail(db, id, p.Username, p.Name, p.Email)
		if res == "true" {
			return successMsg(c, "Successfully updated user details")
		} else {
			return errorMsg(c, res)
		}
	})

	app.Post("/setUserPassword", func(c *fiber.Ctx) error {
		type User struct {
			UserId      int    `json:"userID"`
			OldPassword string `json:"OldPassword"`
			NewPassword string `json:"NewPassword"`
		}
		p := new(User)
		if err := c.BodyParser(p); err != nil {
			return errorMsg(c, err.Error())
		}

		res := setUserPassword(db, p.UserId, p.OldPassword, p.NewPassword)
		if res == "true" {
			return successMsg(c, "Successfully updated user password")
		} else {
			return errorMsg(c, res)
		}
	})

	app.Get("/getPPByID/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		response := getPPByID(db, id)
		// if len(response) == 0 {
		// 	return errorMsg(c, "No picture found")
		// }
		return c.Status(fiber.StatusOK).Send(response)
	},
	)

	app.Put("/setUserPP/:id", func(c *fiber.Ctx) error {
		file, err := c.FormFile("picture")
		if err != nil {
			return errorMsg(c, err.Error())
		}
		buffer, err := file.Open()
		if err != nil {
			return errorMsg(c, err.Error())
		}
		defer buffer.Close()

		data, err := ioutil.ReadAll(buffer)
		if err != nil {
			errorMsg(c, err.Error())
		}

		userID := c.Params("id")

		res := setUserPP(db, data, userID)

		if res == "true" {
			return successMsg(c, "Successfully changed profile picture")
		} else {
			return errorMsg(c, res)
		}
	})

	app.Delete("/deleteUser/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		response := deleteUser(db, id)

		if response == "true" {
			return successMsg(c, "User successfully deleted")
		} else {
			return errorMsg(c, response)
		}
	})

	// Trips API

	app.Post("/addTrip", func(c *fiber.Ctx) error {
		type Trip struct {
			UserId       int    `json:"UserID"`
			StartDate    string `json:"StartDate"`
			EndDate      string `json:"EndDate"`
			TripName     string `json:"TripName"`
			LocationName string `json:"LocationName"`
		}
		p := new(Trip)
		if err := c.BodyParser(p); err != nil {
			return errorMsg(c, err.Error())
		}

		res, id := addTrip(db, p.UserId, p.StartDate, p.EndDate, p.TripName, p.LocationName)
		if res == "true" {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"TripID": id,
				"error":  false,
			})
		} else {
			return errorMsg(c, res)
		}
	})

	app.Get("/getTripDetailByID/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		response := getTripDetailById(db, id)

		if response != nil {
			return c.Status(fiber.StatusOK).JSON(response)

		} else {
			return errorMsg(c, "Trip not found")
		}
	},
	)

	app.Put("/setTripDetail/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		type Trip struct {
			StartDate    string `json:"StartDate"`
			EndDate      string `json:"EndDate"`
			TripName     string `json:"TripName"`
			LocationName string `json:"LocationName"`
		}
		p := new(Trip)
		if err := c.BodyParser(p); err != nil {
			return errorMsg(c, err.Error())
		}

		res := setTripDetail(db, id, p.StartDate, p.EndDate, p.TripName, p.LocationName)
		if res == "true" {
			return successMsg(c, "Successfully updated trip details")
		} else {
			return errorMsg(c, res)
		}

	})

	app.Delete("/deleteTrip/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		response := deleteTrip(db, id)

		if response == "true" {
			return successMsg(c, "Trip successfully deleted")
		} else {
			return errorMsg(c, response)
		}
	})

	// Event API

	app.Post("/addEvent", func(c *fiber.Ctx) error {
		file, err := c.FormFile("picture")
		if err != nil {
			return errorMsg(c, err.Error())
		}
		buffer, err := file.Open()
		if err != nil {
			return errorMsg(c, err.Error())
		}
		defer buffer.Close()

		data, err := ioutil.ReadAll(buffer)
		if err != nil {
			errorMsg(c, err.Error())
		}

		TripId := c.FormValue("tripID")
		UserId := c.FormValue("userID")
		Caption := c.FormValue("caption")
		EventDate := c.FormValue("eventDate")

		res, id := addEvent(db, TripId, UserId, Caption, EventDate, data)
		if res == "true" {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"EventID": id,
				"error":   false,
			})
		} else {
			return errorMsg(c, res)
		}
	})

	app.Get("/getEventDetailByID/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		response := getEventDetailByID(db, id)

		if response != nil {
			return c.Status(fiber.StatusOK).JSON(response)

		} else {
			return errorMsg(c, "Event not found")
		}
	},
	)

	app.Put("/setEventDetail/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		type Event struct {
			Caption   string `json:"Caption"`
			EventDate string `json:"EventDate"`
		}
		p := new(Event)
		if err := c.BodyParser(p); err != nil {
			return errorMsg(c, err.Error())
		}

		res := setEventDetail(db, id, p.Caption, p.EventDate)
		if res == "true" {
			return successMsg(c, "Successfully updated event details")
		} else {
			return errorMsg(c, res)
		}

	})

	app.Put("/setEventPicture/:id", func(c *fiber.Ctx) error {
		file, err := c.FormFile("picture")
		if err != nil {
			return errorMsg(c, err.Error())
		}
		buffer, err := file.Open()
		if err != nil {
			return errorMsg(c, err.Error())
		}
		defer buffer.Close()

		data, err := ioutil.ReadAll(buffer)
		if err != nil {
			errorMsg(c, err.Error())
		}

		eventID := c.Params("id")

		res := setEventPicture(db, data, eventID)

		if res == "true" {
			return successMsg(c, "Successfully changed picture")
		} else {
			return errorMsg(c, res)
		}

	})

	app.Get("/getEventPictureByID/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		response := getEventPictureByID(db, id)
		// if len(response) == 0 {
		// 	return errorMsg(c, "No picture found")
		// }
		return c.Status(fiber.StatusOK).Send(response)
	},
	)

	app.Get("/getEventsFromTrip/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		tripId, events := getAllEventsFromTrip(db, id)

		if events != nil {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"error":    false,
				"tripID":   tripId,
				"eventIDs": events,
			})
		} else {
			return errorMsg(c, "Trip or Events not found")
		}
	})

	app.Get("/getTimeline/:id", func(c *fiber.Ctx) error {
		userID := c.Params("id")
		response := getTimeline(db, userID, 10)

		if len(response) == 0 {
			return errorMsg(c, "Empty timeline")
		} else {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"timeline_data": response,
				"error":         false,
			})
		}

	})

	app.Delete("/deleteEvent/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		response := deleteEvent(db, id)

		if response == "true" {
			return successMsg(c, "Event successfully deleted")
		} else {
			return errorMsg(c, response)
		}
	})

	// Comment API

	app.Post("/addComment", func(c *fiber.Ctx) error {
		type Comment struct {
			EventId int    `json:"EventID"`
			UserId  int    `json:"UserID"`
			Text    string `json:"Text"`
		}
		p := new(Comment)
		if err := c.BodyParser(p); err != nil {
			return errorMsg(c, err.Error())
		}

		res, id := addComment(db, p.EventId, p.UserId, p.Text)
		if res == "true" {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"CommentID": id,
				"error":     false,
			})
		} else {
			return errorMsg(c, res)
		}
	})

	app.Get("/getCommentDetailByID/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		response := getCommentDetailById(db, id)

		if response != nil {
			return c.Status(fiber.StatusOK).JSON(response)

		} else {
			return errorMsg(c, "Comment not found")
		}
	},
	)

	app.Put("/setCommentDetail/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		type Comment struct {
			Text string `json:"Text"`
		}
		p := new(Comment)
		if err := c.BodyParser(p); err != nil {
			return errorMsg(c, err.Error())
		}

		res := setCommentDetail(db, id, p.Text)
		if res == "true" {
			return successMsg(c, "Successfully updated comment details")
		} else {
			return errorMsg(c, res)
		}

	})

	app.Delete("/deleteComment/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		response := deleteComment(db, id)

		if response == "true" {
			return successMsg(c, "Comment successfully deleted")
		} else {
			return errorMsg(c, response)
		}
	})

	app.Get("/getCommentsFromEvent/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		eventId, comments := getAllCommentsFromEvent(db, id)

		if comments != nil {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"error":        false,
				"eventID":      eventId,
				"commentTexts": comments,
			})
		} else {
			return errorMsg(c, "Comments or event not found")
		}
	})

	app.Listen(":3000")
}

func errorMsg(c *fiber.Ctx, err string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": true,
		"msg":   err,
	})
}

func successMsg(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   msg,
	})
}
