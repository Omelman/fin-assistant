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
	"time"
)

//requests
func CreateTransactionRequest(r *http.Request) (resources.CreateTransactionResponse, error) {
	var request resources.CreateTransactionResponse
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, validateCreateTransactionRequest(request)
}

func validateCreateTransactionRequest(r resources.CreateTransactionResponse) error {
	return validation.Errors{
		"/data/attributes/date":        validation.Validate(&r.Data.Attributes.Date, validation.Required),
		"/data/attributes/amount":      validation.Validate(&r.Data.Attributes.Amount, validation.Required),
		"/data/attributes/category":    validation.Validate(&r.Data.Attributes.Category, validation.Required),
		"/data/attributes/description": validation.Validate(&r.Data.Attributes.Description, validation.Required),
		"/data/attributes/balance_id":  validation.Validate(&r.Data.Attributes.BalaceId, validation.Required),
		"/data/attributes/include":     validation.Validate(&r.Data.Attributes.Include, validation.NotNil),
	}.Filter()
}

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	request, err := CreateTransactionRequest(r)
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

	newTime, err := time.Parse("2006-01-02", request.Data.Attributes.Date)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse date")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	newTransaction := interfaces.Transaction{
		Date:        newTime,
		Description: request.Data.Attributes.Description,
		Amount:      request.Data.Attributes.Amount,
		Category:    request.Data.Attributes.Category,
		Include:     *request.Data.Attributes.Include,
		BalanceId:   request.Data.Attributes.BalaceId,
	}
	transactionId, err := Transaction(r).Create(newTransaction)
	if err != nil {
		Log(r).WithError(err).Error("failed to create transaction")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	response := resources.ReturnIdResponse{
		Data: resources.ReturnId{
			Attributes: resources.ReturnIdAttributes{
				Id: transactionId,
			},
		},
	}
	ape.Render(w, response)
}
