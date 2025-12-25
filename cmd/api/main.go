package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"loyalty-points-service/internal/handlers"

	"github.com/gorilla/mux"
)

// Metric globale : Compteur de requêtes
var requestCount uint64

// Structure pour les logs JSON
type LogEntry struct {
	Time      string `json:"time"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

func main() {
	r := mux.NewRouter()

	// 1. Ajout du Middleware d'Observabilité (Logs + Metrics + Tracing)
	r.Use(observabilityMiddleware)

	// Routes API
	api := r.PathPrefix("/api/v1").Subrouter()

	// Points endpoints
	api.HandleFunc("/customers/{customerID}/points", handlers.GetCustomerPoints).Methods("GET")
	api.HandleFunc("/customers/{customerID}/points/add", handlers.AddPoints).Methods("POST")
	api.HandleFunc("/customers/{customerID}/points/redeem", handlers.RedeemPoints).Methods("POST")

	// Health check
	api.HandleFunc("/health", healthCheck).Methods("GET")

	// 2. Endpoint Metrics (Nouvelle consigne)
	api.HandleFunc("/metrics", metricsHandler).Methods("GET")

	// Configuration du serveur avec Timeouts (Pour la sécurité/gosec)
	srv := &http.Server{
		Handler:      r,
		Addr:         ":8081",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logJSON("INFO", "🚀 Loyalty Points Service starting on :8081", "")
	log.Fatal(srv.ListenAndServe())
}

// Middleware : Gère le Tracing, les Logs et les Métriques pour chaque requête
func observabilityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// A. Tracing : Créer un ID unique pour la requête
		requestID := fmt.Sprintf("req-%d", time.Now().UnixNano())
		w.Header().Set("X-Request-ID", requestID)

		// B. Metrics : Incrémenter le compteur
		atomic.AddUint64(&requestCount, 1)

		// C. Logs Structurés (Début de requête)
		logJSON("INFO", fmt.Sprintf("Started %s %s", r.Method, r.URL.Path), requestID)

		// Exécuter la vraie requête
		next.ServeHTTP(w, r)
	})
}

// Fonction utilitaire pour écrire des logs en JSON
func logJSON(level, message, reqID string) {
	entry := LogEntry{
		Time:      time.Now().Format(time.RFC3339),
		Level:     level,
		Message:   message,
		RequestID: reqID,
	}
	json.NewEncoder(os.Stdout).Encode(entry)
}

// Handler pour afficher les métriques
func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// On récupère la valeur actuelle du compteur
	count := atomic.LoadUint64(&requestCount)

	response := map[string]interface{}{
		"total_requests": count,
		"status":         "up",
	}
	json.NewEncoder(w).Encode(response)
}

// Health check endpoint
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"status":  "healthy",
		"service": "loyalty-points",
		"version": "1.0.0",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Encoding error", http.StatusInternalServerError)
	}
}