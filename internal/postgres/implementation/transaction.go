package implementation

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/fin-assistant/internal/postgres/interfaces"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const transactionTableName = "transaction"

var transactionSelect = sq.Select(" transaction.*").From(transactionTableName)
var transactionUpdate = sq.Update(transactionTableName)

type Transactions struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
	upd sq.UpdateBuilder
}

func NewTransaction(db *pgdb.DB) interfaces.Transactions {
	return &Transactions{
		db:  db.Clone(),
		sql: transactionSelect,
		upd: transactionUpdate,
	}
}

func (q *Transactions) New() interfaces.Transactions {
	return NewTransaction(q.db)
}

func (q *Transactions) Create(user interfaces.Transaction) error {
	clauses := structs.Map(user)

	var id int64
	stmt := sq.Insert(transactionTableName).SetMap(clauses).Suffix("returning id")
	err := q.db.Get(&id, stmt)
	return err
}

func (q *Transactions) Transaction(fn func(q interfaces.Transactions) error) (err error) {
	return q.db.Transaction(func() error {
		return fn(q)
	})
}

func (q *Transactions) Get() (*interfaces.Transaction, error) {
	var user interfaces.Transaction
	err := q.db.Get(&user, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, err
}

func (q *Transactions) GetById(id int) (*interfaces.Transaction, error) {
	var user interfaces.Transaction
	stmt := sq.Select("*").From(transactionTableName).Where("id = ?", id)
	err := q.db.Get(&user, stmt)
	return &user, err
}
