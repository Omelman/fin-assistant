package handlers

import (
	"encoding/json"
	"github.com/fin-assistant/internal/resources"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")

	exists, err := User(r).CheckUser(email)
	if err != nil {
		Log(r).WithError(err).Error("failed to check user")
		ape.Render(w, problems.InternalError())
		return
	}
	if exists != true {
		Log(r).Error("not found user")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	user, err := User(r).GetByEmail(email)
	if err != nil {
		Log(r).Error("not found user")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response := resources.GetUserResponse {
		Data: resources.GetUser{
			Key: resources.NewKeyInt64(1, resources.GET_USER),
			Attributes: &resources.GetUserAttributes{
				Id: &user.ID,
				Email: &user.Email,
			},
		},
		Included: resources.Included{},
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		Log(r).WithError(err).Error("failed to marshal json")
		ape.Render(w, problems.InternalError())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(responseJson)
	if err != nil {
		Log(r).WithError(err).Error("failed to write a response")
		ape.Render(w, problems.InternalError())
		return
	}
}
