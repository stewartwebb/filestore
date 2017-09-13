package common

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

// LoggingHandler is to use up disk space
func LoggingHandler(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}

// RecoveryHandler in case it crasho
func RecoveryHandler(next http.Handler) http.Handler {
	return handlers.RecoveryHandler()(next)
}
