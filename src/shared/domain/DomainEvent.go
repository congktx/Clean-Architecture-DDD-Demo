package domain

type DomainEvent interface {
	EventName() string
	OccurredOn() int64
}
