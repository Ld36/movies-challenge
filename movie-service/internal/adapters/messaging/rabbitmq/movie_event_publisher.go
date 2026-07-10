package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/luizdavid/movies-challenge/movie-service/internal/domain"
	"github.com/luizdavid/movies-challenge/movie-service/internal/events"
	amqp "github.com/rabbitmq/amqp091-go"
)

const movieEventsExchange = "movies.events"

type MovieEventPublisher struct {
	channel *amqp.Channel
}

func NewMovieEventPublisher(channel *amqp.Channel) *MovieEventPublisher {
	return &MovieEventPublisher{
		channel: channel,
	}
}

func (p *MovieEventPublisher) PublishMovieCreated(ctx context.Context, movie domain.Movie) error {
	payload := events.MovieCreatedPayload{
		ID:    movie.ID,
		Title: movie.Title,
		Year:  movie.Year,
	}

	return p.publish(ctx, events.MovieCreatedEvent, payload)
}

func (p *MovieEventPublisher) PublishMovieDeleted(ctx context.Context, id int64) error {
	payload := events.MovieDeletedPayload{
		ID: id,
	}

	return p.publish(ctx, events.MovieDeletedEvent, payload)
}

func (p *MovieEventPublisher) publish(ctx context.Context, routingKey string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return p.channel.PublishWithContext(
		ctx,
		movieEventsExchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Body:         body,
		},
	)
}

func DeclareMovieEventsExchange(channel *amqp.Channel) error {
	return channel.ExchangeDeclare(
		movieEventsExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
}
