package events

const (
	MovieCreatedEvent = "movie.created"
	MovieDeletedEvent = "movie.deleted"
)

type MovieCreatedPayload struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Year  string `json:"year"`
}

type MovieDeletedPayload struct {
	ID int64 `json:"id"`
}
