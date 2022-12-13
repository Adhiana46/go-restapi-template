package queue

import (
	"encoding/json"
	"fmt"

	"github.com/Adhiana46/go-restapi-template/internal/dto"
	"github.com/Adhiana46/go-restapi-template/internal/entity"
	"github.com/Adhiana46/go-restapi-template/internal/service"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type todoItemWorker struct {
	conn      *amqp.Connection
	queueName string

	svcTodoItem service.TodoItemService
}

func NewTodoItemWorker(conn *amqp.Connection, queueName string, svcTodoItem service.TodoItemService) QueueWorker {
	return &todoItemWorker{
		conn:      conn,
		queueName: queueName,

		svcTodoItem: svcTodoItem,
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
			var payload queueRequestPayload
			_ = json.Unmarshal(d.Body, &payload)

			go w.handlePayload(d, payload)
		}
	}()

	<-forever

	return nil
}

func (w *todoItemWorker) handlePayload(d amqp.Delivery, payload queueRequestPayload) {
	dataJson, err := json.Marshal(payload.Data)
	if err != nil {
		errorResponse(w.conn, w.queueName, payload.Action, err)
		d.Ack(false)
		return
	}

	log.Infof("[%s] Receiving message: %s -> %s", w.queueName, payload.Action, string(dataJson))

	switch payload.Action {
	case "create":
		todoItem, err := w.handleCreate(dataJson)
		if err != nil {
			errorResponse(w.conn, w.queueName, payload.Action, err)
		} else {
			successResponse(w.conn, w.queueName, "created", todoItem)
		}
	case "update":
		todoItem, err := w.handleUpdate(dataJson)
		if err != nil {
			errorResponse(w.conn, w.queueName, payload.Action, err)
		} else {
			successResponse(w.conn, w.queueName, "updated", todoItem)
		}
	case "delete":
		todoItem, err := w.handleDelete(dataJson)
		if err != nil {
			errorResponse(w.conn, w.queueName, payload.Action, err)
		} else {
			successResponse(w.conn, w.queueName, "deleted", todoItem)
		}
	}

	d.Ack(false)
}

func (w *todoItemWorker) handleCreate(data []byte) (*entity.TodoItem, error) {
	reqDto := dto.TodoItemCreateRequest{}

	err := json.Unmarshal(data, &reqDto)
	if err != nil {
		return nil, err
	}

	todoItem, err := w.svcTodoItem.Create(reqDto)
	if err != nil {
		return nil, err
	}

	return todoItem, nil
}

func (w *todoItemWorker) handleUpdate(data []byte) (*entity.TodoItem, error) {
	reqDto := dto.TodoItemUpdateRequest{}

	err := json.Unmarshal(data, &reqDto)
	if err != nil {
		return nil, err
	}

	todoItem, err := w.svcTodoItem.Update(reqDto)
	if err != nil {
		return nil, err
	}

	return todoItem, nil
}

func (w *todoItemWorker) handleDelete(data []byte) (*entity.TodoItem, error) {
	reqDto := dto.TodoItemUuidRequest{}

	err := json.Unmarshal(data, &reqDto)
	if err != nil {
		return nil, err
	}

	todoItem, err := w.svcTodoItem.FindByUuid(reqDto)
	if err != nil {
		return nil, err
	}

	err = w.svcTodoItem.Delete(reqDto)
	if err != nil {
		return nil, err
	}

	return todoItem, nil
}
