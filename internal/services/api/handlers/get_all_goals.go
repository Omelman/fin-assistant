package handlers

import (
	"github.com/fin-assistant/internal/services/api/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
	"strconv"
)

func GetAllGoals(w http.ResponseWriter, r *http.Request) {
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

	goals, err := Goal(r).GetAllGoals(userId)

	response := resources.CreateGoalListResponse{
		Data: make([]resources.CreateGoal, 0, len(*goals)),
	}

	for _, goal := range *goals {
		response.Data = append(response.Data, *goal.Resource())
	}

	ape.Render(w, response)
}
