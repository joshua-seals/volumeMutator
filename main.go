package main

import (
	"encoding/json"
	"log"
	"net/http"
)

const cert = "./webhook/certs/tls.crt"
const key = "./webhook/certs/tls.key"

func main() {

	// handle our core application
	http.HandleFunc("/volume-mutator", volumeMutator)
	http.HandleFunc("/liveness", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]string{
			"status": "ok",
		}
		w.Header().Set("Content-Type", "application/json")
		js, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "The server encountered an error", http.StatusInternalServerError)
		}
		w.Write(js)
	})
	log.Print("Listening on port 8443...")
	log.Fatal(http.ListenAndServeTLS(":8443", cert, key, nil))
}
