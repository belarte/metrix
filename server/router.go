package server

import (
	"html/template"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    t, err := template.ParseFiles("server/main.html")
    if err != nil {
        log.Fatal(err)
    }

    t.Execute(w, nil)
}

func Run(addr string) {
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(addr, nil))
}
