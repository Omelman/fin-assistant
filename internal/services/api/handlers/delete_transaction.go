package handlers

import (
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
	"strconv"
)

func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	transactionId, err := strconv.Atoi(idString)
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

	userId, err := strconv.Atoi(r.Header.Get("user-id"))

	err = Transaction(r).DeleteTransaction(userId, transactionId)
	if err != nil {
		Log(r).WithError(err).Error("failed to delete balance")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusOK)
}
