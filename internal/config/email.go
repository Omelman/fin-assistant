package config

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
)

type Email struct {
	Address  string `fig:"address,required"`
	Password string `fig:"password,required"`
}

func (c *config) Email() Email {
	return c.email.Do(func() interface{} {
		var result Email

		err := figure.
			Out(&result).
			From(kv.MustGetStringMap(c.getter, "email")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out email"))
		}

		return result
	}).(Email)
}
