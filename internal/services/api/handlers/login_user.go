package handlers

import (
	"crypto/rand"
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"strconv"
	"time"

	"encoding/hex"
	"github.com/fin-assistant/internal/services/api/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

//request
func NewLoginUserRequest(r *http.Request) (resources.LoginUserResponse, error) {
	var request resources.LoginUserResponse
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, validateLoginUserRequest(request)
}

func validateLoginUserRequest(r resources.LoginUserResponse) error {
	return validation.Errors{
		"/data/attributes/email":    validation.Validate(&r.Data.Attributes.Email, validation.Required),
		"/data/attributes/password": validation.Validate(&r.Data.Attributes.Password, validation.Required),
	}.Filter()
}

//
func LoginUser(w http.ResponseWriter, r *http.Request) {
	request, err := NewLoginUserRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	data, err := User(r).GetByEmail(request.Data.Attributes.Email)
	if err != nil {
		Log(r).WithError(err).Error("failed to get user by email")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	rightPassword := data.Password
	receivedPassword := []byte(request.Data.Attributes.Password)

	err = bcrypt.CompareHashAndPassword(rightPassword, receivedPassword)
	if err != nil {
		Log(r).WithError(err).Error("failed to compare passwords")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	tokenBytes := make([]byte, 32)
	rand.Read(tokenBytes)
	tokenString := hex.EncodeToString(tokenBytes)

	err = User(r).SetTokenByEmail(request.Data.Attributes.Email, tokenString)
	if err != nil {
		Log(r).WithError(err).Error("failed to token by email")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	w.Header().Set("token", tokenString)
	w.Header().Set("user-id", strconv.Itoa(data.ID))
	response := resources.CreateUserResponse{
		Data: resources.CreateUser{
			Key: resources.Key{
				ID:   time.Now().UTC().String(),
				Type: resources.LOGIN_USER,
			},
			Attributes: resources.CreateUserAttributes{
				Email:     data.Email,
				Firstname: data.Firstname,
				Lastname:  data.Lastname,
			},
		},
	}

	ape.Render(w, response)
}
