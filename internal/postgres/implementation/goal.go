package implementation

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/fin-assistant/internal/postgres/interfaces"
	"gitlab.com/distributed_lab/kit/pgdb"
	"time"
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

func (q *Goals) Create(user interfaces.Goal) (int, error) {
	clauses := structs.Map(user)

	var id int
	stmt := sq.Insert(goalTableName).SetMap(clauses).Suffix("returning id")
	err := q.db.Get(&id, stmt)
	return id, err
}

func (q *Goals) Transaction(fn func(q interfaces.Goals) error) (err error) {
	return q.db.Transaction(func() error {
		return fn(q)
	})
}

func (q *Goals) Select() ([]interfaces.Goal, error) {
	var result []interfaces.Goal
	err := q.db.Select(&result, q.sql)
	return result, err
}

func (q *Goals) Get() (*interfaces.Goal, error) {
	var user interfaces.Goal
	err := q.db.Get(&user, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, err
}

func (q *Goals) GetEmail(goalId int) (*string, error) {
	var email string
	stmt := sq.Select(fmt.Sprintf("users.email FROM users INNER JOIN balance ON balance.user_id = users.id "+
		"INNER JOIN goal ON goal.balance_id = balance.id WHERE goal.id = %d", goalId))
	err := q.db.Get(&email, stmt)
	return &email, err
}

func (q *Goals) GetAllGoals(userId int) (*[]interfaces.Goal, error) {
	var goals []interfaces.Goal
	stmt := sq.Select(fmt.Sprintf("goal.* FROM users "+
		"INNER JOIN balance ON users.id = balance.user_id INNER JOIN goal "+
		"ON goal.balance_id = balance.id WHERE users.id = %d", userId))
	err := q.db.Select(&goals, stmt)
	return &goals, err
}

func (q *Goals) GetById(id int) (*interfaces.Goal, error) {
	var user interfaces.Goal
	stmt := sq.Select("*").From(goalTableName).Where("id = ?", id)
	err := q.db.Get(&user, stmt)
	return &user, err
}

func (q *Goals) Update(goal interfaces.Goal, goalId int) error {
	clauses := structs.Map(goal)

	stmt := sq.Update(goalTableName).SetMap(clauses).Where("id = ?", goalId)
	err := q.db.Exec(stmt)
	return err
}

func (q *Goals) DeleteGoal(goalId int) error {
	stmt := sq.Delete(goalTableName).Where("id = ?", goalId)
	err := q.db.Exec(stmt)
	return err
}

func (q *Goals) FilterByStatus(date string) interfaces.Goals {
	q.sql = q.sql.Where("date_finish = ?", date)
	return q
}

func (q *Goals) CountRemainsGoals(userId int, balanceId int) (*int, error) {
	var count sql.NullInt32
	timeNow := time.Now().Format("2006-01-02")
	stmt := sq.Select(fmt.Sprintf("COUNT(goal.id) FROM users "+
		"INNER JOIN balance ON users.id = balance.user_id INNER JOIN goal "+
		"ON goal.balance_id = balance.id WHERE users.id = %d AND goal.date_finish > to_date('%s','yyyy-mm-dd')  "+
		"AND balance_id = %d",
		userId, timeNow, balanceId))
	err := q.db.Get(&count, stmt)
	ans := int(count.Int32)
	return &ans, err
}
