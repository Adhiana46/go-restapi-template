package main

import (
	"database/sql"
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	responsePkg "github.com/Adhiana46/go-restapi-template/pkg/response"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func handleError(c *fiber.Ctx, err error) error {
	var resp responsePkg.JsonResponse

	if err == sql.ErrNoRows {
		resp = responsePkg.JsonError(404, "", nil)
	} else if strings.Contains(strings.ToLower(err.Error()), "query parameter") {
		resp = responsePkg.JsonError(400, "", nil)
	} else {
		switch err.(type) {
		case validator.ValidationErrors:
			errs := err.(validator.ValidationErrors)
			validationErrors := parseValidationErrors(errs, &validateTrans)
			resp = responsePkg.JsonError(400, "", validationErrors)
		default:
			log.Println("handleError: ", err)
			resp = responsePkg.JsonError(500, "", nil)
		}
	}

	return c.JSON(resp)
}

func parseValidationErrors(validationErrs validator.ValidationErrors, trans *ut.Translator) map[string][]string {
	errorFields := map[string][]string{}
	for _, e := range validationErrs {
		if trans != nil {
			errorFields[e.Field()] = append(errorFields[e.Field()], e.Translate(*trans))
		} else {
			errorFields[e.Field()] = append(errorFields[e.Field()], e.Tag())
		}
	}

	return errorFields
}

func handlePanic(c *fiber.Ctx) {
	if r := recover(); r != nil {
		// TODO: log
		log.Println("Recovered in f", r, string(debug.Stack()))

		response := responsePkg.JsonError(http.StatusInternalServerError, "", nil)

		c.JSON(response)
	}
}

func shouldBind(c *fiber.Ctx, req interface{}) error {
	if err := c.ParamsParser(req); err != nil {
		return err
	}
	if err := c.QueryParser(req); err != nil {
		return err
	}
	if len(c.Body()) > 0 {
		if err := c.BodyParser(req); err != nil {
			return err
		}
	}

	return nil
}
