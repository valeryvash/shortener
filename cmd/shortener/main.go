package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var runes = []rune(alphabet)
var links = map[string]string{}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	ru := []rune(r.URL.Path)
	// в "бд" хранится только 8 символьная последовательность
	if len(ru) < 9 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// убираем ведущий слеш
	ru = ru[1:]
	if v, ok := links[string(ru)]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		w.Header().Add("Location", v)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hash := sha256.Sum256([]byte(body))

	result := make([]rune, 8)
	for i := range result {
		result[i] = runes[int(hash[i])%len(runes)]
	}
	links[string(result)] = string(body)

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, _ = fmt.Fprintf(w, "http://localhost:8080/%s", string(result))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/" && r.Method == http.MethodPost:
			createHandler(w, r)
		case r.URL.Path != "/" && r.Method == http.MethodGet:
			redirectHandler(w, r)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
