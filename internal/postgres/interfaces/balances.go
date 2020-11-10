package interfaces

type Balances interface {
	New() Balances

	Create(User Balance) error
	Get() (*Balance, error)
	GetById(id int) (*Balance, error)

	Transaction(fn func(q Balances) error) (err error)
}

type Balance struct {
	ID       int    `db:"id" structs:"-"`
	Currency string `db:"currency" structs:"currency"`
	UserId   int    `db:"user_id" structs:"-"`
}
