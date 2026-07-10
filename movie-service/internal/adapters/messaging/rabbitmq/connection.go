package rabbitmq

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewConnection(url string) (*amqp.Connection, error) {
	var lastErr error

	for attempt := 1; attempt <= 10; attempt++ {
		connection, err := amqp.Dial(url)
		if err == nil {
			return connection, nil
		}

		lastErr = err
		time.Sleep(3 * time.Second)
	}

	return nil, lastErr
}

func NewChannel(connection *amqp.Connection) (*amqp.Channel, error) {
	return connection.Channel()
}
