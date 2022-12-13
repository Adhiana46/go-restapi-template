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

type activityGroupWorker struct {
	conn      *amqp.Connection
	queueName string

	svcActivityGroup service.ActivityGroupService
}

func NewActivityGroupWorker(conn *amqp.Connection, queueName string, svcActivityGroup service.ActivityGroupService) QueueWorker {
	return &activityGroupWorker{
		conn:             conn,
		queueName:        queueName,
		svcActivityGroup: svcActivityGroup,
	}
}

func (w *activityGroupWorker) GetWorkerName() string {
	return w.queueName
}

func (w *activityGroupWorker) Listen() error {
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

func (w *activityGroupWorker) handlePayload(d amqp.Delivery, payload queueRequestPayload) {
	dataJson, err := json.Marshal(payload.Data)
	if err != nil {
		errorResponse(w.conn, w.queueName, payload.Action, err)
		d.Ack(false)
		return
	}

	log.Infof("[%s] Receiving message: %s -> %s", w.queueName, payload.Action, string(dataJson))

	switch payload.Action {
	case "create":
		activityGroup, err := w.handleCreate(dataJson)
		if err != nil {
			errorResponse(w.conn, w.queueName, payload.Action, err)
		} else {
			successResponse(w.conn, w.queueName, "created", activityGroup)
		}
	case "update":
		activityGroup, err := w.handleUpdate(dataJson)
		if err != nil {
			errorResponse(w.conn, w.queueName, payload.Action, err)
		} else {
			successResponse(w.conn, w.queueName, "updated", activityGroup)
		}
	case "delete":
		activityGroup, err := w.handleDelete(dataJson)
		if err != nil {
			errorResponse(w.conn, w.queueName, payload.Action, err)
		} else {
			successResponse(w.conn, w.queueName, "deleted", activityGroup)
		}
	}

	d.Ack(false)
}

func (w *activityGroupWorker) handleCreate(data []byte) (*entity.ActivityGroup, error) {
	reqDto := dto.ActivityGroupCreateRequest{}

	err := json.Unmarshal(data, &reqDto)
	if err != nil {
		return nil, err
	}

	activityGroup, err := w.svcActivityGroup.Create(reqDto)
	if err != nil {
		return nil, err
	}

	return activityGroup, nil
}

func (w *activityGroupWorker) handleUpdate(data []byte) (*entity.ActivityGroup, error) {
	reqDto := dto.ActivityGroupUpdateRequest{}

	err := json.Unmarshal(data, &reqDto)
	if err != nil {
		return nil, err
	}

	activityGroup, err := w.svcActivityGroup.Update(reqDto)
	if err != nil {
		return nil, err
	}

	return activityGroup, nil
}

func (w *activityGroupWorker) handleDelete(data []byte) (*entity.ActivityGroup, error) {
	reqDto := dto.ActivityGroupUuidRequest{}

	err := json.Unmarshal(data, &reqDto)
	if err != nil {
		return nil, err
	}

	activityGroup, err := w.svcActivityGroup.FindByUuid(reqDto)
	if err != nil {
		return nil, err
	}

	err = w.svcActivityGroup.Delete(reqDto)
	if err != nil {
		return nil, err
	}

	return activityGroup, nil
}
