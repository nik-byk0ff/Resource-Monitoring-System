package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddMetricInvalidPayload(t *testing.T) {
	req, err := http.NewRequest("POST", "/metrics", bytes.NewBuffer([]byte(`{"invalid":"data"}`)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddMetric)

	handler.ServeHTTP(rr, req)

	// Since NodeName is required, validation should fail with 400 Bad Request
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}
