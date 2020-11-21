package interfaces

import (
	"github.com/fin-assistant/internal/services/api/resources"
	"time"
)

type Goals interface {
	New() Goals

	Create(User Goal) (int, error)
	Select() ([]Goal, error)
	Get() (*Goal, error)
	GetById(id int) (*Goal, error)
	GetAllGoals(userId int) (*[]Goal, error)
	GetEmail(goalId int) (*string, error)
	Update(goal Goal, goalId int) error

	FilterByStatus(date string) Goals

	DeleteGoal(goalId int) error

	Transaction(fn func(q Goals) error) (err error)
}

type Goal struct {
	ID          int       `db:"id" structs:"-"`
	DateStart   time.Time `db:"date_start" structs:"date_start"`
	DateFinish  time.Time `db:"date_finish" structs:"date_finish"`
	Description string    `db:"description" structs:"description"`
	Amount      int       `db:"amount" structs:"amount"`
	BalanceId   int       `db:"balance_id" structs:"-"`
}

func (r *Goal) Resource() *resources.CreateGoal {
	return &resources.CreateGoal{
		Attributes: resources.CreateGoalAttributes{
			DateStart:   r.DateStart.Format("2006-01-02"),
			DateFinish:  r.DateFinish.Format("2006-01-02"),
			Description: r.Description,
			Amount:      r.Amount,
			BalaceId:    r.BalanceId,
			GoalId:      r.ID,
		},
	}
}
