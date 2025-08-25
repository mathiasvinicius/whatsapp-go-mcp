package cmd

import (
	"net/http"
)

// corsMiddleware adds CORS headers required by Smithery.ai
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers as required by Smithery
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, mcp-session-id, mcp-protocol-version, *")
		w.Header().Set("Access-Control-Expose-Headers", "mcp-session-id, mcp-protocol-version")
		
		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		// Pass to the next handler
		next.ServeHTTP(w, r)
	})
}