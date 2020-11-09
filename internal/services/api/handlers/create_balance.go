package handlers

import (
	"encoding/json"
	"github.com/fin-assistant/internal/services/api/resources"
	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
	"strconv"
)

//requests
func CreateBalanceRequest(r *http.Request) (resources.CreateBalanceResponse, error) {
	var request resources.CreateBalanceResponse
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, validateCreateBalanceRequest(request)
}

func validateCreateBalanceRequest(r resources.CreateBalanceResponse) error {
	return validation.Errors{
		"/data/attributes/currency": validation.Validate(&r.Data.Attributes.Currency, validation.Required),
	}.Filter()
}

func CreateBalance(w http.ResponseWriter, r *http.Request) {
	_, err := CreateBalanceRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	userId, _ := strconv.Atoi(r.Header.Get("user-id"))
	user, err := User(r).GetById(userId)
	if err != nil {
		Log(r).WithError(err).Error("failed to get user by id")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if user.Token != r.Header.Get("token") {
		Log(r).Error("wrong token")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"header/token": errors.New("wrong token")})...)
		return
	}

	response := resources.CheckTokenResponse{
		Data: resources.CheckToken{
			Key: resources.NewKeyInt64(1, resources.GET_USER),
			Attributes: resources.CheckTokenAttributes{
				Id:    &user.ID,
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
