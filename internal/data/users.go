package data

type Users interface {
	New() Users

	Create(User User) error
	Get() (*User, error)
	GetByEmail(email string) (*User, error)

	SetTokenByEmail(email string, token string) error
	SetRecoveryKeyByEmail(email string, recoveryKey string) error
	UpdatePasswordByKey(recoveryKey string, newPassword []byte) error

	ClearKeyByEmail(email string) error

	CheckUser(email string) (bool, error)
	Transaction(fn func(q Users) error) (err error)
}

type User struct {
	ID          int32  `db:"id" structs:"-"`
	Nickname    string `db:"nickname" structs:"nickname"`
	Email       string `db:"email" structs:"email"`
	Details     []byte `db:"details" structs:"-"`
	Password    []byte `db:"password" structs:"password"`
	Token       string `db:"token" structs:"token"`
	RecoveryKey string `db:"recovery_key" structs:"recovery_key"`
}
