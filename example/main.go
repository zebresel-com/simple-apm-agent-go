package main

import (
	"log"
	"net/http"
	"simple-apm-go-agent"
)

func middlewareOne(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middlewareOne")
		next.ServeHTTP(w, r)
		log.Println("Executing middlewareOne again")
	})
}

func middlewareTwo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middlewareTwo")
		if r.URL.Path != "/" {
			return
		}
		next.ServeHTTP(w, r)
		log.Println("Executing middlewareTwo again")
	})
}

func final(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing finalHandler")
	w.Write([]byte("OK"))
}

func main() {
	finalHandler := http.HandlerFunc(final)

	agent := simpleapm.Agent{
		APMUrl:    "http://localhost:3000",
		AppId:     "86b802d8-9669-4fd7-a3fc-9be3c72c4cbc",
		AppSecret: "86b802d8-9669-4fd7-a3fc-9be3c72c4cbc",
	}

	http.Handle("/", agent.Middleware(middlewareTwo(finalHandler)))
	http.ListenAndServe(":3001", nil)
}
