package interfaces

import "github.com/fin-assistant/internal/services/api/resources"

type Balances interface {
	New() Balances

	Create(balance Balance) error
	Get() (*Balance, error)
	GetById(id int) (*Balance, error)
	GetAllBalances(id int) (*[]Balance, error)

	Transaction(fn func(q Balances) error) (err error)
}

type Balance struct {
	ID       int    `db:"id" structs:"-"`
	Currency string `db:"currency" structs:"currency"`
	UserId   int    `db:"user_id" structs:"user_id"`
}

func (r *Balance) Resource() *resources.GetBalance {
	return &resources.GetBalance{
		Attributes: resources.GetBalanceAttributes{
			Currency:  r.Currency,
			BalanceId: r.ID,
		},
	}
}