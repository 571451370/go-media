package pdsample

type RangeEntry struct {
	Min, Max float64
}

type RangeList struct {
	Ranges     []RangeEntry
	RangesSize int
}