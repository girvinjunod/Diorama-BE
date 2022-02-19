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

	app.Post("/setUserDetail", func(c *fiber.Ctx) error {
		type User struct {
			UserId   int    `json:"userID"`
			Username string `json:"Username"`
			Name     string `json:"Name"`
			Email    string `json:"Email"`
		}
		p := new(User)
		if err := c.BodyParser(p); err != nil {
			return errorMsg(c, err.Error())
		}

		if p.UserId == 0 {
			return errorMsg(c, "Invalid user ID")
		}

		res := setUserDetail(db, p.UserId, p.Username, p.Name, p.Email)
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

		if p.UserId == 0 {
			return errorMsg(c, "Invalid user ID")
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

	app.Post("/setUserPP", func(c *fiber.Ctx) error {
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

		userID := c.FormValue("userID")

		res := setUserPP(db, data, userID)

		if res == "true" {
			return successMsg(c, "Successfully added profile picture")
		} else {
			return errorMsg(c, res)
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
		if p.UserId == 0 {
			return errorMsg(c, "Invalid user ID")
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

	app.Post("/setTripDetail", func(c *fiber.Ctx) error {
		type Trip struct {
			TripId       int    `json:"TripID"`
			StartDate    string `json:"StartDate"`
			EndDate      string `json:"EndDate"`
			TripName     string `json:"TripName"`
			LocationName string `json:"LocationName"`
		}
		p := new(Trip)
		if err := c.BodyParser(p); err != nil {
			return errorMsg(c, err.Error())
		}

		if p.TripId == 0 {
			return errorMsg(c, "Invalid user ID")
		}

		res := setTripDetail(db, p.TripId, p.StartDate, p.EndDate, p.TripName, p.LocationName)
		if res == "true" {
			return successMsg(c, "Successfully updated trip details")
		} else {
			return errorMsg(c, res)
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
		PostTime := c.FormValue("postTime")

		res, id := addEvent(db, TripId, UserId, Caption, EventDate, PostTime, data)
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

	app.Post("/setEventDetail", func(c *fiber.Ctx) error {
		type Event struct {
			EventId   int    `json:"EventID"`
			Caption   string `json:"Caption"`
			EventDate string `json:"EventDate"`
		}
		p := new(Event)
		if err := c.BodyParser(p); err != nil {
			return errorMsg(c, err.Error())
		}

		if p.EventId == 0 {
			return errorMsg(c, "Invalid event ID")
		}

		res := setEventDetail(db, p.EventId, p.Caption, p.EventDate)
		if res == "true" {
			return successMsg(c, "Successfully updated event details")
		} else {
			return errorMsg(c, res)
		}

	})

	app.Post("/setEventPicture", func(c *fiber.Ctx) error {
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

		eventID := c.FormValue("eventID")

		res := setEventPicture(db, data, eventID)

		if res == "true" {
			return successMsg(c, "Successfully added picture")
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
