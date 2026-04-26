package services

import "main.go/types"

type StatsRunner struct {
	stats types.RunStats
}

func NewStatsRunner() *StatsRunner {
	return &StatsRunner{}
}

func (r *StatsRunner) Apply(result types.FetchResult) {
	if result.WasDuplicate {
		r.stats.Duplicates++
		return
	}

	r.stats.Queued++
	r.stats.TotalDuration += result.Duration

	if result.FetchErr != nil {
		r.stats.FetchFailed++
		if result.StatusCode != 0 {
			r.stats.Non2xxFailures++
		}
		return
	}

	if result.StoreErr != nil {
		r.stats.StoreFailed++
		return
	}

	if result.Stored {
		r.stats.StoredSuccess++
	}
}

func (r *StatsRunner) Consume(results <-chan types.FetchResult) {
	for result := range results {
		r.Apply(result)
	}
}

func (r *StatsRunner) Stats() types.RunStats {
	return r.stats
}
