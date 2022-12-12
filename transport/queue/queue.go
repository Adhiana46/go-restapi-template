package queue

import amqp "github.com/rabbitmq/amqp091-go"

type QueueWorker interface {
	GetWorkerName() string
	Listen() error
	handlePayload(d amqp.Delivery, payload queueRequestPayload)
	successResponse(action string, data interface{}) error
	errorResponse(action string, err error) error
	// TODO:
}

type queueRequestPayload struct {
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data"`
}
