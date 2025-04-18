package http

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestEnvHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/env", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	srv := NewMockServer()
	// bind handler to mock server
	handler := http.HandlerFunc(srv.infoHandler)
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned bad status code: got %v, want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := ".*hostname.*"
	r := regexp.MustCompile(expected)
	if !r.MatchString(rr.Body.String()) {
		t.Fatalf("handler returned unexpected body:\ngot \n%v \nwant \n%s",
			rr.Body.String(), expected)
	}
}
