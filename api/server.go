package api

import (
	"encoding/json"
	"github.com/divingbeetle/Nonogram-terminal-server/storage"
	"net/http"
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
	http.HandleFunc("/puzzles", s.handlePuzzleRequest)
	return http.ListenAndServe(s.listenAddr, nil)
}

func (s *Server) handlePuzzleRequest(w http.ResponseWriter, r *http.Request) {
	puzzle, err := storage.FetchPuzzle()
	if err != nil {
		http.Error(w, "failed to fetch puzzle", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(puzzle)
}
