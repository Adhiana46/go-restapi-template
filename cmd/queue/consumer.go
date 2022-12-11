package main

import (
	"github.com/Adhiana46/go-restapi-template/transport/queue"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type consumer struct {
	workers []queue.QueueWorker
}

func newConsumer(conn *amqp.Connection) (*consumer, error) {
	c := &consumer{
		workers: []queue.QueueWorker{
			queue.NewActivityGroupWorker(conn, "activity-group"),
			queue.NewTodoItemWorker(conn, "todo-item"),
		},
	}

	return c, nil
}

func (c *consumer) listen() {
	var forever chan struct{}
	// log.Printf("[*] Waiting for messages (%s).", fmt.Sprintf("%s.request", w.queueName))

	log.Infoln("Registering queue workers:")
	for _, worker := range c.workers {
		go worker.Listen()
		log.Infoln(" *", worker.GetWorkerName())
	}

	log.Infof("Waiting for messages. To exit press CTRL+C")
	<-forever
}
