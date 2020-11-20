package api

import (
	"context"
	"fmt"
	"github.com/fin-assistant/internal/config"
	"github.com/fin-assistant/internal/postgres/implementation"
	"github.com/fin-assistant/internal/postgres/interfaces"
	"github.com/go-gomail/gomail"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"time"
)

type GoalsChecker struct {
	Goals interfaces.Goals
	Email config.Email
	log   *logan.Entry
}

func NewGoalsChecker(cfg config.Config) *GoalsChecker {
	return &GoalsChecker{
		Goals: implementation.NewGoal(cfg.DB()),
		Email: cfg.Email(),
		log:   cfg.Log(),
	}
}

func (p *GoalsChecker) Run(ctx context.Context) {
	p.log.Info("goals_checker started")
	running.WithBackOff(ctx, p.log, "goals_checker", p.checkPending,
		30*time.Second, 30*time.Second, 30*time.Second)
	p.log.Info("goals_checker finished")
}

func (p *GoalsChecker) checkPending(_ context.Context) error {
	goals, err := p.Goals.New().
		FilterByStatus(time.Now().Format("2006-01-02")).
		Select()
	if err != nil {
		return errors.Wrap(err, "failed to get pending goals")
	}

	for _, goal := range goals {
		p.log.Debug(fmt.Sprintf("processing goal %d", goal.ID))

		receiver := "cr.frog03@gmail.com"
		body := fmt.Sprintf("Your goal ended : %d", goal.ID)

		address := p.Email.Address
		password := p.Email.Password

		m := gomail.NewMessage()
		m.SetHeader("From", address)
		m.SetHeader("To", receiver)
		m.SetAddressHeader("Cc", receiver, "User")
		m.SetHeader("Subject", "Fin-assistant")
		m.SetBody("text/plain", body)

		d := gomail.NewDialer("smtp.gmail.com", 587, address, password)

		if err := d.DialAndSend(m); err != nil {
			p.log.WithError(err)
		}
	}

	return nil
}
