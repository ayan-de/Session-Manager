package api_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/session-manager/backend/internal/api"
)

func TestClaudeCodeSessionsHandlerReturnsServerErrorWhenProviderFails(t *testing.T) {
	handler := api.NewClaudeCodeSessionsHandler(func() (any, error) {
		return nil, errors.New("boom")
	})

	req := httptest.NewRequest(http.MethodGet, "/api/claude-code/sessions", nil)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusInternalServerError)
	}
}

func TestClaudeCodeSessionsHandlerReturnsMethodNotAllowed(t *testing.T) {
	handler := api.NewClaudeCodeSessionsHandler(func() (any, error) {
		return []any{}, nil
	})

	req := httptest.NewRequest(http.MethodPost, "/api/claude-code/sessions", nil)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

func TestClaudeCodeSessionsHandlerReturnsOKOnSuccess(t *testing.T) {
	handler := api.NewClaudeCodeSessionsHandler(func() (any, error) {
		return []map[string]any{}, nil
	})

	req := httptest.NewRequest(http.MethodGet, "/api/claude-code/sessions", nil)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}

	if res.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("Content-Type = %q, want %q", res.Header().Get("Content-Type"), "application/json")
	}
}
