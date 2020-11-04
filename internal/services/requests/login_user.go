package requests

import (
	"encoding/json"
	"github.com/fin-assistant/internal/resources"
	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

func NewLoginUserRequest(r *http.Request) (resources.LoginUserResponse, error) {
	var request resources.LoginUserResponse
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, validateLoginUserRequest(request)
}

func validateLoginUserRequest(r resources.LoginUserResponse) error {
	return validation.Errors{
		"/data/attributes/email":    validation.Validate(&r.Data.Attributes.Email, validation.Required),
		"/data/attributes/password": validation.Validate(&r.Data.Attributes.Password, validation.Required),
	}.Filter()
}
