package requests

import (
	"encoding/json"
	"github.com/fin-assistant/internal/resources"
	"github.com/badoux/checkmail"
	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

func NewCreateUserRequest(r *http.Request) (resources.CreateUserResponse, error) {
	var request resources.CreateUserResponse
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, validateCreateUserRequest(request)
}

func validateCreateUserRequest(r resources.CreateUserResponse) error {
	err := checkmail.ValidateHost(r.Data.Attributes.Email)
	if smtpErr, ok := err.(checkmail.SmtpError); ok && err != nil {
		return smtpErr.Err
	}

	err =  validation.Errors{
		"/data/attributes/nickname":	validation.Validate(&r.Data.Attributes.Nickname, validation.Required),
		"/data/attributes/password":	validation.Validate(&r.Data.Attributes.Password, validation.Required),
	}.Filter()
	if err != nil {
		return err
	}

	err = ValidateNewUser(r)
	if err != nil {
		return err
	}

	return nil
}
