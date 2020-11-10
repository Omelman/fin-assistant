package handlers

import (
	"github.com/fin-assistant/internal/services/api/resources"
	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
	"strconv"
)

func GetAllBalance(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.Header.Get("user-id"))
	if err != nil {
		Log(r).WithError(err).Error("failed to parse user-id")
		ape.RenderErr(w, problems.InternalError())
		return
	}

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

	balances, err := Balance(r).GetAllBalances(userId)

	response := resources.GetBalanceListResponse{
		Data: make([]resources.GetBalance, 0, len(*balances)),
	}

	for _, balances := range *balances {
		response.Data = append(response.Data, *balances.Resource())
	}

	ape.Render(w, response)
}
