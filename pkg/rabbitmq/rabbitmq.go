package rabbitmq

import (
	"fmt"
	"math"
	"time"

	"github.com/Adhiana46/go-restapi-template/config"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

func OpenConn(cfg *config.Config) (*amqp.Connection, error) {
	var count int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	dsn := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.AmqpUser,
		cfg.AmqpPass,
		cfg.AmqpHost,
		cfg.AmqpPort,
	)

	// Don't continue until rabbit is ready
	for {
		c, err := amqp.Dial(dsn)
		if err != nil {
			log.Errorln("RabbitMQ not yet ready...", err)
			count++
		} else {
			log.Infoln("Connected to RabbitMQ...")
			connection = c
			break
		}

		if count > 5 {
			log.Errorln("Could not connect to RabbitMQ", err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(count), 2)) * time.Second
		log.Infoln("backing off", backOff)
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
