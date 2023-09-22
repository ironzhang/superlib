package httputils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponseRecorder(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponseRecorder(w)

	r.WriteHeader(http.StatusBadRequest)
	r.Header().Set("Content-Type", "text")
	fmt.Fprintf(r, "hello, world")

	if got, want := w.Code, http.StatusBadRequest; got != want {
		t.Errorf("response writer's code: got %v, want %v", got, want)
	} else {
		t.Logf("response writer's code: got %v", got)
	}
	if got, want := r.Status(), http.StatusBadRequest; got != want {
		t.Errorf("response recoder's status: got %v, want %v", got, want)
	} else {
		t.Logf("response recoder's status: got %v", got)
	}
	if got, want := string(w.Body.Bytes()), "hello, world"; got != want {
		t.Errorf("response writer's body: got %v, want %v", got, want)
	} else {
		t.Logf("response writer's body: got %v", got)
	}
	if got, want := string(r.Body()), "hello, world"; got != want {
		t.Errorf("response recorder's body: got %v, want %v", got, want)
	} else {
		t.Logf("response recorder's body: got %v", got)
	}
	if got, want := w.HeaderMap.Get("Content-Type"), "text"; got != want {
		t.Errorf("response writer's Content-Type header: got %v, want %v", got, want)
	} else {
		t.Logf("response writer's Content-Type header: got %v", got)
	}
}
