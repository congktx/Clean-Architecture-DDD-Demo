package domain

// EventDispatcher defines how domain events are dispatched to their handlers.
type EventDispatcher interface {
	Dispatch(events []DomainEvent) error
}
