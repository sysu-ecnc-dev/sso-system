package handler

import "net/http"

func healthCheck(w http.ResponseWriter, r *http.Request) {
	successResponse(w, "Get the health status successfully.", map[string]string{
		"status": "healthy",
	})
}
