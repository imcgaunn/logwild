package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogGenHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/loggen", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	srv := NewMockServer()
	handler := http.HandlerFunc(srv.logGenHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
