package model

type Metric struct {
	ID          int    `form:"id"`
	Title       string `form:"title"`
	Unit        string `form:"unit"`
	Description string `form:"description"`
}

type Metrics []Metric
