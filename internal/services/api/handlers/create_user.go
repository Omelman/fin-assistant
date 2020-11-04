package handlers

import (
	"encoding/json"
	"github.com/badoux/checkmail"
	"github.com/fin-assistant/internal/postgres/interfaces"
	"github.com/fin-assistant/internal/resources"
	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/validator.v2"
	"net/http"
)

//request
func NewCreateUserRequest(r *http.Request) (resources.CreateUserResponse, error) {
	var request resources.CreateUserResponse
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, validateCreateUserRequest(request)
}

func validateCreateUserRequest(r resources.CreateUserResponse) error {
	err := checkmail.ValidateHost(r.Data.Attributes.Email)
	if smtpErr, ok := err.(checkmail.SmtpError); ok && err != nil {
		return smtpErr.Err
	}

	err =  validation.Errors{
		"/implementation/attributes/firstname":	validation.Validate(&r.Data.Attributes.Firstname, validation.Required),
		"/implementation/attributes/lastname":	validation.Validate(&r.Data.Attributes.Lastname, validation.Required),
		"/implementation/attributes/password":	validation.Validate(&r.Data.Attributes.Password, validation.Required),
	}.Filter()
	if err != nil {
		return err
	}

	err = ValidateNewUser(r)
	if err != nil {
		return err
	}

	return nil
}

//regexes
type NewUserRequest struct {
	Firstname string `validate:"min=3,max=40,regexp=^[a-zA-Z]*$"`
	Lastname  string `validate:"min=3,max=40,regexp=^[a-zA-Z]*$"`
	Password  string `validate:"min=8,max=100"`
}

func ValidateNewUser (r resources.CreateUserResponse) error {
	nur := NewUserRequest{
		Firstname: r.Data.Attributes.Firstname,
		Password: r.Data.Attributes.Password,
	}
	errs := validator.Validate(nur)
	return errs
}

//
func CreateUser(w http.ResponseWriter, r *http.Request) {
	request, err := NewCreateUserRequest(r)
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
	err = q.Create(interfaces.User{
		Firstname: request.Data.Attributes.Firstname,
		Lastname:  request.Data.Attributes.Lastname,
		Email:     request.Data.Attributes.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		Log(r).WithError(err).Error("failed to create user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusOK)
}
