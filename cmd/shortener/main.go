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

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if r.URL.Path == "/" && method == http.MethodPost {
			hash := sha256.Sum256([]byte(body))

			result := make([]rune, 8)
			for i := range result {
				result[i] = runes[int(hash[i])%len(runes)]
			}
			links[string(result)] = string(body)

			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, "http://localhost:8080/%s", string(result))
			return
		}

		if method == http.MethodGet {
			r := []rune(r.URL.Path)
			if len(r) < 9 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			r = r[1:]
			if v, ok := links[string(r)]; !ok {
				w.WriteHeader(http.StatusBadRequest)
				return
			} else {
				w.Header().Add("Location", v)
				w.WriteHeader(http.StatusTemporaryRedirect)
				return
			}
		}

		w.WriteHeader(http.StatusBadRequest)
	})

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
