package handlers

import (
	"github.com/fin-assistant/internal/services/api/resources"
	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//query
const (
	transactionListFilterCategory = "filter[category]"
	transactionListFilterBalance  = "filter[balance]"
	transactionListFilterDateFrom = "filter[date_from]"
	transactionListFilterDateTo   = "filter[date_to]"
)

type GetTransactionList struct {
	TransactionFilters
}

func getString(request *http.Request, name string) string {
	result := chi.URLParam(request, name)
	if result != "" {
		return strings.TrimSpace(result)
	}

	return strings.TrimSpace(request.URL.Query().Get(name))
}

type TransactionFilters struct {
	Category, Balance *string
	From, To          *time.Time
}

func NewGetTransRequestList(r *http.Request) (*GetTransactionList, error) {
	request := GetTransactionList{}
	return &request, request.TransactionFilters.populate(r)
}

func (f *TransactionFilters) populate(r *http.Request) error {
	errs := validation.Errors{}

	category := getString(r, transactionListFilterCategory)
	if category != "" {
		f.Category = &category
	}

	balance := getString(r, transactionListFilterBalance)
	if balance != "" {
		f.Balance = &balance
	}

	if getString(r, transactionListFilterDateFrom) != "" {
		dateFrom, err := time.Parse("2006-01-02", getString(r, transactionListFilterDateFrom))
		if err != nil {
			errs[transactionListFilterDateFrom] = err
		}
		if &dateFrom != nil {
			f.From = &dateFrom
		}
	}

	if getString(r, transactionListFilterDateTo) != "" {
		dateTo, err := time.Parse("2006-01-02", getString(r, transactionListFilterDateTo))
		if err != nil {
			errs[transactionListFilterDateTo] = err
		}
		if &dateTo != nil {
			f.To = &dateTo
		}
	}

	return errs.Filter()
}

//endpoint start here
func GetAllTransaction(w http.ResponseWriter, r *http.Request) {
	err := CheckToken(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to check token")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	request, err := NewGetTransRequestList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	q := Transaction(r).New()
	if request.Category != nil {
		q = q.FilterByCategory(*request.Category)
	}
	if request.Balance != nil {
		q = q.FilterByBalance(*request.Balance)
	}
	if request.From != nil {
		q = q.FilterOnlyAfter(*request.From)
	}
	if request.To != nil {
		q = q.FilterOnlyBefore(*request.To)
	}

	userId, err := strconv.Atoi(r.Header.Get("user-id"))
	if err != nil {
		Log(r).WithError(err).Error("failed to parse user-id")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	transactions, err := q.UserJoined().FilterByUserId(userId).Select()

	response := resources.CreateTransactionListResponse{
		Data: make([]resources.CreateTransaction, 0, len(*transactions)),
	}

	for _, transaction := range *transactions {
		response.Data = append(response.Data, *transaction.Resource())
	}

	ape.Render(w, response)
}
