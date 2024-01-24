package main

import (
	"fmt"
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com") // defines where we can load assets on this server from
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin") // full URL on same origin requests, when cross origin, URL path and any query values are stripped
		w.Header().Set("X-Content-Type-Options", "nosniff")           // prevents content sniffing attacks
		w.Header().Set("X-Frame-Options", "deny")                     // prevents clickjacking attacks
		w.Header().Set("X-XSS-Protection", "0")                       // disable XSS filter in the browser (recommended by OWASP standard)

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = ReadUserIP(r)
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		app.logger.Info("Received request", "ip", ip, "protocol", proto, "method", method, "uri", uri)

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event
		// of a panic as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a
			// panic or not. If there has...
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
