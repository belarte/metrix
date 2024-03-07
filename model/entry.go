package model

import "fmt"

type Entry struct {
	MetricID int     `form:"metric" param:"metric_id"`
	Value    float64 `form:"value"`
	Date     string  `form:"date" param:"date"`
}

func NewEntry(metricID int, value float64, date string) Entry {
	return Entry{
		MetricID: metricID,
		Value:    value,
		Date:     date,
	}
}

type Entries []Entry

func (e Entry) String() string {
	return fmt.Sprintf("%2d %s %.2f\n", e.MetricID, e.Date, e.Value)
}

func (e Entries) String() string {
	s := fmt.Sprintln("Id Date       Value")
	for _, entry := range e {
		s += fmt.Sprintf("%2d %s %.2f\n", entry.MetricID, entry.Date, entry.Value)
	}
	return s
}
