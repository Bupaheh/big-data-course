package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReplace(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/replace", nil)

	replace(rr, req)

	resultCode := rr.Result().StatusCode

	if resultCode != http.StatusOK {
		t.Errorf("Incorrect status code, got: %d, want: %d", resultCode, http.StatusOK)
	}
}

func TestGet(t *testing.T) {
	body := "TestGet"
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/replace", strings.NewReader(body))

	replace(rr, req)

	rr = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/get", nil)

	get(rr, req)
	resultCode := rr.Result().StatusCode
	var buffer [1000]byte
	n, _ := rr.Result().Body.Read(buffer[:])
	receivedBody := string(buffer[:n])

	if resultCode != http.StatusOK {
		t.Errorf("Incorrect status code, got: %d, want: %d", resultCode, http.StatusOK)
	}

	if body != receivedBody {
		t.Errorf("Incorrect body code, got: %s, want: %s", receivedBody, body)
	}
}
