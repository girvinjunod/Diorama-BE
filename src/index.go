package main

import (
	"database/sql"
	"diorama/v2/auth"
	"diorama/v2/models"
	"diorama/v2/utils"
	"fmt"
	"io/ioutil"
	"os"

	"log"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
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
	port     = goDotEnvVariable("PQ_PORT")
	user     = goDotEnvVariable("PQ_USER")
	password = goDotEnvVariable("PQ_PASSWORD")
	dbname   = goDotEnvVariable("PQ_DBNAME")
)

func main() {
	serverport := os.Getenv("PORT")
	secretkey := goDotEnvVariable("SECRET_KEY")
	// log.Println(secret_key)
	// cnxn := "postgres://sdrgqiodobzvzq:0ec897ee53f52a65f994301d697abe14f5cac794844ebb127adef380513f0c4d@ec2-3-209-124-113.compute-1.amazonaws.com:5432/dbb0rrl7sa5hb4"
	// log.Println("Starting server on " + host)
	app := fiber.New()

	app.Static("/public", "../public")

	// // connection string
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	// log.Println(psqlconn)
	// open database
	// db, err := sql.Open("postgres", cnxn)
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

	//unrestricted routes
	app.Get("/", func(c *fiber.Ctx) error {
		log.Println("Hello")
		return utils.SuccessMsg(c, "Hello World!")
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
			return utils.ErrorMsg(c, err.Error())
		}
		res := auth.Register(db, p.Username, p.Email, p.Name, p.Password)
		if res == "true" {
			return utils.SuccessMsg(c, "Successfully registered user")
		}
		return utils.ErrorMsg(c, res)

	})

	app.Post("/login", func(c *fiber.Ctx) error {
		type LoginRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		req := new(LoginRequest)
		if err := c.BodyParser(req); err != nil {
			return utils.ErrorMsg(c, err.Error())
		}

		res, id := auth.Login(db, req.Username, req.Password)

		if res == "true" {
			token, err := auth.CreateJWTToken(req.Username)
			if err != nil {
				return utils.ErrorMsg(c, err.Error())
			}
			return c.JSON(fiber.Map{"token": token, "user": req.Username, "user_id": id})
		}
		return utils.ErrorMsg(c, res)

	})

	app.Get("/getPPByID/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		response := models.GetPPByID(db, id)
		// if len(response) == 0 {
		// 	return errorMsg(c, "No picture found")
		// }
		return c.Status(fiber.StatusOK).Send(response)
	},
	)
	app.Get("/getEventPictureByID/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		response := models.GetEventPictureByID(db, id)
		// if len(response) == 0 {
		// 	return errorMsg(c, "No picture found")
		// }
		return c.Status(fiber.StatusOK).Send(response)
	},
	)

	app.Get("/getTripsImage/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		response := models.GetTripsImage(db, id)
		return c.Status(fiber.StatusOK).Send(response)
	},
	)
	//Restricted Routes
	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(secretkey),
	}))

	// User API

	app.Get("/getUserByID/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		response := models.GetUserById(db, id)

		if response != nil {
			return c.Status(fiber.StatusOK).JSON(response)

		}
		return utils.ErrorMsg(c, "User not found")

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
			return utils.ErrorMsg(c, err.Error())
		}

		res := models.SetUserDetail(db, id, p.Username, p.Name, p.Email)
		if res == "true" {
			return utils.SuccessMsg(c, "Successfully updated user details")
		}
		return utils.ErrorMsg(c, res)

	})

	app.Post("/setUserPassword", func(c *fiber.Ctx) error {
		// log.Println(c)
		// log.Println(c.Request())
		// log.Println(string(c.Body()[:]))
		type User struct {
			UserID      int    `json:"userID"`
			OldPassword string `json:"OldPassword"`
			NewPassword string `json:"NewPassword"`
		}
		p := new(User)
		if err := c.BodyParser(p); err != nil {
			return utils.ErrorMsg(c, err.Error())
		}

		res := models.SetUserPassword(db, p.UserID, p.OldPassword, p.NewPassword)
		if res == "true" {
			return utils.SuccessMsg(c, "Successfully updated user password")
		}
		return utils.ErrorMsg(c, res)

	})

	app.Put("/setUserPP/:id", func(c *fiber.Ctx) error {
		userID := c.Params("id")
		file, err := c.FormFile("picture")
		if err != nil {
			log.Println("File is bytes")
			file := c.FormValue("picture")
			if file == "" {
				log.Println("Can't get picture")
				return utils.ErrorMsg(c, err.Error())
			}
			buffer := []byte(file)
			res := models.SetUserPP(db, buffer, userID)

			if res == "true" {
				return utils.SuccessMsg(c, "Successfully changed profile picture")
			}
			return utils.ErrorMsg(c, res)

		}
		log.Println("File is file")
		buffer, err := file.Open()
		if err != nil {
			return utils.ErrorMsg(c, err.Error())
		}
		defer buffer.Close()

		data, err := ioutil.ReadAll(buffer)
		if err != nil {
			utils.ErrorMsg(c, err.Error())
		}

		res := models.SetUserPP(db, data, userID)

		if res == "true" {
			return utils.SuccessMsg(c, "Successfully changed profile picture")
		}
		return utils.ErrorMsg(c, res)

	})

	app.Delete("/deleteUser/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		response := models.DeleteUser(db, id)

		if response == "true" {
			return utils.SuccessMsg(c, "User successfully deleted")
		}
		return utils.ErrorMsg(c, response)

	})

	app.Get("/searchUser/:query", func(c *fiber.Ctx) error {
		query := c.Params("query")
		response := models.SearchUser(db, query)

		if response != nil {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"error": false,
				"users": response,
			})
		}
		return utils.ErrorMsg(c, "Users not found")

	},
	)

	//Follow API
	app.Put("/follow/:followerid/:followedid", func(c *fiber.Ctx) error {
		followerid := c.Params("followerid")
		followedid := c.Params("followedid")
		response := models.Follow(db, followerid, followedid)

		if response == "true" {
			return utils.SuccessMsg(c, "Successfully followed")
		}
		return utils.ErrorMsg(c, response)

	})

	app.Delete("/unfollow/:followerid/:followedid", func(c *fiber.Ctx) error {
		followerid := c.Params("followerid")
		followedid := c.Params("followedid")
		response := models.Unfollow(db, followerid, followedid)

		if response == "true" {
			return utils.SuccessMsg(c, "Successfully unfollowed")
		}
		return utils.ErrorMsg(c, response)

	})

	app.Get("/getFollowedUsers/:id", func(c *fiber.Ctx) error {
		userID := c.Params("id")
		response := models.GetAllFollowedUsers(db, userID)

		if len(response) == 0 {
			return utils.ErrorMsg(c, "No followed user found")
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"followed_users": response,
			"error":          false,
		})

	})

	app.Get("/getFollowers/:id", func(c *fiber.Ctx) error {
		userID := c.Params("id")
		response := models.GetAllFollowers(db, userID)

		if len(response) == 0 {
			return utils.ErrorMsg(c, "No followed user found")
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"followed_users": response,
			"error":          false,
		})

	})

	app.Get("/checkIfFollowed/:followerid/:followedid", func(c *fiber.Ctx) error {
		followerid := c.Params("followerid")
		followedid := c.Params("followedid")
		err, response := models.CheckIfFollowed(db, followerid, followedid)

		if err == "" {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"is_followed": response,
				"error":       false,
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg":   "Failed to check",
			"error": true,
		})

	})

	app.Get("/getCountFollowers/:id", func(c *fiber.Ctx) error {
		userID := c.Params("id")
		response := models.GetCountFollower(db, userID)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"count": response,
			"error": false,
		})
	})

	app.Get("/getCountFollowing/:id", func(c *fiber.Ctx) error {
		userID := c.Params("id")
		response := models.GetCountFollowing(db, userID)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"count": response,
			"error": false,
		})
	})

	// Trips API

	app.Post("/addTrip", func(c *fiber.Ctx) error {
		type Trip struct {
			UserID       int    `json:"UserID"`
			StartDate    string `json:"StartDate"`
			EndDate      string `json:"EndDate"`
			TripName     string `json:"TripName"`
			LocationName string `json:"LocationName"`
		}
		p := new(Trip)
		if err := c.BodyParser(p); err != nil {
			return utils.ErrorMsg(c, err.Error())
		}

		res, id := models.AddTrip(db, p.UserID, p.StartDate, p.EndDate, p.TripName, p.LocationName)
		if res == "true" {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"TripID": id,
				"error":  false,
			})
		}
		return utils.ErrorMsg(c, res)

	})

	app.Get("/getTripDetailByID/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		response := models.GetTripDetailById(db, id)

		if response != nil {
			return c.Status(fiber.StatusOK).JSON(response)

		}
		return utils.ErrorMsg(c, "Trip not found")

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
			return utils.ErrorMsg(c, err.Error())
		}

		res := models.SetTripDetail(db, id, p.StartDate, p.EndDate, p.TripName, p.LocationName)
		if res == "true" {
			return utils.SuccessMsg(c, "Successfully updated trip details")
		}
		return utils.ErrorMsg(c, res)

	})

	app.Delete("/deleteTrip/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		response := models.DeleteTrip(db, id)

		if response == "true" {
			return utils.SuccessMsg(c, "Trip successfully deleted")
		}
		return utils.ErrorMsg(c, response)

	})

	app.Get("/getTripsByUser/:user_id", func(c *fiber.Ctx) error {
		userID := c.Params("user_id")
		id, res := models.GetTripsByUser(db, userID)

		if res != nil {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"error":   false,
				"user_id": id,
				"tripIds": res,
			})
		}
		return utils.ErrorMsg(c, "Trip not found")

	},
	)

	// Event API

	app.Post("/addEvent", func(c *fiber.Ctx) error {
		log.Println("Add Event")
		TripID := c.FormValue("tripID")
		UserID := c.FormValue("userID")
		Caption := c.FormValue("caption")
		EventDate := c.FormValue("eventDate")
		file, err := c.FormFile("picture")
		if err != nil {
			log.Println("File is bytes")
			file := c.FormValue("picture")
			if file == "" {
				log.Println("Can't get picture")
				return utils.ErrorMsg(c, err.Error())
			}
			buffer := []byte(file)
			res, id := models.AddEvent(db, TripID, UserID, Caption, EventDate, buffer)
			if res == "true" {
				return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"EventID": id,
					"error":   false,
				})
			}
			return utils.ErrorMsg(c, res)

		}
		log.Println("File is file")
		buffer, err := file.Open()
		if err != nil {
			return utils.ErrorMsg(c, err.Error())
		}
		defer buffer.Close()

		data, err := ioutil.ReadAll(buffer)
		if err != nil {
			utils.ErrorMsg(c, err.Error())
		}

		res, id := models.AddEvent(db, TripID, UserID, Caption, EventDate, data)
		if res == "true" {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"EventID": id,
				"error":   false,
			})
		}
		return utils.ErrorMsg(c, res)

	})

	app.Get("/getEventDetailByID/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		response := models.GetEventDetailByID(db, id)

		if response != nil {
			return c.Status(fiber.StatusOK).JSON(response)

		}
		return utils.ErrorMsg(c, "Event not found")

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
			return utils.ErrorMsg(c, err.Error())
		}

		res := models.SetEventDetail(db, id, p.Caption, p.EventDate)
		if res == "true" {
			return utils.SuccessMsg(c, "Successfully updated event details")
		}
		return utils.ErrorMsg(c, res)

	})

	app.Put("/setEventPicture/:id", func(c *fiber.Ctx) error {
		file, err := c.FormFile("picture")
		if err != nil {
			return utils.ErrorMsg(c, err.Error())
		}
		buffer, err := file.Open()
		if err != nil {
			return utils.ErrorMsg(c, err.Error())
		}
		defer buffer.Close()

		data, err := ioutil.ReadAll(buffer)
		if err != nil {
			utils.ErrorMsg(c, err.Error())
		}

		eventID := c.Params("id")

		res := models.SetEventPicture(db, data, eventID)

		if res == "true" {
			return utils.SuccessMsg(c, "Successfully changed picture")
		}
		return utils.ErrorMsg(c, res)

	})

	app.Get("/getEventsFromTrip/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		tripID, events := models.GetAllEventsFromTrip(db, id)

		if events != nil {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"error":    false,
				"tripID":   tripID,
				"eventIDs": events,
			})
		}
		return utils.ErrorMsg(c, "Trip or Events not found")

	})

	app.Get("/getTimeline/:id", func(c *fiber.Ctx) error {
		userID := c.Params("id")
		response := models.GetTimeline(db, userID)

		if len(response) == 0 {
			return utils.ErrorMsg(c, "Empty timeline")
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"timeline_data": response,
			"error":         false,
		})

	})

	app.Delete("/deleteEvent/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		response := models.DeleteEvent(db, id)

		if response == "true" {
			return utils.SuccessMsg(c, "Event successfully deleted")
		}
		return utils.ErrorMsg(c, response)

	})

	// Comment API

	app.Post("/addComment", func(c *fiber.Ctx) error {
		type Comment struct {
			EventID int    `json:"EventID"`
			UserID  int    `json:"UserID"`
			Text    string `json:"Text"`
		}
		p := new(Comment)
		if err := c.BodyParser(p); err != nil {
			return utils.ErrorMsg(c, err.Error())
		}

		res, id := models.AddComment(db, p.EventID, p.UserID, p.Text)
		if res == "true" {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"CommentID": id,
				"error":     false,
			})
		}
		return utils.ErrorMsg(c, res)

	})

	app.Get("/getCommentDetailByID/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		response := models.GetCommentDetailById(db, id)

		if response != nil {
			return c.Status(fiber.StatusOK).JSON(response)

		}
		return utils.ErrorMsg(c, "Comment not found")

	},
	)

	app.Put("/setCommentDetail/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		type Comment struct {
			Text string `json:"Text"`
		}
		p := new(Comment)
		if err := c.BodyParser(p); err != nil {
			return utils.ErrorMsg(c, err.Error())
		}

		res := models.SetCommentDetail(db, id, p.Text)
		if res == "true" {
			return utils.SuccessMsg(c, "Successfully updated comment details")
		}
		return utils.ErrorMsg(c, res)

	})

	app.Delete("/deleteComment/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		response := models.DeleteComment(db, id)

		if response == "true" {
			return utils.SuccessMsg(c, "Comment successfully deleted")
		}
		return utils.ErrorMsg(c, response)

	})

	app.Get("/getCommentsFromEvent/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		eventID, comments := models.GetAllCommentsFromEvent(db, id)

		if comments != nil {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"error":    false,
				"eventID":  eventID,
				"comments": comments,
			})
		}
		return utils.ErrorMsg(c, "Comments or event not found")

	})

	if serverport != "" {
		fmt.Println("Server is running on port: " + port)
		app.Listen(":" + port)
	} else {
		fmt.Println("Server is running on port: 8080")
		app.Listen(":8080")
	}
}
