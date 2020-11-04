package interfaces

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
	Firstname   string `db:"firstname" structs:"firstname"`
	Lastname    string `db:"lastname" structs:"lastname"`
	Email       string `db:"email" structs:"email"`
	Password    []byte `db:"password" structs:"password"`
	Token       string `db:"token" structs:"token"`
}
