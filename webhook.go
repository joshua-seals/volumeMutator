package main

import (
	"log"
	"net/http"
)

func volumeMutator(w http.ResponseWriter, r *http.Request) {
	log.Printf("Mutator Contacted")
}
