package api

import (
	"github.com/fin-assistant/internal/config"
	"github.com/fin-assistant/internal/postgres/implementation"
	"github.com/fin-assistant/internal/services/api/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gitlab.com/distributed_lab/ape"
)

func Router(cfg config.Config) chi.Router {
	r := chi.NewRouter()
	r.Use(
		middleware.SetHeader("Access-Control-Allow-Origin", "*"),
		middleware.SetHeader("Access-Control-Allow-Methods", "*"),
		middleware.SetHeader("Access-Control-Allow-Headers", "user-id,token"),
		ape.RecoverMiddleware(cfg.Log()),
		ape.LoganMiddleware(cfg.Log()),
		ape.CtxMiddleWare(
			handlers.CtxLog(cfg.Log()),
			handlers.CtxUser(implementation.NewUsers(cfg.DB())),
			handlers.CtxBalance(implementation.NewBalance(cfg.DB())),
			handlers.CtxWork(implementation.NewWork(cfg.DB())),
			handlers.CtxTransaction(implementation.NewTransaction(cfg.DB())),
			handlers.CtxGoal(implementation.NewGoal(cfg.DB())),
			handlers.CtxEmail(cfg.Email()),
		),
	)

	r.Route("/", func(r chi.Router) {
		r.Route("/assistant", func(r chi.Router) {
			r.Post("/sign-up", handlers.CreateUser)
			r.Post("/sign-in", handlers.LoginUser)
			r.Post("/check_token", handlers.CheckToken)
			r.Get("/get_user/{email}", handlers.GetUser)
			r.Post("/balance/create", handlers.CreateBalance)
		})
	})

	return r
}
