package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Adhiana46/go-restapi-template/config"
	"github.com/Adhiana46/go-restapi-template/internal/repository"
	"github.com/Adhiana46/go-restapi-template/internal/service"
	"github.com/Adhiana46/go-restapi-template/pkg/sqldb"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	id_translations "github.com/go-playground/validator/v10/translations/id"
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

var (
	db *sqlx.DB

	// utils
	validate      *validator.Validate
	validateTrans ut.Translator

	// Repository
	repoActivityGroup repository.ActivityGroupRepository
	repoTodoItem      repository.TodoItemRepository

	// Services
	svcActivityGroup service.ActivityGroupService
	svcTodoItem      service.TodoItemService
)

var cfg *config.Config

func main() {
	boot()
	defer db.Close()

	r := httpRoutes()

	if err := r.Listen(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)); err != nil {
		log.Panicf("Can't start the server, error: %s", err)
	}
}

func boot() {
	log.SetReportCaller(true)
	log.SetFormatter(&easy.Formatter{
		TimestampFormat: time.RFC3339,
		LogFormat:       "[%lvl%][%time%]: %msg%\n",
	})
	log.SetOutput(os.Stdout)

	// Load environment variables
	cfg = &config.Config{}
	log.Infoln("load environment variables")
	var err error
	if _, err := os.Stat(".env"); err == nil {
		err = cleanenv.ReadConfig(".env", cfg)
	} else {
		err = cleanenv.ReadEnv(cfg)
	}

	if err != nil {
		log.Panicf("Can't read environment variable: %s", err)
	}

	// validation & validation trans
	log.Infoln("setup validation")
	id := id.New()
	uni := ut.New(id, id)
	validateTrans, _ = uni.GetTranslator("id")
	validate = validator.New()
	id_translations.RegisterDefaultTranslations(validate, validateTrans)

	log.Infoln("Connecting to database...")
	db, err = sqldb.OpenConn(cfg)
	if err != nil {
		log.Panicf("Can't open database connection: %s", err)
	}

	// repositories
	repoActivityGroup = repository.NewPostgresActivityGroupRepository(db)
	repoTodoItem = repository.NewPostgresTodoItemRepository(db)

	// services
	svcActivityGroup = service.NewActivityGroupService(validate, repoActivityGroup)
	svcTodoItem = service.NewTodoItemService(validate, repoTodoItem, repoActivityGroup)
}
