package main

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Adhiana46/go-restapi-template/handlers"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func routes() *fiber.App {
	r := fiber.New(fiber.Config{
		ErrorHandler: handleError,
	})

	// Handle Panic
	r.Use(func(c *fiber.Ctx) error {
		log.Println("Handle Panic Middleware")

		defer handlePanic(c)

		return c.Next()
	})

	api := r.Group("/api/v1")

	routesGoogleSSO(api)

	// Register Handlers
	handlers.
		NewActivityGroupHandler(svcActivityGroup).
		RegisterRoutes(api.Group("/activity-group"))
	handlers.
		NewTodoItemHandler(svcTodoItem).
		RegisterRoutes(api.Group("/activity-group/:activity_uuid/todo-items"))

	return r
}

var randomString string = "random-string-secret"
var sso *oauth2.Config

func routesGoogleSSO(r fiber.Router) {
	sso = &oauth2.Config{
		RedirectURL:  cfg.GoogleSSO.RedirectURL,
		ClientID:     cfg.GoogleSSO.ClientID,
		ClientSecret: cfg.GoogleSSO.ClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	r.Get("oauth/google/signin", func(c *fiber.Ctx) error {
		url := sso.AuthCodeURL(randomString)
		log.Println("oauth/google/signin", url)

		return c.Redirect(url, http.StatusTemporaryRedirect)
	})

	r.Get("oauth/google/callback", func(c *fiber.Ctx) error {
		state := c.Query("state")
		code := c.Query("code")

		data, err := getUserData(state, code)
		if err != nil {
			log.Println("error getUserData: ", err)
		}

		log.Printf("userData: %s", data)

		return c.JSON(string(data))
	})
}

func getUserData(state, code string) ([]byte, error) {
	if state != randomString {
		return nil, errors.New("Invalid user state")
	}

	token, err := sso.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v3/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
