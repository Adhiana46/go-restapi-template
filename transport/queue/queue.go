package queue

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueWorker interface {
	GetWorkerName() string
	Listen() error
	handlePayload(d amqp.Delivery, payload queueRequestPayload)
}

type queueRequestPayload struct {
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data"`
}

func successResponse(conn *amqp.Connection, queueName string, action string, data interface{}) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		fmt.Sprintf("%s.%s", queueName, action), // name
		false,                                   // durable
		false,                                   // delete when unused
		false,                                   // exclusive
		false,                                   // no-wait
		nil,                                     // arguments
	)

	dataJson, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        dataJson,
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func errorResponse(conn *amqp.Connection, queueName string, action string, errAct error) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		fmt.Sprintf("%s.error", queueName), // name
		false,                              // durable
		false,                              // delete when unused
		false,                              // exclusive
		false,                              // no-wait
		nil,                                // arguments
	)

	data := map[string]interface{}{
		"action": action,
		"error":  errAct.Error(),
	}
	dataJson, _ := json.Marshal(data)

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        dataJson,
		},
	)

	if err != nil {
		return err
	}

	return nil
}
