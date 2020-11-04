package handlers

import (
	"github.com/fin-assistant/internal/services/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func CompleteRecovery(w http.ResponseWriter, r *http.Request) {
	request, err := requests.CompleteRecoveryRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse the request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	givenKey := request.Data.Attributes.Key

	q := User(r)
	data, err := q.GetByEmail(request.Data.Attributes.Email)
	if err != nil {
		Log(r).WithError(err).Error("a user doesn't exist")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if data.RecoveryKey != givenKey {
		Log(r).Error("wrong recovery key")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	givenPasswordBytes := []byte(request.Data.Attributes.NewPassword)
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(givenPasswordBytes, bcrypt.DefaultCost)

	err = User(r).UpdatePasswordByKey(givenKey, hashedPasswordBytes)
	if err != nil {
		Log(r).WithError(err).Error("failed to update a password")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	//TODO Clean recovery_key field in order to prevent reuse

	w.WriteHeader(http.StatusOK)
}
