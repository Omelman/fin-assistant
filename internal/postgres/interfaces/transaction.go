package interfaces

import "time"

type Transactions interface {
	New() Transactions

	Create(User Transaction) error
	Get() (*Transaction, error)
	GetById(id int) (*Transaction, error)

	Transaction(fn func(q Transactions) error) (err error)
}

type Transaction struct {
	ID          int       `db:"id" structs:"-"`
	Date        time.Time `db:"date" structs:"date"`
	Description string    `db:"description" structs:"description"`
	Amount      int       `db:"amount" structs:"amount"`
	Category    string    `db:"category" structs:"category"`
	Include     bool      `db:"include" structs:"include"`
	BalanceId   int       `db:"balance_id" structs:"-"`
}
