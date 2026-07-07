package domain

type BudgetRepository interface {
	Save(budget *BudgetEntity) error
	FindByID(id string) (*BudgetEntity, error)
	FindByOwnerID(ownerId string) ([]*BudgetEntity, error)
}
