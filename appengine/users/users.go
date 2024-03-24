package users

import (
	"context"
	"net/http"

	"google.golang.org/appengine/v2"
	"google.golang.org/appengine/v2/user"
)

var currentUser *User

// User ...
type User struct {
	email   string
	id      string
	isAdmin bool
}

// Email ...
func (u *User) Email() string { return u.email }

// ID ...
func (u *User) ID() string { return u.id }

// IsAdmin ...
func (u *User) IsAdmin() bool { return u.isAdmin }

func setCurrentUser(email, id string, isAdmin bool) {
	currentUser = &User{
		email:   email,
		id:      id,
		isAdmin: isAdmin,
	}
}

func resetCurrentUser() {
	currentUser = nil
}

// AddDevHandlers ...
func AddDevHandlers(mux *http.ServeMux) {
	mux.HandleFunc("GET /_login_", func(w http.ResponseWriter, req *http.Request) {
		d := req.URL.Query().Get("d")
		setCurrentUser("abcd@email.dummy", "_local_id_", true)
		http.Redirect(w, req, d, http.StatusSeeOther)
	})

	mux.HandleFunc("GET /_logout_", func(w http.ResponseWriter, req *http.Request) {
		d := req.URL.Query().Get("d")
		resetCurrentUser()
		http.Redirect(w, req, d, http.StatusSeeOther)
	})
}

// Current ...
func Current(ctx context.Context) *User {
	if appengine.IsAppEngine() {
		if u := user.Current(ctx); u != nil {
			return &User{
				email:   u.Email,
				id:      u.ID,
				isAdmin: u.Admin,
			}
		}
		return nil
	}
	return currentUser
}

// LoginURL ...
func LoginURL(ctx context.Context, dest string) (string, error) {
	if appengine.IsAppEngine() {
		return user.LoginURL(ctx, dest)
	}
	return "/_login_?d=" + dest, nil
}

// LogoutURL ...
func LogoutURL(ctx context.Context, dest string) (string, error) {
	if appengine.IsAppEngine() {
		return user.LogoutURL(ctx, dest)
	}
	return "/_logout_?d=" + dest, nil
}
