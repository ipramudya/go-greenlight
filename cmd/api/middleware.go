package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ipramudya/go-greenlight/internal/data"
	"github.com/ipramudya/go-greenlight/internal/validator"
	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				rw.Header().Set("Connection", "close")
				app.serverErrorResponse(rw, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(rw, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	type client struct {
		lastSeen time.Time
		limitter *rate.Limiter
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// launch a background goroutine which removes old entries
	// from the clients map once every minute
	go func() {
		for {
			time.Sleep(time.Minute)

			// lock mutex to prevent any rate limiter checks from happening
			// while the cleanup is taking place
			mu.Lock()

			// if any client haven't been seen within the last three minutes,
			// delete the corresponding entry from map
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// extract client's IP address from the request
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			app.serverErrorResponse(rw, r, err)
			return
		}

		// prevent this code to run cuncurrently
		mu.Lock()

		// initialize a new rate limiter and add the IP address and
		// limitter to the map if it doesn't already exist
		if _, found := clients[ip]; !found {
			clients[ip] = &client{limitter: rate.NewLimiter(rate.Limit(app.limiter.rps), app.limiter.burst)}
		}

		// call Allow limitter method for the current IP address,
		// if the request isn't allowed, unlock the mutex and send 429 response
		if !clients[ip].limitter.Allow() {
			mu.Unlock()
			app.rateLimitExceedResponse(rw, r)
			return
		}

		mu.Unlock()
		next.ServeHTTP(rw, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Vary", "Authorization")

		authorizationH := r.Header.Get("Authorization")

		// when there's no authorization header found, add the AnonymousUser to the request context
		if authorizationH == "" {
			r = app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(rw, r)
			return
		}

		headerParts := strings.Split(authorizationH, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(rw, r)
			return
		}

		token := headerParts[1]

		v := validator.New()

		if data.ValidateTokenPlaintext(v, token); !v.IsValid() {
			app.invalidAuthenticationTokenResponse(rw, r)
			return
		}

		user, err := app.models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(rw, r)
			default:
				app.serverErrorResponse(rw, r, err)
			}
			return
		}

		r = app.contextSetUser(r, user)
		next.ServeHTTP(rw, r)
	})
}
