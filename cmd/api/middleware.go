package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

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
