package handlers

import (
	"github.com/fin-assistant/internal/services/api/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
	"strconv"
)

func GetAllTransaction(w http.ResponseWriter, r *http.Request) {
	err := CheckToken(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to check token")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	userId, err := strconv.Atoi(r.Header.Get("user-id"))
	if err != nil {
		Log(r).WithError(err).Error("failed to parse user-id")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	transactions, err := Transaction(r).GetAllTransaction(userId)

	response := resources.CreateTransactionListResponse{
		Data: make([]resources.CreateTransaction, 0, len(*transactions)),
	}

	for _, transaction := range *transactions {
		response.Data = append(response.Data, *transaction.Resource())
	}

	ape.Render(w, response)
}
