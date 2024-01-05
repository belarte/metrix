package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/belarte/metrix/database"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("server/templates/main.tmpl", "server/templates/home.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, nil)
}

type manageParams struct {
    Metrics []database.Metric
}

func manageHandler(db database.InMemory) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("server/templates/main.tmpl", "server/templates/manage.tmpl")
		if err != nil {
			log.Fatal(err)
		}

        metrics, _ := db.GetMetrics()

		t.Execute(w, manageParams{metrics})
	}
}

func Run(addr string, db database.InMemory) error {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/manage", manageHandler(db))
	return http.ListenAndServe(addr, nil)
}
