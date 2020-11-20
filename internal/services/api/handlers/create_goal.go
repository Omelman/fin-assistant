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
func CreateGoalRequest(r *http.Request) (resources.CreateGoalResponse, error) {
	var request resources.CreateGoalResponse
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, validateCreateGoalRequest(request)
}

func validateCreateGoalRequest(r resources.CreateGoalResponse) error {
	return validation.Errors{
		"/data/attributes/date_start":    validation.Validate(&r.Data.Attributes.DateStart, validation.Required),
		"/data/attributes/date_finished": validation.Validate(&r.Data.Attributes.DateFinish, validation.Required),
		"/data/attributes/amount":        validation.Validate(&r.Data.Attributes.Amount, validation.Required),
		"/data/attributes/description":   validation.Validate(&r.Data.Attributes.Description, validation.Required),
		"/data/attributes/balance_id":    validation.Validate(&r.Data.Attributes.BalaceId, validation.Required),
	}.Filter()
}

func CreateGoal(w http.ResponseWriter, r *http.Request) {
	request, err := CreateGoalRequest(r)
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

	newTimeStart, err := time.Parse("2006-01-02", request.Data.Attributes.DateStart)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse date start ")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	newTimeFinish, err := time.Parse("2006-01-02", request.Data.Attributes.DateFinish)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse date start ")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	newGoal := interfaces.Goal{
		DateStart:   newTimeStart,
		DateFinish:  newTimeFinish,
		Description: request.Data.Attributes.Description,
		Amount:      request.Data.Attributes.Amount,
		BalanceId:   request.Data.Attributes.BalaceId,
	}
	goalId, err := Goal(r).Create(newGoal)
	if err != nil {
		Log(r).WithError(err).Error("failed to create goal")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	response := resources.ReturnIdResponse{
		Data: resources.ReturnId{
			Attributes: resources.ReturnIdAttributes{
				Id: goalId,
			},
		},
	}
	ape.Render(w, response)
}
