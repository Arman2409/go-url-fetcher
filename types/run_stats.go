package types

import "time"

type FetchResult struct {
	URL          string
	StatusCode   int
	Duration     time.Duration
	Stored       bool
	FetchErr     error
	StoreErr     error
	WasDuplicate bool
}

type RunStats struct {
	Queued         int
	Duplicates     int
	FetchFailed    int
	StoreFailed    int
	StoredSuccess  int
	Non2xxFailures int
	TotalDuration  time.Duration
}

func (s RunStats) AvgDuration() time.Duration {
	if s.Queued == 0 {
		return 0
	}
	return s.TotalDuration / time.Duration(s.Queued)
}
