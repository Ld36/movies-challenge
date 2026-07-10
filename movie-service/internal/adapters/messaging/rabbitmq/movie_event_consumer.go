package rabbitmq

import (
	"context"
	"encoding/json"

	"github.com/luizdavid/movies-challenge/movie-service/internal/domain"
	"github.com/luizdavid/movies-challenge/movie-service/internal/events"
	"github.com/luizdavid/movies-challenge/movie-service/internal/ports"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

const movieEventsQueue = "movies.events.worker"

type MovieEventConsumer struct {
	channel    *amqp.Channel
	repository ports.MovieRepository
	logger     *zap.Logger
}

func NewMovieEventConsumer(channel *amqp.Channel, repository ports.MovieRepository, logger *zap.Logger) *MovieEventConsumer {
	return &MovieEventConsumer{
		channel:    channel,
		repository: repository,
		logger:     logger,
	}
}

func (c *MovieEventConsumer) Start(ctx context.Context) error {
	queue, err := c.channel.QueueDeclare(
		movieEventsQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	bindings := []string{
		events.MovieCreatedEvent,
		events.MovieDeletedEvent,
	}

	for _, routingKey := range bindings {
		if err := c.channel.QueueBind(
			queue.Name,
			routingKey,
			movieEventsExchange,
			false,
			nil,
		); err != nil {
			return err
		}
	}

	messages, err := c.channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	c.logger.Info("movie event consumer started", zap.String("queue", queue.Name))

	go func() {
		for {
			select {
			case <-ctx.Done():
				c.logger.Info("movie event consumer stopped")
				return

			case message, ok := <-messages:
				if !ok {
					c.logger.Warn("movie event messages channel closed")
					return
				}

				if err := c.handleMessage(ctx, message); err != nil {
					c.logger.Error("failed to handle movie event",
						zap.String("routing_key", message.RoutingKey),
						zap.Error(err),
					)

					_ = message.Nack(false, true)
					continue
				}

				_ = message.Ack(false)
			}
		}
	}()

	return nil
}

func (c *MovieEventConsumer) handleMessage(ctx context.Context, message amqp.Delivery) error {
	switch message.RoutingKey {
	case events.MovieCreatedEvent:
		return c.handleMovieCreated(ctx, message.Body)

	case events.MovieDeletedEvent:
		return c.handleMovieDeleted(ctx, message.Body)

	default:
		c.logger.Warn("unknown movie event", zap.String("routing_key", message.RoutingKey))
		return nil
	}
}

func (c *MovieEventConsumer) handleMovieCreated(ctx context.Context, body []byte) error {
	var payload events.MovieCreatedPayload

	if err := json.Unmarshal(body, &payload); err != nil {
		return err
	}

	movie := domain.Movie{
		ID:    payload.ID,
		Title: payload.Title,
		Year:  payload.Year,
	}

	_, err := c.repository.Create(ctx, movie)
	if err != nil {
		return err
	}

	c.logger.Info("movie created event processed", zap.Int64("movie_id", movie.ID))

	return nil
}

func (c *MovieEventConsumer) handleMovieDeleted(ctx context.Context, body []byte) error {
	var payload events.MovieDeletedPayload

	if err := json.Unmarshal(body, &payload); err != nil {
		return err
	}

	if err := c.repository.Delete(ctx, payload.ID); err != nil {
		return err
	}

	c.logger.Info("movie deleted event processed", zap.Int64("movie_id", payload.ID))

	return nil
}
