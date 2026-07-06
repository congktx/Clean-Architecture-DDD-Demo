package domain

type AggregateRoot struct {
	domainEvents []DomainEvent
}

func (a *AggregateRoot) AddDomainEvent(event DomainEvent) {
	a.domainEvents = append(a.domainEvents, event)
}

func (a *AggregateRoot) GetDomainEvents() []DomainEvent {
	return a.domainEvents
}

func (a *AggregateRoot) ClearDomainEvents() {
	a.domainEvents = []DomainEvent{}
}
