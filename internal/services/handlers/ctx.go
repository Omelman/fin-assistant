package handlers

import (
	"context"
	"github.com/fin-assistant/internal/config"
	"github.com/fin-assistant/internal/data"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	usersCtxKey
	emailCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxUser(entry data.Users) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, usersCtxKey, entry)
	}
}

func User(r *http.Request) data.Users {
	return r.Context().Value(usersCtxKey).(data.Users)
}

func CtxEmail(entry config.Email) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, emailCtxKey, entry)
	}
}

func Email(r *http.Request) config.Email {
	return r.Context().Value(emailCtxKey).(config.Email)
}
