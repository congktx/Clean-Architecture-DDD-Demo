package domain

type WalletRepository interface {
	FindByID(id string) (*WalletEntity, error)
	FindByOwnerID(ownerId string) ([]*WalletEntity, error)
	Save(wallet *WalletEntity) error
}
