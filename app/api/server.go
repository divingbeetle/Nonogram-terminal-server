package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/divingbeetle/Nonogram-terminal-server/storage"
)

type Server struct {
	listenAddr string
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
	}
}

func (s *Server) Start() error {
	http.HandleFunc("GET /puzzles", s.getPuzzles)
	http.HandleFunc("GET /puzzles/{id}", s.getPuzzleById)
	http.HandleFunc("GET /puzzles/random", s.getRandomPuzzle)
	return http.ListenAndServe(s.listenAddr, nil)
}

func (s *Server) getPuzzles(w http.ResponseWriter, r *http.Request) {
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		http.Error(w, "invalid offset", http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		http.Error(w, "invalid limit", http.StatusBadRequest)
		return
	}

	const maxLimit = 100
	if limit > maxLimit {
		http.Error(w, "limit is too large", http.StatusBadRequest)
		return
	}

	puzzles, err := storage.FetchPuzzles(offset, limit)
	if err != nil {
		http.Error(w, "failed to fetch puzzles", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
    w.Header().Set("X-Total-Count", strconv.Itoa(len(puzzles)))
	json.NewEncoder(w).Encode(puzzles)
}

func (s *Server) getPuzzleById(w http.ResponseWriter, r *http.Request) {
	idstr := r.PathValue("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	puzzle, err := storage.FetchPuzzle(id)
	if err != nil {
		http.Error(w, "failed to fetch puzzle", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(puzzle)
}

func (s *Server) getRandomPuzzle(w http.ResponseWriter, r *http.Request) {
	puzzle, err := storage.FetchRandomPuzzle()
	if err != nil {
		http.Error(w, "failed to fetch puzzle", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(puzzle)
}
