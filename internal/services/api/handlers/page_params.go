package handlers

import (
	"fmt"
	"github.com/fin-assistant/internal/services/api/resources"
	"math"
	"net/http"
	"strconv"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	pageParamLimit  = "page[limit]"
	pageParamNumber = "page[number]"
	pageParamCursor = "page[cursor]"
	pageParamOrder  = "page[order]"
)

const defaultLimit uint64 = 15
const maxLimit uint64 = 100

const (
	// OrderAscending is used to indicate an ascending order in request params
	OrderAscending = "asc"

	// OrderDescending is used to indicate an descending order in request params
	OrderDescending = "desc"
)

//OrderType - represents sorting order of the query
type OrderType string

//Invert - inverts order by
func (o OrderType) Invert() OrderType {
	switch o {
	case OrderDesc:
		return OrderAsc
	case OrderAsc:
		return OrderDesc
	default:
		panic(errors.From(errors.New("unexpected order type"), logan.F{
			"order_type": o,
		}))
	}
}

const (
	// OrderAsc - ascending order
	OrderAsc OrderType = "asc"
	// OrderDesc - descending order
	OrderDesc OrderType = "desc"
)

type PageParams struct {
	Limit      uint64    `url:"page[limit]"`
	PageNumber uint64    `url:"page[number]"`
	Order      OrderType `url:"page[order]"`
}

func (p *PageParams) ApplyTo(sql sq.SelectBuilder, cols ...string) sq.SelectBuilder {
	offset := p.Limit * p.PageNumber

	sql = sql.Limit(p.Limit).Offset(offset)

	switch p.Order {
	case OrderAsc:
		for _, col := range cols {
			sql = sql.OrderBy(fmt.Sprintf("%s %s", col, "asc"))
		}
	case OrderDesc:
		for _, col := range cols {
			sql = sql.OrderBy(fmt.Sprintf("%s %s", col, "desc"))
		}
	default:
		panic(errors.From(errors.New("unexpected order type"), logan.F{
			"order_type": p.Order,
		}))
	}

	return sql
}

func GetPageParams(request *http.Request) (*PageParams, error) {
	limit, err := getLimit(request, defaultLimit, maxLimit)
	if err != nil {
		return nil, err
	}

	pageNumber, err := getPageNumber(request)
	if err != nil {
		return nil, err
	}

	order, err := getOrder(request)
	if err != nil {
		return nil, err
	}

	return &PageParams{
		Order:      order,
		Limit:      limit,
		PageNumber: pageNumber,
	}, nil
}

func getOrder(request *http.Request) (OrderType, error) {
	order := getString(request, pageParamOrder)
	switch order {
	case OrderAscending, OrderDescending:
		return OrderType(order), nil
	case "":
		return OrderAscending, nil
	default:
		return OrderDescending, validation.Errors{
			pageParamOrder: fmt.Errorf("allowed order types: %s, %s", OrderAscending, OrderDescending),
		}
	}
}

func getPageNumber(request *http.Request) (uint64, error) {
	result, err := getUint64(request, pageParamNumber)
	if err != nil {
		return 0, validation.Errors{
			pageParamNumber: err,
		}
	}

	return result, nil
}

func getLimit(request *http.Request, defaultLimit, maxLimit uint64) (uint64, error) {
	result, err := getUint64(request, pageParamLimit)
	if err != nil {
		return 0, validation.Errors{
			pageParamLimit: errors.New("Must be a valid uint64 value"),
		}
	}

	if result == 0 {
		return defaultLimit, nil
	}

	if result > maxLimit {
		return 0, validation.Errors{
			pageParamLimit: fmt.Errorf("limit must not exceed %d", maxLimit),
		}
	}

	return result, nil
}
func getCursor(request *http.Request) (uint64, error) {
	result, err := getUint64(request, pageParamCursor)
	if err != nil {
		return 0, validation.Errors{
			pageParamCursor: err,
		}
	}

	if result > math.MaxInt64 {
		return 0, validation.Errors{
			pageParamCursor: fmt.Errorf("cursor %d exceed max allowed %d", result, math.MaxInt64),
		}
	}

	return result, nil
}

func getStringSlice(r *http.Request, name string) []string {
	concatenated := getString(r, name)
	if concatenated == "" {
		return nil
	}

	return strings.Split(concatenated, ",")
}

func getString(request *http.Request, name string) string {
	result := chi.URLParam(request, name)
	if result != "" {
		return strings.TrimSpace(result)
	}

	return strings.TrimSpace(request.URL.Query().Get(name))
}

func getUint64(request *http.Request, name string) (uint64, error) {
	strVal := getString(request, name)
	if strVal == "" {
		return 0, nil
	}

	return strconv.ParseUint(strVal, 0, 64)
}

func GetOffsetLinks(r *http.Request, p PageParams) *resources.Links {
	result := resources.Links{
		Next: getOffsetLink(r, p.PageNumber+1, p.Limit, p.Order),
		Self: getOffsetLink(r, p.PageNumber, p.Limit, p.Order),
	}

	return &result
}

func getOffsetLink(r *http.Request, pageNumber, limit uint64, order OrderType) string {
	u := r.URL
	query := u.Query()
	query.Set(pageParamNumber, strconv.FormatUint(pageNumber, 10))
	query.Set(pageParamLimit, strconv.FormatUint(limit, 10))
	query.Set(pageParamOrder, string(order))
	u.RawQuery = query.Encode()
	return u.String()
}
