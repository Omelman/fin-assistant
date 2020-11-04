package requests

import (
	"github.com/fin-assistant/internal/resources"
	"gopkg.in/validator.v2"
)

type NewUserRequest struct {
	Nickname string `validate:"min=3,max=40,regexp=^[a-zA-Z]*$"`
	Password string `validate:"min=8,max=100"`
}

type RecoveryRequest struct {
	Password string `validate:"min=8,max=100"`
}

func ValidateNewUser (r resources.CreateUserResponse) error {
	nur := NewUserRequest{
		Nickname: r.Data.Attributes.Nickname,
		Password: r.Data.Attributes.Password,
	}
	errs := validator.Validate(nur)
	return errs
}

func ValidateRecovery (r resources.CompleteRecoveryResponse) error {
	nur := RecoveryRequest{
		Password: r.Data.Attributes.NewPassword,
	}
	errs := validator.Validate(nur)
	return errs
}
