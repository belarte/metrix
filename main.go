package main

import (
	"html/template"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    t, err := template.ParseFiles("main.html")
    if err != nil {
        log.Fatal(err)
    }

    t.Execute(w, nil)
}

func main() {
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
