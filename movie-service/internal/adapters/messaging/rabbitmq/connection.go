package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func NewConnection(url string) (*amqp.Connection, error) {
	return amqp.Dial(url)
}

func NewChannel(connection *amqp.Connection) (*amqp.Channel, error) {
	return connection.Channel()
}
