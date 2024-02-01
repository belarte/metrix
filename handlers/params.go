package handlers

import "github.com/belarte/metrix/database"

type templateParams struct {
	Metrics           []database.Metric
	Selected          database.Metric
	Content           string
	ButtonLabel       string
	AdditionalMessage string
}
