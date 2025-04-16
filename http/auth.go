package http

import (
	"csm/auth"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

const (
	CSM_CLIENT_ID     = "test_id"
	CSM_CLIENT_SECRET = "test_secret"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		prefix := "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			log.Printf("Authorization header not prefixed with bearer: %s", authHeader)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token := authHeader[len(prefix):]
		valid, err := auth.IsJwTValid(token)
		if err != nil || !valid {
			log.Printf("Access token not valid: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func handleGenerateJWT(w http.ResponseWriter, r *http.Request) {
	clientID, clientSecret, err := auth.DecodeBasicAuth(r.Header.Get("Authorization"))
	if clientID != CSM_CLIENT_ID || clientSecret != CSM_CLIENT_SECRET || err != nil {
		log.Printf("Generate jwt failed. ClientID-recieved: %s, ClientSecret-recieved: %s, err: %v. Expected clientID:%s and clientSecret: %s", clientID, clientSecret, err, CSM_CLIENT_ID, CSM_CLIENT_SECRET)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateJWT(clientID)
	if err != nil {
		log.Printf("Failed to generate jwt. Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Issuing jwt:", token)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
