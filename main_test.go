package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReturnById(t *testing.T) {
	req, err := http.NewRequest("GET", "/articles/5fa84a48ce99b037713c6b9f", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ReturnById)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code, got %v expected %v", status, http.StatusOK)
	}

	expected := `{"_id":"5fa84a48ce99b037713c6b9f","Title":"sup","Subtitle":"sup sub","Content":"suppp content"}`
	if rr.Body.String() == expected {
		t.Errorf("handler returned unexpected body, got %v expected %v", rr.Body.String(), expected)
	}
}
