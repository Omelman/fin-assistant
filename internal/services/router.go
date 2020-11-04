package services

import (
	"github.com/fin-assistant/internal/config"
	"github.com/fin-assistant/internal/data/postgres"
	"github.com/fin-assistant/internal/services/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gitlab.com/distributed_lab/ape"
)

func Router(s *service, cfg config.Config) chi.Router {
	r := chi.NewRouter()
	r.Use(
		middleware.SetHeader("Access-Control-Allow-Origin","*"),
		ape.RecoverMiddleware(cfg.Log()),
		ape.LoganMiddleware(cfg.Log()),
		ape.CtxMiddleWare(
			handlers.CtxLog(cfg.Log()),
			handlers.CtxUser(postgres.NewUsers(cfg.DB())),
			handlers.CtxEmail(cfg.Email()),
		),
	)

	r.Route("/", func(r chi.Router) {
		r.Route("/assistant", func(r chi.Router) {
			r.Post("/sign-up", handlers.CreateUser)
			r.Post("/sign-in", handlers.LoginUser)
			r.Post("/check_token", handlers.CheckToken)
			r.Get("/get_user/{email}", handlers.GetUser)

			r.Route("/recovery", func(r chi.Router) {
				r.Get("/{email}", handlers.AskRecovery)
				r.Post("/", handlers.CompleteRecovery)
			})
		})
	})

	return r
}
