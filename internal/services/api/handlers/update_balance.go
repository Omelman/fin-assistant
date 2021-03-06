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
	"strings"
)

//requests
func UpdateBalanceRequest(r *http.Request) (resources.GetBalanceResponse, error) {
	var request resources.GetBalanceResponse
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, validateUpdateBalanceRequest(request)
}

func validateUpdateBalanceRequest(r resources.GetBalanceResponse) error {
	return validation.Errors{
		"/data/attributes/currency":   validation.Validate(&r.Data.Attributes.Currency, validation.Required),
		"/data/attributes/balance_id": validation.Validate(&r.Data.Attributes.BalanceId, validation.Required),
	}.Filter()
}

func UpdateBalance(w http.ResponseWriter, r *http.Request) {
	request, err := UpdateBalanceRequest(r)
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

	balances, err := Balance(r).GetAllBalances(userId)
	if err != nil {
		Log(r).WithError(err).Error("failed to get balances")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	for _, value := range *balances {
		if strings.ToUpper(request.Data.Attributes.Currency) == value.Currency {
			Log(r).Error("This currency already exist")
			ape.RenderErr(w, problems.BadRequest(err)...)
			return
		}
	}

	balance := interfaces.Balance{
		ID:       request.Data.Attributes.BalanceId,
		Currency: strings.ToUpper(request.Data.Attributes.Currency),
		UserId:   userId,
	}
	err = Balance(r).Update(balance)
	if err != nil {
		Log(r).WithError(err).Error("failed to update balance")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusOK)
}
