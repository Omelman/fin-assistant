package interfaces

import "time"

type Goals interface {
	New() Goals

	Create(User Goal) error
	Get() (*Goal, error)
	GetById(id int) (*Goal, error)

	Transaction(fn func(q Goals) error) (err error)
}

type Goal struct {
	ID          int           `db:"id" structs:"-"`
	DateStart   time.Duration `db:"date_start" structs:"date_start"`
	DateFinish  time.Duration `db:"date_finish" structs:"date_finish"`
	Description string        `db:"description" structs:"description"`
	Amount      int           `db:"amount" structs:"amount"`
	BalanceId   int           `db:"balance_id" structs:"-"`
}
