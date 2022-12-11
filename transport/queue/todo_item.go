package queue

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type todoItemWorker struct {
	conn      *amqp.Connection
	queueName string
}

func NewTodoItemWorker(conn *amqp.Connection, queueName string) QueueWorker {
	return &todoItemWorker{
		conn:      conn,
		queueName: queueName,
	}
}

func (w *todoItemWorker) GetWorkerName() string {
	return w.queueName
}

func (w *todoItemWorker) Listen() error {
	ch, err := w.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		fmt.Sprintf("%s.request", w.queueName), // name
		true,                                   // durable
		false,                                  // delete when unused
		false,                                  // exclusive
		false,                                  // no-wait
		nil,                                    // arguments
	)
	if err != nil {
		return err
	}

	// set Qos
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	var forever chan struct{}
	go func() {
		for d := range msgs {
			var payload queuePayload
			_ = json.Unmarshal(d.Body, &payload)

			go w.handlePayload(d, payload)
			d.Ack(false)
		}
	}()

	<-forever

	return nil
}

func (w *todoItemWorker) handlePayload(d amqp.Delivery, payload queuePayload) {
	// TODO: do something

	d.Ack(false)
}

func (w *todoItemWorker) successResponse() {
	//
}

func (w *todoItemWorker) errorResponse() {
	//
}
