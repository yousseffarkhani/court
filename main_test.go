package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	hf := http.HandlerFunc(handler)
	hf.ServeHTTP(res, req)

	assertStatusOK(t, res.Code)

	got := res.Body.String()
	want := "hello world"
	if got != want {
		t.Errorf("got : %s, want : %s", got, want)
	}
}

func TestBirdHandler(t *testing.T) {
	t.Run("POST '/birds' adds bird to list of birds", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "", nil)
		if err != nil {
			t.Fatal(err)
		}

		res := httptest.NewRecorder()
		hf := http.HandlerFunc(createBirdHandler)

		hf.ServeHTTP(res, req)

		assertStatusOK(t, res.Code)

	})
	t.Run("POST '/birds' should redirect to index.html", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "", nil)
		if err != nil {
			t.Fatal(err)
		}

		res := httptest.NewRecorder()
		hf := http.HandlerFunc(createBirdHandler)

		hf.ServeHTTP(res, req)
		// assertStatusOK(t, re)
	})
	t.Run("GET '/birds' should return a list of birds", func(t *testing.T) {
		birds := []Bird{
			{"Sparrow", "A small harmless bird"},
		}
		req, err := http.NewRequest(http.MethodGet, "", nil)
		res := httptest.NewRecorder()

		hf := http.HandlerFunc(getBirdHandler)

		hf.ServeHTTP(res, req)

		assertStatusOK(t, res.Code)

		want := birds
		var got []Bird
		err = json.NewDecoder(res.Body).Decode(&got)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(want, got) {
			t.Errorf("got : %s, want : %s", got, want)
		}
	})

}

func TestRouter(t *testing.T) {
	r := newRouter()
	mockServer := httptest.NewServer(r)
	t.Run("Get '/hello' route should return OK", func(t *testing.T) {
		res, err := http.Get(mockServer.URL + "/hello")
		if err != nil {
			t.Fatal(err)
		}

		assertStatusOK(t, res.StatusCode)
	})
	t.Run("Get '/XXXX' route should return 404", func(t *testing.T) {
		res, err := http.Get(mockServer.URL + "/helloo")
		if err != nil {
			t.Fatal(err)
		}

		assertStatusNotFound(t, res.StatusCode)
	})
	t.Run("POST '/hello' route should be forbidden", func(t *testing.T) {
		res, err := http.Post(mockServer.URL+"/hello", "", nil)
		if err != nil {
			t.Fatal(err)
		}

		assertStatusMethodNotAllowed(t, res.StatusCode)
	})
}

func TestStaticFileServer(t *testing.T) {
	r := newRouter()
	mockServer := httptest.NewServer(r)

	res, err := http.Get(mockServer.URL + "/assets/")
	if err != nil {
		t.Fatal(err)
	}

	assertStatusOK(t, res.StatusCode)

	got := res.Header.Get("Content-Type")
	want := "text/html; charset=utf-8"
	if got != want {
		t.Errorf("got : %s, want : %s", got, want)
	}

}

func assertStatusOK(t *testing.T, status int) {
	t.Helper()
	if status != http.StatusOK {
		t.Error("Expected status OK, got ", status)
	}
}

func assertStatusNotFound(t *testing.T, status int) {
	t.Helper()
	if status != http.StatusNotFound {
		t.Error("Expected status not found, got ", status)
	}
}

func assertStatusMethodNotAllowed(t *testing.T, status int) {
	t.Helper()
	if status != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d\n", http.StatusMethodNotAllowed, status)
	}
}
