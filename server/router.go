package server

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/belarte/metrix/database"
	"github.com/gorilla/schema"
)

type templateParams struct {
	Metrics  []database.Metric
	Selected database.Metric
	Content  string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("server/templates/main.tmpl", "server/templates/home.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, templateParams{Content: "content"})
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

		t.Execute(w, templateParams{
			Metrics:  metrics,
			Selected: database.Metric{ID: -1},
			Content:  "content",
		})
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

		t.ExecuteTemplate(w, "content", templateParams{
			metrics,
			metric,
			"content",
		})
	}
}

func selectHandler(db *database.InMemory) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(r.Form.Get("manage-select"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("id: %d", id)
		log.Println(db.GetMetric(id))

		metric, err := db.GetMetric(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		metrics, err := db.GetMetrics()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		t, err := template.ParseFiles("server/templates/manage.tmpl")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		t.ExecuteTemplate(w, "content", templateParams{
			metrics,
			metric,
			"content",
		})
	}
}

func Run(addr string, db *database.InMemory) error {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/manage", manageHandler(db))
	http.HandleFunc("/manage/click", clickHandler(db))
	http.HandleFunc("/manage/select", selectHandler(db))
	return http.ListenAndServe(addr, nil)
}
