package diagram

import (
	"bytes"
	"fmt"

	"github.com/belarte/metrix/model"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func graphData(entries []model.Entry) ([]string, []opts.LineData) {
	var x []string
	var y []opts.LineData

	for _, e := range entries {
		x = append(x, e.Date)
		y = append(y, opts.LineData{Value: e.Value})
	}

	return x, y
}

func DataToGraph(metric model.Metric, entries []model.Entry) string {
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
	return s
}
