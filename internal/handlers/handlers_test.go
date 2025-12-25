package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// TestGetCustomerPointsHandler teste directement le handler
func TestGetCustomerPointsHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/customers/cust-001/points", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetCustomerPoints)

	// Simule les variables mux
	req = mux.SetURLVars(req, map[string]string{"customerID": "cust-001"})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

// TestAddPointsHandler teste le handler d'ajout de points
func TestAddPointsHandler(t *testing.T) {
	requestBody := map[string]interface{}{
		"points": 15,
		"reason": "Handler test",
	}
	jsonData, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/customers/cust-001/points/add", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddPoints)

	// Simule les variables mux
	req = mux.SetURLVars(req, map[string]string{"customerID": "cust-001"})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("AddPoints handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
