package api

import (
	"github.com/fin-assistant/internal/config"
	"github.com/fin-assistant/internal/postgres/implementation"
	"github.com/fin-assistant/internal/services/api/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"gitlab.com/distributed_lab/ape"
)

func Router(cfg config.Config) chi.Router {
	r := chi.NewRouter()
	r.Use(
		cors.AllowAll().Handler,
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
			//user
			r.Post("/sign-up", handlers.CreateUser)
			r.Post("/sign-in", handlers.LoginUser)
			r.Get("/check_session", handlers.CheckSession)
			r.Get("/get_user/{email}", handlers.GetUser)
			//balance
			r.Post("/balance", handlers.CreateBalance)
			r.Get("/balance", handlers.GetAllBalance)
			r.Put("/balance", handlers.UpdateBalance)
			r.Delete("/balance/{id}", handlers.DeleteBalance)
			//trans
			r.Post("/transaction", handlers.CreateTransaction)
			r.Get("/transaction", handlers.GetAllTransaction)
			r.Delete("/transaction/{id}", handlers.DeleteTransaction)
			r.Put("/transaction", handlers.UpdateTransaction)
			//goal
			r.Post("/goal", handlers.CreateGoal)
			r.Delete("/goal/{id}", handlers.DeleteGoal)
			r.Put("/goal", handlers.UpdateGoal)
			r.Get("/goal", handlers.GetAllGoals)
			r.Get("/goal/remain", handlers.GetRemainGoals)
			//
			r.Get("/expenses", handlers.GetAllExpenses)

		})
	})

	return r
}
