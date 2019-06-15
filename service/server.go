package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/partha2312/nlp/nlp"
)

type Server struct {
	router *mux.Router
	ngram  nlp.NGram
}

func New(ngram nlp.NGram) *Server {
	router := mux.NewRouter()
	return &Server{
		router: router,
		ngram:  ngram,
	}
}

func (s *Server) InitRoutes() {
	s.router.HandleFunc("/check", s.Check)
	s.router.HandleFunc("/search/{word}", s.Search)
	fmt.Println("serving on 8080")
	http.ListenAndServe(":8080", s.router)
}

func (s *Server) Check(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Server) Search(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	word := vars["word"]

	if len(word) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(fmt.Sprintf("serving for %s", word))

	matrix := s.ngram.Fetch(word)
	j, err := json.Marshal(matrix)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}