package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Adhiana46/go-restapi-template/internal/repository"
	"github.com/Adhiana46/go-restapi-template/internal/service"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	id_translations "github.com/go-playground/validator/v10/translations/id"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

var (
	// utils
	validate      *validator.Validate
	validateTrans ut.Translator

	// Repository
	repoActivityGroup repository.ActivityGroupRepository

	// Services
	svcActivityGroup service.ActivityGroupService
)

func main() {
	boot()

	r := routes()

	if err := r.Listen(":8000"); err != nil {
		log.Panicf("Can't start the server, error: %s", err)
	}
}

func boot() {
	// validation & validation trans
	id := id.New()
	uni := ut.New(id, id)
	validateTrans, _ = uni.GetTranslator("id")
	validate = validator.New()
	id_translations.RegisterDefaultTranslations(validate, validateTrans)

	db, err := openDB()
	if err != nil {
		log.Panicf("Can't open database connection: %s", err)
	}

	// repositories
	repoActivityGroup = repository.NewPostgresActivityGroupRepository(db)

	// services
	svcActivityGroup = service.NewActivityGroupService(validate, repoActivityGroup)
}

func openDB() (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		"localhost",
		"5432",
		"user",
		"todoapp",
		"disable",
		"secret",
	)

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(60)
	db.SetConnMaxLifetime(120 * time.Second)
	db.SetMaxIdleConns(30)
	db.SetConnMaxIdleTime(20 * time.Second)
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
