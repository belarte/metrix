package model

type Entry struct {
	ID       int     `form:"id"`
	MetricID int     `form:"metric"`
	Value    float64 `form:"value"`
	Date     string  `form:"date"`
}

func NewEntry(id, metricID int, value float64, date string) Entry {
	return Entry{
		ID:       id,
		MetricID: metricID,
		Value:    value,
		Date:     date,
	}
}

type Entries []Entry
