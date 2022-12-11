package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Adhiana46/go-restapi-template/internal/repository"
	"github.com/Adhiana46/go-restapi-template/internal/service"
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

type Config struct {
	Host string `env:"HOST" env-default:""`
	Port string `env:"PORT" env-default:"8000"`

	DbHost string `env:"DB_HOST" env-default:"localhost"`
	DbPort string `env:"DB_PORT" env-default:"5432"`
	DbUser string `env:"DB_USER" env-default:"user"`
	DbPass string `env:"DB_PASS" env-default:"secret"`
	DbName string `env:"DB_NAME" env-default:"todoapp"`
	DbSSL  string `env:"DB_SSL" env-default:"disable"`
}

var cfg Config

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
	log.Infoln("load environment variables")
	var err error
	if _, err := os.Stat(".env"); err == nil {
		err = cleanenv.ReadConfig(".env", &cfg)
	} else {
		err = cleanenv.ReadEnv(&cfg)
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
	db, err = openDB()
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

func openDB() (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		cfg.DbHost,
		cfg.DbPort,
		cfg.DbUser,
		cfg.DbName,
		cfg.DbSSL,
		cfg.DbPass,
	)

	dbConn, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, err
	}

	dbConn.SetMaxOpenConns(60)
	dbConn.SetConnMaxLifetime(120 * time.Second)
	dbConn.SetMaxIdleConns(30)
	dbConn.SetConnMaxIdleTime(20 * time.Second)
	if err = dbConn.Ping(); err != nil {
		return nil, err
	} else {
		log.Infoln("Pinged database successfully!")
	}

	return dbConn, nil
}
