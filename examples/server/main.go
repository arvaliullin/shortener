package main

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func mainPage(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "Привет!")
}

func apiPage(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "Это страница /api")
}

func main() {
	router := chi.NewRouter()

	router.Get(`/`, mainPage)
	router.Get(`/api`, apiPage)

	if err := http.ListenAndServe(`:8080`, router); err != nil {
		panic(err)
	}
}
