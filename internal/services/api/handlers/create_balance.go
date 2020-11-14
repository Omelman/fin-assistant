package handlers

import (
	"encoding/json"
	"github.com/fin-assistant/internal/postgres/interfaces"
	"github.com/fin-assistant/internal/services/api/resources"
	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
	"strconv"
)

//requests
func BalanceRequest(r *http.Request) (resources.BalanceResponse, error) {
	var request resources.BalanceResponse
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, validateBalanceRequest(request)
}

func validateBalanceRequest(r resources.BalanceResponse) error {
	return validation.Errors{
		"/data/attributes/currency": validation.Validate(&r.Data.Attributes.Currency, validation.Required),
	}.Filter()
}

func CreateBalance(w http.ResponseWriter, r *http.Request) {
	request, err := BalanceRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	err = CheckToken(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to check token")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	userId, err := strconv.Atoi(r.Header.Get("user-id"))
	balance := interfaces.Balance{
		Currency: request.Data.Attributes.Currency,
		UserId:   userId,
	}
	err = Balance(r).Create(balance)
	if err != nil {
		Log(r).WithError(err).Error("failed to create balance")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusOK)
}
