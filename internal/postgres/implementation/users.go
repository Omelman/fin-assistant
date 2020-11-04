package implementation

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/fin-assistant/internal/postgres/interfaces"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const usersTableName = "users"

var usersSelect = sq.Select("users.*").From(usersTableName)
var usersUpdate = sq.Update(usersTableName)

type Users struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
	upd sq.UpdateBuilder
}

func NewUsers(db *pgdb.DB) interfaces.Users {
	return &Users{
		db:  db.Clone(),
		sql: usersSelect,
		upd: usersUpdate,
	}
}

func (q *Users) New() interfaces.Users {
	return NewUsers(q.db)
}

func (q *Users) Create(user interfaces.User) error {
	clauses := structs.Map(user)

	var id int64
	stmt := sq.Insert(usersTableName).SetMap(clauses).Suffix("returning id")
	err := q.db.Get(&id, stmt)
	return err
}

func (q *Users) Transaction(fn func(q interfaces.Users) error) (err error) {
	return q.db.Transaction(func() error {
		return fn(q)
	})
}

func (q *Users) Get() (*interfaces.User, error) {
	var user interfaces.User
	err := q.db.Get(&user, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, err
}

func (q *Users) GetByEmail(email string) (*interfaces.User, error) {
	var user interfaces.User
	stmt := sq.Select("*").From(usersTableName).Where("email = ?", email)
	err := q.db.Get(&user, stmt)
	return &user, err
}

func (q *Users) SetTokenByEmail(email string, token string) error {
	stmt := usersUpdate.
		Where("email = ?", email).
		Set("token", token)
	res, err := q.db.ExecWithResult(stmt)
	if err != nil {
		return errors.Wrap(err, "unable to update row")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "unable to get affected rows")
	}
	if rowsAffected == 0 {
		return errors.Wrap(nil, "Raw somehow hasn't been updated")
	}
	return nil
}

func (q *Users) SetRecoveryKeyByEmail(email string, recoveryKey string) error {
	stmt := usersUpdate.
		Where("email = ?", email).
		Set("recovery_key", recoveryKey)
	res, err := q.db.ExecWithResult(stmt)
	if err != nil {
		return errors.Wrap(err, "unable to update row")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "unable to get affected rows")
	}
	if rowsAffected == 0 {
		return errors.Wrap(nil, "Raw somehow hasn't been updated")
	}
	return nil
}

func (q *Users) UpdatePasswordByKey(recoveryKey string, newPassword []byte) error {
	stmt := usersUpdate.
		Where("recovery_key = ?", recoveryKey).
		Set("password", newPassword)
	res, err := q.db.ExecWithResult(stmt)
	if err != nil {
		return errors.Wrap(err, "unable to update row")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "unable to get affected rows")
	}
	if rowsAffected == 0 {
		return errors.Wrap(nil, "Raw somehow hasn't been updated")
	}
	return nil
}

func (q *Users) ClearKeyByEmail(email string) error {
	stmt := usersUpdate.
		Where("email = ?", email).
		Set("recovery_key", "")
	res, err := q.db.ExecWithResult(stmt)
	if err != nil {
		return errors.Wrap(err, "unable to update row")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "unable to get affected rows")
	}
	if rowsAffected == 0 {
		return errors.Wrap(nil, "Raw somehow hasn't been updated")
	}
	return nil
}

func (q *Users) CheckUser(email string) (bool, error) {
	var count bool
	stmt := sq.Select("COUNT(1)").From(usersTableName).Where("email = ?", email)
	err := q.db.Get(&count, stmt)
	return count, err
}

