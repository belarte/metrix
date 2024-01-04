package server

import (
	"html/template"
	"log"
	"net/http"
)

func handler(templateFileName string) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        t, err := template.ParseFiles("server/templates/main.tmpl", "server/templates/"+templateFileName)
        if err != nil {
            log.Fatal(err)
        }
        t.Execute(w, nil)
    }
}

func Run(addr string) error {
    http.HandleFunc("/", handler("home.tmpl"))
    http.HandleFunc("/manage", handler("manage.tmpl"))
    return http.ListenAndServe(addr, nil)
}
