package implementation

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/fin-assistant/internal/postgres/interfaces"
	"gitlab.com/distributed_lab/kit/pgdb"
	"strings"
	"time"
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

func (q *Transactions) Create(transaction interfaces.Transaction) (int, error) {
	clauses := structs.Map(transaction)

	var id int
	stmt := sq.Insert(transactionTableName).SetMap(clauses).Suffix("returning id")
	err := q.db.Get(&id, stmt)
	return id, err
}

func (q *Transactions) Select() (*[]interfaces.Transaction, error) {
	var result []interfaces.Transaction
	err := q.db.Select(&result, q.sql)
	return &result, err
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

func (q *Transactions) GetExpenses(userId int, category string) (*int, error) {
	var sum sql.NullInt32
	stmt := sq.Select(fmt.Sprintf("SUM(transaction.amount) FROM users "+
		"INNER JOIN balance ON users.id = balance.user_id INNER JOIN transaction "+
		"ON transaction.balance_id = balance.id WHERE users.id = %d AND transaction.category = '%s'",
		userId, category))
	err := q.db.Get(&sum, stmt)
	ans := int(sum.Int32)
	return &ans, err
}

func (q *Transactions) Update(transaction interfaces.Transaction, transactionId int) error {
	clauses := structs.Map(transaction)

	stmt := sq.Update(transactionTableName).SetMap(clauses).Where("id = ?", transactionId)
	err := q.db.Exec(stmt)
	return err
}

func (q *Transactions) DeleteTransaction(transactionId int) error {
	stmt := sq.Delete(transactionTableName).Where("id = ?", transactionId)
	err := q.db.Exec(stmt)
	return err
}

func (q *Transactions) FilterByCategory(code string) interfaces.Transactions {
	q.sql = q.sql.Where("category = ?", code)
	return q
}

func (q *Transactions) FilterByBalance(code string) interfaces.Transactions {
	q.sql = q.sql.Where("transaction.balance_id = ?", code)
	return q
}

func (q *Transactions) FilterOnlyBefore(time time.Time) interfaces.Transactions {
	q.sql = q.sql.Where(sq.Lt{"transaction.date": time})
	return q
}

func (q *Transactions) FilterOnlyAfter(time time.Time) interfaces.Transactions {
	q.sql = q.sql.Where(sq.Gt{"transaction.date": time})
	return q
}

func (q *Transactions) UserJoined() interfaces.Transactions {
	q.sql = q.sql.Join(fmt.Sprintf("%s on %s.balance_id = %s.id", balanceTableName,
		transactionTableName, balanceTableName))
	q.sql = q.sql.Join(fmt.Sprintf("%s on %s.user_id = %s.id", usersTableName,
		balanceTableName, usersTableName))
	return q
}

func (q *Transactions) FilterByUserId(code int) interfaces.Transactions {
	q.sql = q.sql.Where(sq.Eq{fmt.Sprintf("%s.id", usersTableName): code})
	return q
}

func (q *Transactions) OrderByLatest() interfaces.Transactions {
	q.sql = q.sql.OrderBy(fmt.Sprintf("transaction.date %s", "desc"))
	return q
}

func (q *Transactions) Search(param string, col string) interfaces.Transactions {
	param = strings.Replace(param, " ", "%", -1)
	param = fmt.Sprint("%", param, "%")

	q.sql = q.sql.Where(sq.ILike{fmt.Sprintf("transaction.%s", col): param})
	return q
}

func (q *Transactions) GetIncome(userId int, balanceId int) (*int, error) {
	var sum sql.NullInt32
	stmt := sq.Select(fmt.Sprintf("SUM(transaction.amount) FROM users "+
		"INNER JOIN balance ON users.id = balance.user_id INNER JOIN transaction "+
		"ON transaction.balance_id = balance.id WHERE users.id = %d AND transaction.balance_id= %d AND transaction.amount>0",
		userId, balanceId))
	err := q.db.Get(&sum, stmt)
	ans := int(sum.Int32)
	return &ans, err
}

func (q *Transactions) GetOutcome(userId int, balanceId int) (*int, error) {
	var sum sql.NullInt32
	stmt := sq.Select(fmt.Sprintf("SUM(transaction.amount) FROM users "+
		"INNER JOIN balance ON users.id = balance.user_id INNER JOIN transaction "+
		"ON transaction.balance_id = balance.id WHERE users.id = %d AND transaction.balance_id= %d AND transaction.amount<0",
		userId, balanceId))
	err := q.db.Get(&sum, stmt)
	ans := int(sum.Int32)
	return &ans, err
}
