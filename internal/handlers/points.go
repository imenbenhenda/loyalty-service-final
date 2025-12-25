package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"loyalty-points-service/internal/models"

	"github.com/gorilla/mux"
)

var (
	// Stockage en mémoire (simplifié pour le projet)
	customers    = make(map[string]*models.Customer)
	transactions []models.Transaction
	mutex        sync.RWMutex
)

// Initialiser quelques données de test
func init() {
	customers["cust-001"] = &models.Customer{ID: "cust-001", Points: 150, Name: "John Doe", Email: "john@email.com"}
	customers["cust-002"] = &models.Customer{ID: "cust-002", Points: 75, Name: "Jane Smith", Email: "jane@email.com"}
	customers["cust-003"] = &models.Customer{ID: "cust-003", Points: 200, Name: "Bob Wilson", Email: "bob@email.com"}
}

// GetCustomerPoints retourne les points d'un client
func GetCustomerPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID := vars["customerID"]

	mutex.RLock()
	defer mutex.RUnlock()

	customer, exists := customers[customerID]
	if !exists {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

// AddPoints ajoute des points à un client
func AddPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID := vars["customerID"]

	var request models.AddPointsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.Points <= 0 {
		http.Error(w, "Points must be positive", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	customer, exists := customers[customerID]
	if !exists {
		// Créer un nouveau client s'il n'existe pas
		customer = &models.Customer{ID: customerID, Points: 0}
		customers[customerID] = customer
	}

	// Ajouter les points
	customer.Points += request.Points

	// Enregistrer la transaction
	transaction := models.Transaction{
		ID:         generateID(),
		CustomerID: customerID,
		Type:       "earn",
		Points:     request.Points,
		Reason:     request.Reason,
		Timestamp:  time.Now(),
	}
	transactions = append(transactions, transaction)

	// Réponse
	response := map[string]interface{}{
		"message":     "Points added successfully",
		"customer":    customer,
		"transaction": transaction,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// RedeemPoints utilise des points pour une récompense
func RedeemPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID := vars["customerID"]

	var request models.RedeemPointsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.Points <= 0 {
		http.Error(w, "Points must be positive", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	customer, exists := customers[customerID]
	if !exists {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	// Vérifier si le client a assez de points
	if customer.Points < request.Points {
		http.Error(w, "Insufficient points", http.StatusBadRequest)
		return
	}

	// Déduire les points
	customer.Points -= request.Points

	// Enregistrer la transaction
	transaction := models.Transaction{
		ID:         generateID(),
		CustomerID: customerID,
		Type:       "redeem",
		Points:     request.Points,
		Reason:     request.Reward,
		Timestamp:  time.Now(),
	}
	transactions = append(transactions, transaction)

	// Réponse
	response := map[string]interface{}{
		"message":     "Points redeemed successfully",
		"reward":      request.Reward,
		"customer":    customer,
		"transaction": transaction,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Fonction utilitaire pour générer un ID simple
func generateID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
