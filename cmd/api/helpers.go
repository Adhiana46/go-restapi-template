package main

import (
	"database/sql"
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/Adhiana46/go-restapi-template/pkg/monitoring"
	parserPkg "github.com/Adhiana46/go-restapi-template/pkg/parser"
	responsePkg "github.com/Adhiana46/go-restapi-template/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func handleError(c *fiber.Ctx, err error) error {
	var statusCode int = 500
	var message string = ""
	var errorsData any = nil

	if err == sql.ErrNoRows {
		statusCode = 404
	} else if strings.Contains(strings.ToLower(err.Error()), "query parameter") {
		statusCode = 400
		message = err.Error()
	} else {
		switch err.(type) {
		case validator.ValidationErrors:
			errs := err.(validator.ValidationErrors)

			statusCode = 400
			errorsData = parserPkg.ValidationErrors(errs, &validateTrans)
		default:
			statusCode = 500

			// log monitoring sentry.io
			logMonitoring.LogPanic(err.Error(), monitoring.LogData{
				Method:        string(c.Request().Header.Method()),
				Endpoint:      string(c.Request().URI().Path()),
				RequestBody:   string(c.Body()),
				RequestParams: c.AllParams(),
				RequestQuery:  string(c.Context().QueryArgs().QueryString()),
				StackTrace:    string(debug.Stack()),
			})
		}
	}

	// TODO: log errors
	log.Println("handleError: ", err)

	return c.Status(statusCode).JSON(responsePkg.JsonError(statusCode, message, errorsData))
}

func handlePanic(c *fiber.Ctx) {
	if r := recover(); r != nil {
		// TODO: log
		log.Println("Recovered in f", r, string(debug.Stack()))

		// log monitoring sentry.io
		logMonitoring.LogPanic(r.(string), monitoring.LogData{
			Method:        string(c.Request().Header.Method()),
			Endpoint:      string(c.Request().URI().Path()),
			RequestBody:   string(c.Body()),
			RequestParams: c.AllParams(),
			RequestQuery:  string(c.Context().QueryArgs().QueryString()),
			StackTrace:    string(debug.Stack()),
		})

		response := responsePkg.JsonError(http.StatusInternalServerError, "", nil)

		c.Status(http.StatusInternalServerError).JSON(response)
	}
}
