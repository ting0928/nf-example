package sbi

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// Todo represents a single todo item.
type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// TodoStore keeps todo items in memory.
type TodoStore struct {
	mu     sync.Mutex
	nextID int
	items  []Todo
}

var todoStore = &TodoStore{
	nextID: 1,
	items:  []Todo{},
}

// Register GET and POST routes.
func (s *Server) getTodoRoute() []Route {
	return []Route{
		{
			Name:    "Get Todos",
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: getTodos,
		},
		{
			Name:    "Create Todo",
			Method:  http.MethodPost,
			Pattern: "/",
			APIFunc: createTodo,
		},
	}
}

// GET /todo/
func getTodos(c *gin.Context) {
	todoStore.mu.Lock()
	items := append([]Todo{}, todoStore.items...)
	todoStore.mu.Unlock()

	c.JSON(http.StatusOK, items)
}

// POST /todo/
func createTodo(c *gin.Context) {
	var request struct {
		Title string `json:"title"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON request",
		})
		return
	}

	title := strings.TrimSpace(request.Title)
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Title cannot be empty",
		})
		return
	}

	todoStore.mu.Lock()

	todo := Todo{
		ID:    todoStore.nextID,
		Title: title,
	}

	todoStore.nextID++
	todoStore.items = append(todoStore.items, todo)

	todoStore.mu.Unlock()

	c.JSON(http.StatusCreated, todo)
}
