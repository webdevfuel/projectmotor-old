package handler

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type Handler struct {
	store *sessions.CookieStore
}

type HandlerOptions struct {
	Store *sessions.CookieStore
}

func NewHandler(options HandlerOptions) *Handler {
	return &Handler{
		store: options.Store,
	}
}

func (h Handler) GetSessionStore(r *http.Request) (*sessions.Session, error) {
	return h.store.Get(r, "_projectmotor_session")
}
