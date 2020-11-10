package implementation

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/fin-assistant/internal/postgres/interfaces"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const goalTableName = "goal"

var goalSelect = sq.Select("goal.*").From(goalTableName)
var goalUpdate = sq.Update(goalTableName)

type Goals struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
	upd sq.UpdateBuilder
}

func NewGoal(db *pgdb.DB) interfaces.Goals {
	return &Goals{
		db:  db.Clone(),
		sql: goalSelect,
		upd: goalUpdate,
	}
}

func (q *Goals) New() interfaces.Goals {
	return NewGoal(q.db)
}

func (q *Goals) Create(user interfaces.Goal) error {
	clauses := structs.Map(user)

	var id int64
	stmt := sq.Insert(goalTableName).SetMap(clauses).Suffix("returning id")
	err := q.db.Get(&id, stmt)
	return err
}

func (q *Goals) Transaction(fn func(q interfaces.Goals) error) (err error) {
	return q.db.Transaction(func() error {
		return fn(q)
	})
}

func (q *Goals) Get() (*interfaces.Goal, error) {
	var user interfaces.Goal
	err := q.db.Get(&user, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, err
}

func (q *Goals) GetById(id int) (*interfaces.Goal, error) {
	var user interfaces.Goal
	stmt := sq.Select("*").From(goalTableName).Where("id = ?", id)
	err := q.db.Get(&user, stmt)
	return &user, err
}
