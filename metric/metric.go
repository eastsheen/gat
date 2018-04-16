package metric

import (
	"fmt"
	"io"
	"time"
)

type (
	Result struct {
		Name     string
		Err      error
		Request  string
		Response string
		Latency  time.Duration
	}

	BenchResult struct{}

	// Metric hold metrics
	Metric struct {
		name    string
		rate    *AvgRateCounter
		request Counter
		fail    Counter
		success Counter
	}
)

// NewMetric create a metric with interval
func NewMetric(name string, interval time.Duration) *Metric {
	return &Metric{
		name: name,
		rate: NewAvgRateCounter(interval),
	}
}

// OnSuccess ...
func (m *Metric) OnSuccess(success bool, latency time.Duration) {
	m.request.Incr(1)
	if success {
		m.success.Incr(1)
	} else {
		m.fail.Incr(1)
	}
	m.rate.Incr(latency.Nanoseconds())
}

// Report metric to writer
func (m *Metric) Report(w io.Writer) {
	fmt.Fprintf(w, "BENCHMARK:[%s]\tTOTAL:[%d]\tSUCCESS:[%d]\tFAIL:[%d]\tQPS:[%d]\tLATENCY:[%d ms]\n", m.name, m.request.Value(), m.success.Value(), m.fail.Value(), m.rate.Hits(), time.Duration(m.rate.Rate())/time.Millisecond)
}

func (r *Result) Report(w io.Writer) {
	report := fmt.Sprintf("REQUEST:\t%s\nRESPONSE:\t%s\nLATENCY:\t[%d\tms]\nERROR:\t\t%#v\n", r.Request, r.Response, r.Latency/time.Millisecond, r.Err)
	fmt.Fprintln(w, report)
}
