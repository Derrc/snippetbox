package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

// sets HTTP security headers (inline with OWASP guidance)
func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
		"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		w.Header().Set("Server", "Go")

		next.ServeHTTP(w, r)
	})
}

// logs information for received requests
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip = r.RemoteAddr
			proto = r.Proto
			method = r.Method
			uri = r.URL.RequestURI()
		)

		app.logger.Info("received request", "ip", ip, "proto", proto, "method", method, "uri", uri)

		next.ServeHTTP(w, r)
	})
}

// gracefully shuts down in the event of a panic by setting the 'Connection' header
// and writing an error response before closing the underyling HTTP connection for the
// affected goroutine
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create a deferred function (always will be run in the event of a panic)
		// as Go unwinds stack
		defer func() {
			if err := recover(); err != nil {
				// header triggers Go's HTTP server to automatically close current connection after response is sent
				w.Header().Set("Connection", "close")
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// redirects user to login page for protected routes
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if user is not authenticated, redirect them to the login page
		// don't allow subsequent handlers to execute in middleware chain
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// set 'Cache-Control: no-store' header so pages that require authentication
		// are not stored in user's local browser cache or other shared cache (i.e. proxy)
		w.Header().Set("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

// checks if session data contains 'authenticatedUserID' and if so
// adds (isAuthenticatedContextKey, true) to the request context
// for all future middlewares/handlers
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}
		
		// make sure that user id exists
		exists, err := app.users.Exists(id)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		if exists {
			// reassigns request context, adding (isAuthenticatedContextKey, true)
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

// prevents CSRF attacks by using the double-submit cookie pattern
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		// can't be modified by javascript client-side
		HttpOnly: true,
		Path: "/",
		// only sent to to server with an encrypted request over HTTPS
		Secure: true,
	})

	return csrfHandler
}