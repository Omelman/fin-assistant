package interfaces

import (
	"github.com/fin-assistant/internal/services/api/resources"
	"time"
)

type Transactions interface {
	New() Transactions

	Create(User Transaction) (int, error)
	Get() (*Transaction, error)
	GetById(id int) (*Transaction, error)
	GetAllTransaction(userId int) (*[]Transaction, error)

	Transaction(fn func(q Transactions) error) (err error)
}

type Transaction struct {
	ID          int       `db:"id" structs:"-"`
	Date        time.Time `db:"date" structs:"date"`
	Description string    `db:"description" structs:"description"`
	Amount      int       `db:"amount" structs:"amount"`
	Category    string    `db:"category" structs:"category"`
	Include     bool      `db:"include" structs:"include"`
	BalanceId   int       `db:"balance_id" structs:"balance_id"`
}

func (r *Transaction) Resource() *resources.CreateTransaction {
	return &resources.CreateTransaction{
		Attributes: resources.CreateTransactionAttributes{
			Date:        r.Date.Format("2006-01-02"),
			Description: r.Description,
			Amount:      r.Amount,
			Category:    r.Category,
			Include:     r.Include,
			BalaceId:    r.BalanceId,
		},
	}
}
