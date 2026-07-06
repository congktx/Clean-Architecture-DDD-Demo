package infrastructure

import (
	"database/sql"
	"fmt"
	"src/src/budget/domain"
	shared "src/src/shared/domain"

	"github.com/shopspring/decimal"
)

type PostgresBudgetRepository struct {
	db *sql.DB
}

func NewPostgresBudgetRepository(db *sql.DB) *PostgresBudgetRepository {
	return &PostgresBudgetRepository{db: db}
}

func (r *PostgresBudgetRepository) Save(budget *domain.BudgetEntity) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO budgets (id, owner_id, month, year, total_limit_amount, total_limit_currency)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			total_limit_amount = EXCLUDED.total_limit_amount,
			total_limit_currency = EXCLUDED.total_limit_currency;
	`
	_, err = tx.Exec(query,
		budget.ID(),
		budget.OwnerID(),
		budget.Period().Month(),
		budget.Period().Year(),
		budget.TotalLimit().Amount().String(),
		budget.TotalLimit().Currency(),
	)
	if err != nil {
		return fmt.Errorf("failed to save budget: %w", err)
	}

	for _, item := range budget.BudgetItems() {
		itemQuery := `
			INSERT INTO budget_items (id, budget_id, category_id, allocated_amount, spent_amount, currency)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (id) DO UPDATE SET
				allocated_amount = EXCLUDED.allocated_amount,
				spent_amount = EXCLUDED.spent_amount;
		`
		_, err = tx.Exec(itemQuery,
			item.ID(),
			budget.ID(),
			item.CategoryID(),
			item.AllocatedAmount().Amount().String(),
			item.SpentAmount().Amount().String(),
			item.AllocatedAmount().Currency(),
		)
		if err != nil {
			return fmt.Errorf("failed to save budget item: %w", err)
		}
	}

	return tx.Commit()
}

func (r *PostgresBudgetRepository) FindByID(id string) (*domain.BudgetEntity, error) {
	query := `SELECT owner_id, month, year, total_limit_amount, total_limit_currency FROM budgets WHERE id = $1`
	var ownerID, limitCurrency string
	var month, year int
	var limitAmountStr string

	err := r.db.QueryRow(query, id).Scan(&ownerID, &month, &year, &limitAmountStr, &limitCurrency)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	period, _ := shared.NewMonthYear(month, year)
	limitAmount, _ := decimal.NewFromString(limitAmountStr)
	totalLimit := shared.NewMoneyObject(limitAmount, limitCurrency)

	itemQuery := `SELECT id, category_id, allocated_amount, spent_amount, currency FROM budget_items WHERE budget_id = $1`
	rows, err := r.db.Query(itemQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.BudgetItemEntity
	for rows.Next() {
		var itemID, catID, allocStr, spentStr, currency string
		if err := rows.Scan(&itemID, &catID, &allocStr, &spentStr, &currency); err != nil {
			return nil, err
		}
		allocAmt, _ := decimal.NewFromString(allocStr)
		spentAmt, _ := decimal.NewFromString(spentStr)
		
		allocMoney := shared.NewMoneyObject(allocAmt, currency)
		spentMoney := shared.NewMoneyObject(spentAmt, currency)
		
		item := domain.LoadBudgetItem(itemID, catID, allocMoney, spentMoney)
		items = append(items, item)
	}

	budget := domain.LoadBudget(id, ownerID, period, totalLimit, items)
	return &budget, nil
}
