package handlers

import "net/http"

// CORSMiddleware enables Cross-Origin Resource Sharing for our API.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set headers to allow requests from any origin, with any method and headers.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// If this is a pre-flight "OPTIONS" request, we just send back the headers and a 200 OK.
		// The browser sends this automatically to check if the actual request is safe to send.
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}
