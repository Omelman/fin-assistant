package handlers

import (
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
	"strconv"
)

func DeleteGoal(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	goalId, err := strconv.Atoi(idString)
	if err != nil {
		Log(r).WithError(err).Error("failed to convert id")
		ape.Render(w, problems.BadRequest(err))
		return
	}

	err = CheckToken(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to check token")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	err = Goal(r).DeleteGoal(goalId)
	if err != nil {
		Log(r).WithError(err).Error("failed to delete goal")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusOK)
}
