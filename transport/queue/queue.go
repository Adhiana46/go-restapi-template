package queue

import amqp "github.com/rabbitmq/amqp091-go"

type QueueWorker interface {
	GetWorkerName() string
	Listen() error
	handlePayload(d amqp.Delivery, payload queuePayload)
	successResponse()
	errorResponse()
	// TODO:
}

type queuePayload struct {
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data"`
}
