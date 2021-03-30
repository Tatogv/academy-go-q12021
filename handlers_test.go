package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"api/handlers"

	"github.com/gorilla/mux"
)

func TestGetAll(t *testing.T) {
	req, err := http.NewRequest("GET", "/read", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.GetAll)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"1":"bulbasaur","2":"ivysaur","3":"venusaur"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetById(t *testing.T) {
	req, err := http.NewRequest("GET", "/read/2", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	vars := map[string]string{
		"pokemonId": "2",
	}
	req = mux.SetURLVars(req, vars)

	handler := http.HandlerFunc(handlers.GetById)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"id":"2","name":"ivysaur"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetBerries(t *testing.T) {
	os.Remove(handlers.FilePath)
	req, err := http.NewRequest("GET", "/getBerries", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.GetBerries)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	_, err = os.Open(handlers.FilePath)
	if err != nil {
		t.Errorf("File not saved: error %v",
			err)
	}

}
