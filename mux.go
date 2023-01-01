package main

// import "net/http"
import (
	"context"
	"net/http"
	
	"github.com/mongolmongol2022/go_todo_app_HandsOnTraining/auth"
	"github.com/mongolmongol2022/go_todo_app_HandsOnTraining/service"

	"github.com/mongolmongol2022/go_todo_app_HandsOnTraining/clock"
	"github.com/mongolmongol2022/go_todo_app_HandsOnTraining/config"
	"github.com/mongolmongol2022/go_todo_app_HandsOnTraining/handler"
	"github.com/mongolmongol2022/go_todo_app_HandsOnTraining/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func NewMux(ctx context.Context, cfg *config.Config) (http.Handler, func(), error) {
	mux := chi.NewRouter()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	v := validator.New()

	db, cleanup, err := store.New(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	clocker := clock.RealClocker{}
	r := store.Repository{Clocker: clocker}
	rcli, err := store.NewKVS(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	jwter, err := auth.NewJWTer(rcli, clocker)
	if err != nil {
		return nil, cleanup, err
	}

	// r := store.Repository{Clocker: clock.RealClocker{}}
	// at := &handler.AddTask{
	// 	Service:   &service.AddTask{DB: db, Repo: &r},
	// 	Validator: v,
	// }
	// mux.Post("/tasks", at.ServeHTTP)
	// lt := &handler.ListTask{
	// 	Service: &service.ListTask{DB: db, Repo: &r},
	// }
	// mux.Get("/tasks", lt.ServeHTTP)

	ru := &handler.RegisterUser{
		Service:   &service.RegisterUser{DB: db, Repo: &r},
		Validator: v,
	}
	mux.Post("/register", ru.ServeHTTP)

	l := &handler.Login{
		Service: &service.Login{
			DB:             db,
			Repo:           &r,
			TokenGenerator: jwter,
		},
		Validator: v,
	}
	mux.Post("/login", l.ServeHTTP)

	at := &handler.AddTask{
		Service:   &service.AddTask{DB: db, Repo: &r},
		Validator: v,
	}
	lt := &handler.ListTask{
		Service: &service.ListTask{DB: db, Repo: &r},
	}
	mux.Route("/tasks", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwter))
		r.Post("/", at.ServeHTTP)
		r.Get("/", lt.ServeHTTP)
	})
	mux.Route("/admin", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwter), handler.AdminMiddleware)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			_, _ = w.Write([]byte(`{"message": "admin only"}`))
		})
	})

	return mux, cleanup, nil
	
}
