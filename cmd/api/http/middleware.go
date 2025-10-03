package http

import (
	"log"
	"net/http"
	"os"
	"time"
)

// Package http define los controladores (handlers), wrappers de error y middleware
// para la capa de presentación HTTP.

// AuthAndLoggingMiddleware es un middleware que combina dos responsabilidades:
//  1. **Autenticación Básica (Basic Auth):** Verifica las credenciales Basic Auth
//     contra variables de entorno (BASIC_AUTH_USER y BASIC_AUTH_PASS).
//     Si las variables no están seteadas, la autenticación es omitida (Warning).
//  2. **Logging de Peticiones:** Registra el método HTTP, la URL, el protocolo y
//     el tiempo que tardó el procesamiento del request.
//
// Recibe un http.Handler (next) y retorna un nuevo http.Handler.
func AuthAndLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		userEnv := os.Getenv("BASIC_AUTH_USER")
		passEnv := os.Getenv("BASIC_AUTH_PASS")

		// --- Lógica de Autenticación Básica ---
		if userEnv == "" || passEnv == "" {
			log.Println("WARNING: BASIC_AUTH environment variables not set. Skipping authentication.")
		} else {
			user, pass, ok := r.BasicAuth()
			// Verifica que las credenciales coincidan con las variables de entorno.
			if !ok || user != userEnv || pass != passEnv {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return // Detiene el flujo si la autenticación falla.
			}
		}

		// Llama al siguiente handler/middleware en la cadena.
		next.ServeHTTP(w, r)

		// --- Lógica de Logging (ejecutada después de next.ServeHTTP) ---
		duration := time.Since(start)
		// Nota: http.StatusOK (200) se usa aquí para el log, pero no refleja
		// necesariamente el código final si el handler/wrapper retornó un error.
		// Para logging más preciso, se necesitaría un ResponseWriter personalizado.
		log.Printf(
			"[%s] %s %s | Status: %d | Duration: %v",
			r.Method,
			r.URL.Path,
			r.Proto,

			http.StatusOK, // Status code logueado (ver nota)
			duration,
		)
	})
}
