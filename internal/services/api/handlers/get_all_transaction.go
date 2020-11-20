package handlers

import (
	"github.com/fin-assistant/internal/services/api/resources"
	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
	"strconv"
	"time"
)

const (
	buyRequestListFilterPromo         = "filter[promo]"
	buyRequestListFilterPurchaseAsset = "filter[purchase_asset]"
	buyRequestListFilterBoughtAsset   = "filter[bought_asset]"
	buyRequestListFilterDateFrom      = "filter[date_from]"
	buyRequestListFilterDateTo        = "filter[date_to]"
	buyRequestListFilterSender        = "filter[sender]"
	buyRequestListFilterSeller        = "filter[seller]"
)

type GetBuyRequestList struct {
	*PageParams
	BuyRequestFilters
}

type BuyRequestFilters struct {
	Sender        *string
	Seller        *string
	Offers        []int64
	Promo         *string
	PurchaseAsset *string
	BoughtAsset   *string
	Status        *string
	From, To      *time.Time
}

func NewGetTransRequestList(r *http.Request) (*GetBuyRequestList, error) {
	request := GetBuyRequestList{}

	params, err := GetPageParams(r)
	if err != nil {
		return nil, validation.Errors{"page": errors.Wrap(err, "failed to populate")}
	}
	request.PageParams = params

	return &request, request.BuyRequestFilters.populate(r)
}

func (f *BuyRequestFilters) populate(r *http.Request) error {
	errs := validation.Errors{}

	seller := getString(r, buyRequestListFilterSeller)
	if seller != "" {
		f.Seller = &seller
	}

	sender := getString(r, buyRequestListFilterSender)
	if sender != "" {
		f.Sender = &sender
	}

	promo := getString(r, buyRequestListFilterPromo)
	if promo != "" {
		f.Promo = &promo
	}

	purchaseAsset := getString(r, buyRequestListFilterPurchaseAsset)
	if purchaseAsset != "" {
		f.PurchaseAsset = &purchaseAsset
	}

	boughtAsset := getString(r, buyRequestListFilterBoughtAsset)
	if boughtAsset != "" {
		f.BoughtAsset = &boughtAsset
	}

	return errs.Filter()
}

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
