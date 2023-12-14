package main

import (
	"context"
	"net/http"
	"os"

	"github.com/flosch/pongo2/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
	"github.com/webdevfuel/projectmotor/auth"
	"github.com/webdevfuel/projectmotor/handler"
	"github.com/webdevfuel/projectmotor/template"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

func main() {
	h := handler.NewHandler(handler.HandlerOptions{
		Store: store,
	})
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	fs := http.FileServer(http.Dir("./dist"))
	r.Handle("/dist/*", http.StripPrefix("/dist/", fs))
	r.Get("/login", h.Login)
	// Group protected routes into one function to run middleware
	r.Group(protectedRouter(h))
	http.ListenAndServe("localhost:3000", r)
}

// Router with user ensured
//
// Add routes here where user has to be logged in
func protectedRouter(h *handler.Handler) func(chi.Router) {
	return func(r chi.Router) {
		r.Use(ProtectedCtx(h))
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			template.Dashboard.ExecuteWriter(pongo2.Context{}, w)
		})
	}
}

// Protected context
//
// Middleware checks if user exists within current session
func ProtectedCtx(h *handler.Handler) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			session, err := h.GetSessionStore(r)
			// redirect in case of error
			if err != nil {
				redirectToLogin(w, r)
				return
			}
			user := session.Values["user"]
			// redirect in case of missing user
			if user == nil {
				redirectToLogin(w, r)
				return
			}
			// TODO: check if user is db.User
			ctx := r.Context()
			ctx = context.WithValue(ctx, auth.UserIDKey{}, 0)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

// Redirect to public auth route
//
// Use this when session user doesn't exist
func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://localhost:3000/login", http.StatusSeeOther)
}
