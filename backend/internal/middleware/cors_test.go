package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORS_DefaultOrigin(t *testing.T) {
	t.Setenv("CORS_ORIGIN", "")

	handler := CORS()(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	origin := rec.Header().Get("Access-Control-Allow-Origin")
	if origin != defaultCORSOrigin {
		t.Errorf("Allow-Origin: got %q, want %q", origin, defaultCORSOrigin)
	}
}

func TestCORS_CustomOrigin(t *testing.T) {
	t.Setenv("CORS_ORIGIN", "http://example.com")

	handler := CORS()(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	origin := rec.Header().Get("Access-Control-Allow-Origin")
	if origin != "http://example.com" {
		t.Errorf("Allow-Origin: got %q, want %q", origin, "http://example.com")
	}
}

func TestCORS_PreflightOptions(t *testing.T) {
	t.Setenv("CORS_ORIGIN", "")

	handlerCalled := false
	handler := CORS()(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusNoContent)
	}
	if handlerCalled {
		t.Error("next handler should not be called for preflight OPTIONS")
	}

	methods := rec.Header().Get("Access-Control-Allow-Methods")
	if methods == "" {
		t.Error("missing Access-Control-Allow-Methods header")
	}
	headers := rec.Header().Get("Access-Control-Allow-Headers")
	if headers == "" {
		t.Error("missing Access-Control-Allow-Headers header")
	}
}

func TestCORS_PassesThroughNonOptions(t *testing.T) {
	t.Setenv("CORS_ORIGIN", "")

	handlerCalled := false
	handler := CORS()(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if !handlerCalled {
		t.Error("next handler should be called for non-OPTIONS request")
	}
	if rec.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusOK)
	}

	origin := rec.Header().Get("Access-Control-Allow-Origin")
	if origin != defaultCORSOrigin {
		t.Errorf("Allow-Origin on non-OPTIONS: got %q, want %q", origin, defaultCORSOrigin)
	}
}
