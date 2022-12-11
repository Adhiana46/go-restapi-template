package sqldb

import (
	"fmt"
	"time"

	"github.com/Adhiana46/go-restapi-template/config"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func OpenConn(cfg *config.Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		cfg.DbHost,
		cfg.DbPort,
		cfg.DbUser,
		cfg.DbName,
		cfg.DbSSL,
		cfg.DbPass,
	)

	dbConn, err := sqlx.Connect(cfg.DbDialect, dsn)
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
