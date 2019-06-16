package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/partha2312/nlp/data"
	"github.com/partha2312/nlp/nlp"
)

type Server struct {
	router *mux.Router
	ngram  nlp.NGram
}

func New() *Server {
	router := mux.NewRouter()
	ngram := nlp.NewNGram()
	return &Server{
		router: router,
		ngram:  ngram,
	}
}

func (s *Server) Init() {
	s.buildNGram()
	s.router.HandleFunc("/check", s.check)
	s.router.HandleFunc("/search/{word}", s.search)
	fmt.Println("serving on 8080")
	http.ListenAndServe(":8080", s.router)
}

func (s *Server) buildNGram() {
	reader := data.NewReader()
	body, err := reader.Read("data/big.txt")
	if err != nil {
		fmt.Println(err)
	}
	start := time.Now()
	s.ngram.ConstructNGrams(string(body))
	fmt.Println(fmt.Sprintf("all grams done in %v", time.Since(start)))
}

func (s *Server) check(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Server) search(w http.ResponseWriter, r *http.Request) {
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
