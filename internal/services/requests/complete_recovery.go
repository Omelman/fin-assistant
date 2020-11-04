package requests

import (
	"encoding/json"
	"github.com/fin-assistant/internal/resources"
	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

func CompleteRecoveryRequest(r *http.Request) (resources.CompleteRecoveryResponse, error) {
	var request resources.CompleteRecoveryResponse
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, validateCompleteRecoveryRequest(request)
}

func validateCompleteRecoveryRequest(r resources.CompleteRecoveryResponse) error {
	err := validation.Errors{
		"/data/attributes/key": validation.Validate(&r.Data.Attributes.Key, validation.Required),
		"/data/attributes/email": validation.Validate(&r.Data.Attributes.Email, validation.Required),
		"/data/attributes/new_password": validation.Validate(&r.Data.Attributes.NewPassword, validation.Required),
	}.Filter()

	if err != nil {
		return err
	}

	err = ValidateRecovery(r)

	return err
}
