package handlers

import (
	"github.com/fin-assistant/internal/services/api/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
	"strconv"
	"time"
)

const (
	foodCategory           = "food"
	transportationCategory = "transportation"
	housingCategory        = "housing"
	vehicleCategory        = "vehicle"
	communicationCategory  = "communication"
	shoppingCategory       = "shopping"
	entertainmentCategory  = "entertainment"
	expensesCategory       = "expenses"
)

func GetAllExpenses(w http.ResponseWriter, r *http.Request) {
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

	dict := make(map[string]int)
	dict = map[string]int{
		foodCategory:           0,
		transportationCategory: 0,
		housingCategory:        0,
		vehicleCategory:        0,
		communicationCategory:  0,
		shoppingCategory:       0,
		entertainmentCategory:  0,
		expensesCategory:       0,
	}

	for category, _ := range dict {
		amount, err := Transaction(r).GetExpenses(userId, category)
		if err != nil {
			Log(r).WithError(err).Error("failed to get category")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		dict[category] = *amount
	}

	response := resources.Category{
		Key: resources.Key{
			ID: time.Now().UTC().String(),
		},
		Attributes: resources.CategoryAttributes{
			Food:           dict[foodCategory],
			Transportation: dict[transportationCategory],
			Housing:        dict[housingCategory],
			Vehicle:        dict[vehicleCategory],
			Communication:  dict[communicationCategory],
			Shopping:       dict[shoppingCategory],
			Entertainment:  dict[entertainmentCategory],
			Expenses:       dict[expensesCategory],
		},
	}

	ape.Render(w, response)
}
