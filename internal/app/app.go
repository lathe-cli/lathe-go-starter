package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type task struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type server struct {
	mu     sync.Mutex
	tasks  map[string]task
	nextID int
}

func NewHandler() http.Handler {
	s := &server{tasks: map[string]task{}, nextID: 1}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("GET /tasks", s.listTasks)
	mux.HandleFunc("POST /tasks", s.createTask)
	mux.HandleFunc("GET /tasks/{id}", s.getTask)
	mux.HandleFunc("PATCH /tasks/{id}", s.updateTask)
	mux.HandleFunc("DELETE /tasks/{id}", s.deleteTask)
	return mux
}

func (s *server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *server) listTasks(w http.ResponseWriter, _ *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tasks := make([]task, 0, len(s.tasks))
	for i := 1; i < s.nextID; i++ {
		if item, ok := s.tasks[strconv.Itoa(i)]; ok {
			tasks = append(tasks, item)
		}
	}
	writeJSON(w, http.StatusOK, tasks)
}

func (s *server) createTask(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title string `json:"title"`
	}
	if err := readJSON(w, r, &body); err != nil {
		writeError(w, err)
		return
	}
	body.Title = strings.TrimSpace(body.Title)
	if body.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "title is required"})
		return
	}

	s.mu.Lock()
	item := task{ID: strconv.Itoa(s.nextID), Title: body.Title}
	s.nextID++
	s.tasks[item.ID] = item
	s.mu.Unlock()
	writeJSON(w, http.StatusCreated, item)
}

func (s *server) getTask(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	item, ok := s.tasks[r.PathValue("id")]
	s.mu.Unlock()
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "task not found"})
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (s *server) updateTask(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title     *string `json:"title"`
		Completed *bool   `json:"completed"`
	}
	if err := readJSON(w, r, &body); err != nil {
		writeError(w, err)
		return
	}
	if body.Title == nil && body.Completed == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "title or completed is required"})
		return
	}
	if body.Title != nil {
		trimmed := strings.TrimSpace(*body.Title)
		if trimmed == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "title must be a non-empty string"})
			return
		}
		body.Title = &trimmed
	}

	s.mu.Lock()
	item, ok := s.tasks[r.PathValue("id")]
	if ok {
		if body.Title != nil {
			item.Title = *body.Title
		}
		if body.Completed != nil {
			item.Completed = *body.Completed
		}
		s.tasks[item.ID] = item
	}
	s.mu.Unlock()
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "task not found"})
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (s *server) deleteTask(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	_, ok := s.tasks[r.PathValue("id")]
	delete(s.tasks, r.PathValue("id"))
	s.mu.Unlock()
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "task not found"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		return errors.New("content-type must be application/json")
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1_000_000)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(dst); err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			return errors.New("request body exceeds 1 MB")
		}
		return fmt.Errorf("invalid JSON: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return errors.New("request body must contain one JSON object")
	}
	return nil
}

func writeError(w http.ResponseWriter, err error) {
	writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
