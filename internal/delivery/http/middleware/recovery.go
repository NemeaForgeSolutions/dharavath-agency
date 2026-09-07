package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC in HTTP handler: %v\nStack trace:\n%s", err, debug.Stack())
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `
					<!DOCTYPE html>
					<html>
					<head><title>500 - Internal Server Error</title></head>
					<body style="font-family: sans-serif; text-align: center; padding: 4rem; background: #071522; color: #fff;">
						<h1 style="color: #e1b34c;">500 - Internal Server Error</h1>
						<p>Something unexpected occurred. Our engineering desk has been notified.</p>
						<a href="/" style="color: #38bdf8; text-decoration: underline;">Return to Home</a>
					</body>
					</html>
				`)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
