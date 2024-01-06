package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/belarte/metrix/database"
	"github.com/gorilla/schema"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("server/templates/main.tmpl", "server/templates/home.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, nil)
}

type manageParams struct {
	Metrics  []database.Metric
	Selected database.Metric
}

func manageHandler(db *database.InMemory) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("server/templates/main.tmpl", "server/templates/manage.tmpl")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		metrics, err := db.GetMetrics()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		t.Execute(w, manageParams{metrics, database.Metric{ID: -1}})
	}
}

func clickHandler(db *database.InMemory) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var metric database.Metric
		decoder := schema.NewDecoder()
		if err := decoder.Decode(&metric, r.PostForm); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		metric, err := db.AddMetric(metric)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		t, err := template.ParseFiles("server/templates/manage.tmpl")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		metrics, err := db.GetMetrics()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		t.ExecuteTemplate(w, "content", manageParams{metrics, metric})
	}
}

func Run(addr string, db *database.InMemory) error {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/manage", manageHandler(db))
	http.HandleFunc("/manage/click", clickHandler(db))
	return http.ListenAndServe(addr, nil)
}
