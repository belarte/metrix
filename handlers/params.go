package handlers

import "github.com/belarte/metrix/model"

type templateParams struct {
	Metrics           []model.Metric
	Selected          model.Metric
	Content           string
	ButtonLabel       string
	AdditionalMessage string
}
