package interfaces

import "time"

type Works interface {
	New() Works

	Create(User Work) error
	Get() (*Work, error)
	GetById(id int) (*Work, error)

	Transaction(fn func(q Works) error) (err error)
}

type Work struct {
	ID         int           `db:"id" structs:"-"`
	DateStart  time.Duration `db:"date_start" structs:"date_start"`
	DateFinish time.Duration `db:"date_finish" structs:"date_finish"`
	Position   string        `db:"position" structs:"position"`
	City       string        `db:"city" structs:"city"`
	Company    string        `db:"company" structs:"company"`
	UserId     int           `db:"user_id" structs:"-"`
}
