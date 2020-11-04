package cli

import (
	"context"
	"github.com/fin-assistant/internal/services/api"

	"github.com/fin-assistant/internal/config"
	"github.com/urfave/cli"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func Run(args []string) bool {
	var cfg config.Config
	log := logan.New()

	defer func() {
		if rvr := recover(); rvr != nil {
			log.WithRecover(rvr).Error("app panicked")
		}
	}()

	app := cli.NewApp()

	before := func(_ *cli.Context) error {
		getter, err := kv.FromEnv()
		if err != nil {
			return errors.Wrap(err, "failed to get config")
		}
		cfg = config.NewConfig(getter)
		log = cfg.Log()
		return nil
	}

	app.Commands = cli.Commands{
		{
			Name: "migrate",
			Subcommands: cli.Commands{
				{
					Name:   "up",
					Before: before,
					Action: func(ctx *cli.Context) error {
						return MigrateUp(cfg)
					},
				},
				{
					Name:   "down",
					Before: before,
					Action: func(ctx *cli.Context) error {
						return MigrateDown(cfg)
					},
				},
			},
		},
		{
			Name:   "run",
			Before: before,
			Action: func(_ *cli.Context) error {
				service, err := api.NewService(context.Background(), cfg)
				if err != nil {
					return errors.Wrap(err, "unable to create service")
				}

				return service.Run()
			},
		},
	}

	if err := app.Run(args); err != nil {
		log.WithError(err).Error("app finished")
		return false
	}
	return true
}
