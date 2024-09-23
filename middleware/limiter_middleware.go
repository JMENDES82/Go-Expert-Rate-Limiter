package middleware

import (
	"net"
	"net/http"
	"strings"

	"github.com/JMENDES82/Go-Expert-Rate-Limiter/limiter"
)

func RateLimiterMiddleware(l limiter.Limiter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			identifier := getIPFromRemoteAddr(r.RemoteAddr) // Extrai apenas o IP
			isToken := false

			token := r.Header.Get("API_KEY")
			if token != "" {
				identifier = token // Priorizar o token se existir
				isToken = true
			}

			if !l.AllowRequest(identifier, isToken) {
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Função para extrair apenas o IP do RemoteAddr (ignora a porta)
func getIPFromRemoteAddr(remoteAddr string) string {
	// Remove a parte da porta se existir
	if strings.Contains(remoteAddr, ":") {
		ip, _, err := net.SplitHostPort(remoteAddr)
		if err == nil {
			return ip
		}
	}
	return remoteAddr // Caso não tenha porta, retorna o endereço completo (em IPv6, por exemplo)
}
