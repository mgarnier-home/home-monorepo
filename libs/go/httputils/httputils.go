package httputils

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"mgarnier11.fr/go/libs/logger"
)

func LogRequestMiddleware(logger *logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Infof("%s %s Request", r.Method, r.URL.Path)

			next.ServeHTTP(w, r)
		})
	}
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace * with specific origin(s) if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Api-Token, Authorization")

		// Handle preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CheckApiTokenMiddleware(authorizedToken string, header string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiToken := r.Header.Get(header)

			if apiToken != authorizedToken {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func GetRequest(url string) error {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // WARNING: disables all certificate checks!
	}
	client := &http.Client{
		Transport: transport,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get %s: %s", url, resp.Status)
	}
	logger.Infof("Successfully got %s", url)
	return nil
}
