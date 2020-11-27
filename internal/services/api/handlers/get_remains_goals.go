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

func GetRemainGoals(w http.ResponseWriter, r *http.Request) {
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

	response := resources.RemainGoalsListResponse{
		Data: make([]resources.RemainGoals, 0, len(*balances)),
	}

	for _, balance := range *balances {
		remains, _ := Goal(r).New().CountRemainsGoals(userId, balance.ID)
		income, _ := Transaction(r).New().GetIncome(userId, balance.ID)
		outcome, _ := Transaction(r).New().GetOutcome(userId, balance.ID)
		response.Data = append(response.Data, *balance.RemainResource(*remains, *income, *outcome))
	}

	ape.Render(w, response)
}
