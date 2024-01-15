package server

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/belarte/metrix/database"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/labstack/echo/v4"
)

type ReportsHandler struct {
	db *database.InMemory
}

func NewReportsHandler(db *database.InMemory) *ReportsHandler {
	return &ReportsHandler{
		db: db,
	}
}

func (h *ReportsHandler) Reports(c echo.Context) error {
	metrics, err := h.db.GetMetrics()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Render(http.StatusOK, "page", templateParams{
		Content: "reports",
		Metrics: metrics,
	})
}

func graphData(entries []database.Entry) ([]string, []opts.LineData) {
	var x []string
	var y []opts.LineData

	for _, e := range entries {
		x = append(x, e.Date)
		y = append(y, opts.LineData{Value: e.Value})
	}

	return x, y
}

func (h *ReportsHandler) Select(c echo.Context) error {
	val := c.FormValue("reports-select")
	if val == "select-reports" {
		return c.String(http.StatusOK, "Please select a metric.")
	}

	id, err := strconv.Atoi(val)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	metric, err := h.db.GetMetric(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	entries, err := h.db.GetSortedEntriesForMetric(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	line := charts.NewLine()
	line.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    metric.Title,
		Subtitle: metric.Description,
	}))

	x, y := graphData(entries)
	label := fmt.Sprintf("%s (%s)", metric.Title, metric.Unit)
	line.SetXAxis(x).AddSeries(label, y)

	buff := new(bytes.Buffer)
	line.Render(buff)
	lines := bytes.Split(buff.Bytes(), []byte("\n"))
	lines = lines[10 : len(lines)-8]
	s := string(bytes.Join(lines, []byte("\n")))

	return c.String(http.StatusOK, s)
}
