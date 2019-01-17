package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestQuote(t *testing.T) {
	quote := PrepareQuotes()
	if len(*quote) <= 10 {
		t.Errorf("Quotes was too short, got: %d, want at least %d.", len(*quote), 10)
	}
}

func TestQuoteApi(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/quote", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(QuoteResponse)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := true
	if strings.Contains(rr.Body.String(), "data") != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
