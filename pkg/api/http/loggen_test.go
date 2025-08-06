package http

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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
		t.Errorf("handler returned bad status code: got %v want %v",
			status, http.StatusOK)
	}
	t.Log(rr.Body.String())
}

func TestLogGenHandlerUsesParams(t *testing.T) {
	req, err := http.NewRequest("GET", "/loggen?per_second=4000&burst_dur=1&message_size=64", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	srv := NewMockServer()
	handler := http.HandlerFunc(srv.logGenHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned bad status code: got %v want %v",
			status, http.StatusOK)
	}
	t.Log(rr.Body.String())
}

func TestLogGenHandlerAppendsToExistingFile(t *testing.T) {
	// Create a temporary file for logging output
	tmpfile, err := os.CreateTemp("", "test.log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Write some initial content to the file
	initialContent := "initial content\n"
	if _, err := tmpfile.Write([]byte(initialContent)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Setup the server with the temporary file as output
	srv := NewMockServer()
	srv.config.LogwildOutFile = tmpfile.Name()

	// First request to the handler
	req1, err := http.NewRequest("GET", "/loggen?per_second=1&burst_dur=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr1 := httptest.NewRecorder()
	handler := http.HandlerFunc(srv.logGenHandler)
	handler.ServeHTTP(rr1, req1)

	// Check the status code
	if status := rr1.Code; status != http.StatusOK {
		t.Errorf("handler returned bad status code: got %v want %v", status, http.StatusOK)
	}

	// Second request to the handler
	req2, err := http.NewRequest("GET", "/loggen?per_second=1&burst_dur=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)

	if status := rr2.Code; status != http.StatusOK {
		t.Errorf("handler returned bad status code: got %v want %v", status, http.StatusOK)
	}

	// Read the file content
	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Assert that the file still has the lines we added before involving LogMaker
	lines := strings.Split(string(content), "\n")
	if !strings.Contains(string(content), initialContent) {
		t.Errorf("expected file to contain initial content, but it does not")
	}
	if len(lines) <= 2 { // expecting at least 2 lines of logs + initial content
		t.Errorf("expected file to have more lines, but it has %d", len(lines))
	}
}

