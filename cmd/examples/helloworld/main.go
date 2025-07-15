package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/starfederation/datastar-go/datastar"
)

//go:embed hello-world.html
var helloWorldHTML []byte

const port = 1337

func main() {
	r := chi.NewRouter()

	const message = "Hello, world!"
	type Store struct {
		Delay time.Duration `json:"delay"`
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(helloWorldHTML)
	})

	r.Get("/hello-world", func(w http.ResponseWriter, r *http.Request) {
		store := &Store{}
		if err := datastar.ReadSignals(r, store); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sse := datastar.NewSSE(w, r)

		for i := 0; i < len(message); i++ {
			if err := sse.PatchElements(`<div id="message">` + message[:i+1] + `</div>`); err != nil {
				return
			}
			time.Sleep(store.Delay * time.Millisecond)
		}
	})

	log.Printf("Starting server on http://localhost:%d", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), r); err != nil {
		panic(err)
	}
}
