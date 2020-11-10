package handlers

import (
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
	"strconv"
)

func CheckSession(w http.ResponseWriter, r *http.Request) {
	err := CheckToken(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to check token")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func CheckToken(r *http.Request) error {
	userId, err := strconv.Atoi(r.Header.Get("user-id"))
	if err != nil {
		Log(r).WithError(err).Error("failed to parse user-id")
		return err
	}

	user, err := User(r).GetById(userId)
	if err != nil {
		Log(r).WithError(err).Error("failed to get user by id")
		return err
	}

	if user.Token != r.Header.Get("token") {
		Log(r).Error("wrong token")
		return err
	}

	return nil
}
