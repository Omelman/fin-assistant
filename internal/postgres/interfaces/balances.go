package interfaces

import "github.com/fin-assistant/internal/services/api/resources"

type Balances interface {
	New() Balances

	Create(balance Balance) (int, error)
	Get() (*Balance, error)
	GetById(id int) (*Balance, error)
	GetAllBalances(id int) (*[]Balance, error)
	DeleteBalance(userId int, balanceId int) error
	GetRestOnBalance(userId int, balanceId int) (*int, error)

	Transaction(fn func(q Balances) error) (err error)
}

type Balance struct {
	ID       int    `db:"id" structs:"-"`
	Currency string `db:"currency" structs:"currency"`
	UserId   int    `db:"user_id" structs:"user_id"`
}

func (r *Balance) Resource(amount int) *resources.GetBalance {
	return &resources.GetBalance{
		Attributes: resources.GetBalanceAttributes{
			Currency:  r.Currency,
			BalanceId: r.ID,
			Amount:    amount,
		},
	}
}
