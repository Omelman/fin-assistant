package handlers

import (
	"github.com/fin-assistant/internal/data"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"

	"github.com/fin-assistant/internal/services/requests"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewCreateUserRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	passwordBytes := []byte(request.Data.Attributes.Password)

	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		Log(r).WithError(err).Error("failed to generate from password")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	q := User(r)
	err = q.Create(data.User{
		Nickname: request.Data.Attributes.Nickname,
		Email:    request.Data.Attributes.Email,
		Password: hashedPassword,
	})
	if err != nil {
		Log(r).WithError(err).Error("failed to create user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusOK)
}
