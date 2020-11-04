package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-gomail/gomail"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func AskRecovery(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")

	q := User(r)

	user, err := q.CheckUser(email)
	if err != nil {
		Log(r).WithError(err).Error("failed to check the user")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if user == false {
		Log(r).Error("user doesn't exist")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	recoveryKeyBytes := make([]byte, 16)
	if _, err := rand.Read(recoveryKeyBytes); err != nil {
		Log(r).WithError(err).Error("failed to create a recovery key")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	recoveryKeyString := hex.EncodeToString(recoveryKeyBytes)

	receiver := email
	body := fmt.Sprintf("Your recovery key is: %s", recoveryKeyString)

	address := Email(r).Address
	password := Email(r).Password
	
	m := gomail.NewMessage()
	m.SetHeader("From", address)
	m.SetHeader("To", receiver)
	m.SetAddressHeader("Cc", receiver, "User")
	m.SetHeader("Subject", "ShelfLife recovery")
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, address, password)

	if err := d.DialAndSend(m); err != nil {
		Log(r).WithError(err).Error("Failed to send a key")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	err = q.SetRecoveryKeyByEmail(receiver, recoveryKeyString)
	if err != nil {
		Log(r).WithError(err).Error("Failed to set a recovery key")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusOK)
}
