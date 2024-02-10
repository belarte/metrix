package model

type Entry struct {
	MetricID int     `form:"metric"`
	Value    float64 `form:"value"`
	Date     string  `form:"date"`
}

func NewEntry(metricID int, value float64, date string) Entry {
	return Entry{
		MetricID: metricID,
		Value:    value,
		Date:     date,
	}
}

type Entries []Entry
