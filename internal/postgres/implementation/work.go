package implementation

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/fin-assistant/internal/postgres/interfaces"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const workTableName = "work"

var workSelect = sq.Select("work.*").From(workTableName)
var workUpdate = sq.Update(workTableName)

type Works struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
	upd sq.UpdateBuilder
}

func NewWork(db *pgdb.DB) interfaces.Works {
	return &Works{
		db:  db.Clone(),
		sql: workSelect,
		upd: workUpdate,
	}
}

func (q *Works) New() interfaces.Works {
	return NewWork(q.db)
}

func (q *Works) Create(user interfaces.Work) error {
	clauses := structs.Map(user)

	var id int64
	stmt := sq.Insert(workTableName).SetMap(clauses).Suffix("returning id")
	err := q.db.Get(&id, stmt)
	return err
}

func (q *Works) Transaction(fn func(q interfaces.Works) error) (err error) {
	return q.db.Transaction(func() error {
		return fn(q)
	})
}

func (q *Works) Get() (*interfaces.Work, error) {
	var user interfaces.Work
	err := q.db.Get(&user, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, err
}

func (q *Works) GetById(id int) (*interfaces.Work, error) {
	var user interfaces.Work
	stmt := sq.Select("*").From(workTableName).Where("id = ?", id)
	err := q.db.Get(&user, stmt)
	return &user, err
}
