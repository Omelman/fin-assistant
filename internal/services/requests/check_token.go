package requests

import (
	"encoding/json"
	"github.com/fin-assistant/internal/resources"
	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

func CheckTokenRequest(r *http.Request) (resources.CheckTokenResponse, error) {
	var request resources.CheckTokenResponse
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, validateCheckTokenRequest(request)
}

func validateCheckTokenRequest(r resources.CheckTokenResponse) error {
	return validation.Errors{
		"/data/attributes/email": validation.Validate(&r.Data.Attributes.Email, validation.Required),
	}.Filter()
}
