package sbi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestTodoAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Use a fresh store for testing.
	originalStore := todoStore
	todoStore = &TodoStore{
		nextID: 1,
		items:  []Todo{},
	}
	defer func() {
		todoStore = originalStore
	}()

	router := gin.New()
	todoGroup := router.Group("/todo")
	applyRoutes(todoGroup, (&Server{}).getTodoRoute())

	t.Run("GET empty todo list", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/todo/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}

		if strings.TrimSpace(w.Body.String()) != "[]" {
			t.Fatalf("unexpected response: %s", w.Body.String())
		}
	})

	t.Run("POST create todo", func(t *testing.T) {
		body := strings.NewReader(`{"title":"Learn Git"}`)
		req := httptest.NewRequest(http.MethodPost, "/todo/", body)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d", w.Code)
		}

		if !strings.Contains(w.Body.String(), `"title":"Learn Git"`) {
			t.Fatalf("unexpected response: %s", w.Body.String())
		}
	})

	t.Run("GET created todo", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/todo/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}

		if !strings.Contains(w.Body.String(), `"title":"Learn Git"`) {
			t.Fatalf("todo not found: %s", w.Body.String())
		}
	})

	t.Run("POST empty title", func(t *testing.T) {
		body := strings.NewReader(`{"title":""}`)
		req := httptest.NewRequest(http.MethodPost, "/todo/", body)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})
}
