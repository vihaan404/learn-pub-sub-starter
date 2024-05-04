package pubsub

import (
	"context"
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp091.Channel, exchange, key string, val T) error {
	res, err := json.Marshal(val)
	if err != nil {
		return err

	}
	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        res,
	})
	if err != nil {
		return err
	}
	return nil

}

//exchange: peril_direct (this is a constant in the internal/routing package)
//queueName: pause.username where username is the user's input. The pause section of the name is the routing key constant in the internal/routing package.
//routingKey: pause (this is a constant in the internal/routing package)
//simpleQueueType: transient

func DeclareAndBind(
	conn *amqp091.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int, // an enum to represent "durable" or "transient"
) (*amqp091.Channel, amqp091.Queue, error) {
	ch, _ := conn.Channel()
	var isDurable bool
	var toAutoDelete bool
	var isExclusive bool
	if simpleQueueType == 2 {
		isDurable = true
		isExclusive = false
		toAutoDelete = false
	} else {
		isDurable = false
		isExclusive = true
		toAutoDelete = true
	}

	q, _ := ch.QueueDeclare(queueName, isDurable, toAutoDelete, isExclusive, false, nil)
	_ = ch.QueueBind("", key, exchange, false, nil)

	return ch, q, nil
}
