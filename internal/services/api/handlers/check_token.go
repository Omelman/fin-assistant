package handlers

import (
	"encoding/json"
	"github.com/fin-assistant/internal/resources"
	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

//requests
func CheckTokenRequest(r *http.Request) (resources.CheckTokenResponse, error) {
	var request resources.CheckTokenResponse
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, validateCheckTokenRequest(request)
}

func validateCheckTokenRequest(r resources.CheckTokenResponse) error {
	return validation.Errors{
		"/implementation/attributes/email": validation.Validate(&r.Data.Attributes.Email, validation.Required),
	}.Filter()
}

//
func CheckToken(w http.ResponseWriter, r *http.Request) {
	request, err := CheckTokenRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	user, err := User(r).GetByEmail(request.Data.Attributes.Email)
	if err != nil {
		Log(r).WithError(err).Error("failed to get user by email")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if user.Token != r.Header.Get("token") {
		Log(r).Error("wrong token")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"header/token": errors.New("wrong token")})...)
		return
	}

	response := resources.CheckTokenResponse {
		Data: resources.CheckToken{
			Key: resources.NewKeyInt64(1, resources.GET_USER),
			Attributes: resources.CheckTokenAttributes{
				Id: &user.ID,
				Email: user.Email,
			},
		},
		Included: resources.Included{},
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		Log(r).WithError(err).Error("failed to marshal json")
		ape.Render(w, problems.InternalError())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(responseJson)
	if err != nil {
		Log(r).WithError(err).Error("failed to write a response")
		ape.Render(w, problems.InternalError())
		return
	}
}
