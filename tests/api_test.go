package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"loyalty-points-service/internal/handlers"

	"github.com/gorilla/mux"
)

// TestHealthCheck vérifie que l'endpoint health fonctionne
func TestHealthCheck(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/health", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthCheckHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Health check returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Vérifie le contenu de la réponse
	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Could not parse JSON response: %v", err)
	}

	expectedFields := []string{"status", "service", "version"}
	for _, field := range expectedFields {
		if _, exists := response[field]; !exists {
			t.Errorf("Response missing expected field: %s", field)
		}
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", response["status"])
	}
}

// TestGetCustomerPointsExisting vérifie la récupération des points d'un client existant
func TestGetCustomerPointsExisting(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/customers/cust-001/points", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GetCustomerPoints returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var customer map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &customer); err != nil {
		t.Errorf("Could not parse customer JSON: %v", err)
	}

	// Vérifie les champs obligatoires
	requiredFields := []string{"id", "points"}
	for _, field := range requiredFields {
		if _, exists := customer[field]; !exists {
			t.Errorf("Customer response missing field: %s", field)
		}
	}

	// Vérifie que les points sont un nombre positif
	if points, ok := customer["points"].(float64); ok {
		if points < 0 {
			t.Errorf("Customer points should be positive, got: %v", points)
		}
	}
}

// TestGetCustomerPointsNonExisting vérifie le comportement avec un client inexistant
func TestGetCustomerPointsNonExisting(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/customers/non-existing-cust/points", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(rr, req)

	// Doit retourner 404 pour un client inexistant
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("GetCustomerPoints for non-existing customer returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

// TestAddPoints vérifie l'ajout de points
func TestAddPoints(t *testing.T) {
	requestBody := map[string]interface{}{
		"points": 25,
		"reason": "Test purchase",
	}
	jsonData, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/api/v1/customers/cust-002/points/add", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("AddPoints returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Could not parse response JSON: %v", err)
	}

	// Vérifie le message de succès
	if message, exists := response["message"]; !exists || message != "Points added successfully" {
		t.Errorf("Expected success message, got: %v", response)
	}

	// Vérifie que le client est dans la réponse
	if customer, exists := response["customer"]; exists {
		if customerMap, ok := customer.(map[string]interface{}); ok {
			if points, exists := customerMap["points"]; exists {
				t.Logf("Customer now has %v points", points)
			}
		}
	}
}

// TestAddPointsInvalid vérifie la validation des données
func TestAddPointsInvalid(t *testing.T) {
	// Test avec points négatifs
	requestBody := map[string]interface{}{
		"points": -10, // Points négatifs doivent être rejetés
		"reason": "Invalid points",
	}
	jsonData, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/api/v1/customers/cust-001/points/add", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(rr, req)

	// Doit retourner 400 Bad Request pour points invalides
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("AddPoints with negative points returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

// TestRedeemPoints vérifie l'utilisation de points
func TestRedeemPoints(t *testing.T) {
	requestBody := map[string]interface{}{
		"points": 20,
		"reward": "Test reward",
	}
	jsonData, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/api/v1/customers/cust-001/points/redeem", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("RedeemPoints returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Could not parse response JSON: %v", err)
	}

	if message, exists := response["message"]; !exists || message != "Points redeemed successfully" {
		t.Errorf("Expected redemption success message, got: %v", response)
	}
}

// TestRedeemPointsInsufficient vérifie le cas où le client n'a pas assez de points
func TestRedeemPointsInsufficient(t *testing.T) {
	// Essaye d'utiliser un nombre énorme de points
	requestBody := map[string]interface{}{
		"points": 10000,
		"reward": "Expensive reward",
	}
	jsonData, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/api/v1/customers/cust-003/points/redeem", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(rr, req)

	// Doit retourner 400 pour points insuffisants
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("RedeemPoints with insufficient points returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

// TestIntegration simule un scénario complet
func TestIntegration(t *testing.T) {
	// 1. Vérifie les points initiaux
	req1, _ := http.NewRequest("GET", "/api/v1/customers/cust-001/points", nil)
	rr1 := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(rr1, req1)

	var initialCustomer map[string]interface{}
	json.Unmarshal(rr1.Body.Bytes(), &initialCustomer)
	initialPoints := initialCustomer["points"].(float64)

	// 2. Ajoute des points
	addBody, _ := json.Marshal(map[string]interface{}{"points": 30, "reason": "Integration test"})
	req2, _ := http.NewRequest("POST", "/api/v1/customers/cust-001/points/add", bytes.NewBuffer(addBody))
	req2.Header.Set("Content-Type", "application/json")
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	// 3. Vérifie les nouveaux points
	req3, _ := http.NewRequest("GET", "/api/v1/customers/cust-001/points", nil)
	rr3 := httptest.NewRecorder()
	router.ServeHTTP(rr3, req3)

	var finalCustomer map[string]interface{}
	json.Unmarshal(rr3.Body.Bytes(), &finalCustomer)
	finalPoints := finalCustomer["points"].(float64)

	// Vérifie que les points ont été ajoutés
	expectedPoints := initialPoints + 30
	if finalPoints != expectedPoints {
		t.Errorf("Points integration test failed: initial=%v, expected=%v, final=%v", initialPoints, expectedPoints, finalPoints)
	} else {
		t.Logf("✅ Integration test passed: %v + 30 = %v", initialPoints, finalPoints)
	}
}

// Fonctions utilitaires pour les tests
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "loyalty-points",
		"version": "1.0.0",
	})
}

func setupRouter() *mux.Router {
	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/customers/{customerID}/points", handlers.GetCustomerPoints).Methods("GET")
	api.HandleFunc("/customers/{customerID}/points/add", handlers.AddPoints).Methods("POST")
	api.HandleFunc("/customers/{customerID}/points/redeem", handlers.RedeemPoints).Methods("POST")
	api.HandleFunc("/health", healthCheckHandler).Methods("GET")
	return router
}
