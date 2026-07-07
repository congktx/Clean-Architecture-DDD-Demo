package domain

type EventDispatcher interface {
	Dispatch(events []DomainEvent) error
}
