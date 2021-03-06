package interfaces

import (
	"github.com/fin-assistant/internal/services/api/resources"
	"time"
)

type Transactions interface {
	New() Transactions

	Create(User Transaction) (int, error)

	Get() (*Transaction, error)
	GetById(id int) (*Transaction, error)
	Select() (*[]Transaction, error)
	//statistics
	GetExpenses(userId int, category string) (*int, error)
	GetOutcome(userId int, balanceId int) (*int, error)
	GetIncome(userId int, balanceId int) (*int, error)
	//
	Update(transaction Transaction, transactionId int) error
	DeleteTransaction(transactionId int) error

	FilterByCategory(code string) Transactions
	FilterByBalance(code string) Transactions
	FilterByUserId(code int) Transactions
	FilterOnlyBefore(time time.Time) Transactions
	FilterOnlyAfter(time time.Time) Transactions

	Search(param string, col string) Transactions
	OrderByLatest() Transactions
	UserJoined() Transactions

	Transaction(fn func(q Transactions) error) (err error)
}

type Transaction struct {
	ID          int       `db:"id" structs:"-"`
	Date        time.Time `db:"date" structs:"date"`
	Description string    `db:"description" structs:"description"`
	Amount      int       `db:"amount" structs:"amount"`
	Category    string    `db:"category" structs:"category"`
	Include     bool      `db:"include" structs:"include"`
	BalanceId   int       `db:"balance_id" structs:"balance_id"`
}

func (r *Transaction) Resource() *resources.CreateTransaction {
	return &resources.CreateTransaction{
		Attributes: resources.CreateTransactionAttributes{
			Date:          r.Date.Format("2006-01-02"),
			Description:   r.Description,
			Amount:        r.Amount,
			Category:      r.Category,
			Include:       r.Include,
			BalaceId:      r.BalanceId,
			TransactionId: r.ID,
		},
	}
}
