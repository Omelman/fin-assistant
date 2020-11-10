package implementation

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/fin-assistant/internal/postgres/interfaces"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const balanceTableName = "balance"

var balanceSelect = sq.Select("balance.*").From(balanceTableName)
var balanceUpdate = sq.Update(balanceTableName)

type Balances struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
	upd sq.UpdateBuilder
}

func NewBalance(db *pgdb.DB) interfaces.Balances {
	return &Balances{
		db:  db.Clone(),
		sql: balanceSelect,
		upd: balanceUpdate,
	}
}

func (q *Balances) New() interfaces.Balances {
	return NewBalance(q.db)
}

func (q *Balances) Create(balance interfaces.Balance) error {
	clauses := structs.Map(balance)

	var id int64
	stmt := sq.Insert(balanceTableName).SetMap(clauses).Suffix("returning id")
	err := q.db.Get(&id, stmt)
	return err
}

func (q *Balances) Transaction(fn func(q interfaces.Balances) error) (err error) {
	return q.db.Transaction(func() error {
		return fn(q)
	})
}

func (q *Balances) Get() (*interfaces.Balance, error) {
	var user interfaces.Balance
	err := q.db.Get(&user, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, err
}

func (q *Balances) GetById(id int) (*interfaces.Balance, error) {
	var user interfaces.Balance
	stmt := sq.Select("*").From(balanceTableName).Where("id = ?", id)
	err := q.db.Get(&user, stmt)
	return &user, err
}

func (q *Balances) GetAllBalances(id int) (*[]interfaces.Balance, error) {
	var user []interfaces.Balance
	stmt := sq.Select("*").From(balanceTableName).Where("user_id = ?", id)
	err := q.db.Select(&user, stmt)
	return &user, err
}
