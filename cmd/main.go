package main

import (
	"fmt"
	"time"

	"github.com/partha2312/nlp/data"
	"github.com/partha2312/nlp/nlp"
	"github.com/partha2312/nlp/service"
)

func main() {
	fmt.Println("hello world")
	server := service.New(buildNGram())
	server.InitRoutes()
}

func buildNGram() nlp.NGram {
	nGram := nlp.NewNGram()
	reader := data.NewReader()
	body, err := reader.Read("data/big.txt")
	if err != nil {
		fmt.Println(err)
	}
	start := time.Now()
	nGram.ConstructNGrams(string(body))
	fmt.Println(fmt.Sprintf("all grams done in %v", time.Since(start)))
	return nGram
}
