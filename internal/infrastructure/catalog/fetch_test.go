package catalog

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetch_OK(t *testing.T) {
	want := []byte(`{"schema_version":1,"packages":[]}`)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(want)
	}))
	defer srv.Close()

	ctx := context.Background()
	body, err := Fetch(ctx, srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != string(want) {
		t.Fatalf("body = %q", body)
	}
}

func TestFetch_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	_, err := Fetch(context.Background(), srv.URL)
	if err == nil {
		t.Fatal("expected error")
	}
}
