// Package metrics provides typed metrics with counters, gauges, histograms, and Prometheus export.
// Ported from hypercode/go/internal/metrics/metrics.go
//
// WHAT: Thread-safe metrics collection: counters, gauges, histograms, time series
// WHY: All subsystems need observability - tool execution times, memory usage, cache hit rates
// HOW: Labeled metric storage with Prometheus export, system monitoring, and downsampling
package metrics

import (
	"fmt"
	"math"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// MetricEvent is a legacy event-stream record.
type MetricEvent struct {
	Timestamp int64             `json:"timestamp"`
	Type      string            `json:"type"`
	Value     float64           `json:"value"`
	Tags      map[string]string `json:"tags,omitempty"`
}

// HistogramStats holds statistical summaries for a histogram.
type HistogramStats struct {
	Count int     `json:"count"`
	Sum   float64 `json:"sum"`
	Avg   float64 `json:"avg"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	P50   float64 `json:"p50"`
	P95   float64 `json:"p95"`
	P99   float64 `json:"p99"`
}

// AllMetrics holds a snapshot of all metric types.
type AllMetrics struct {
	Counters   map[string]float64        `json:"counters"`
	Gauges     map[string]float64        `json:"gauges"`
	Histograms map[string]HistogramStats `json:"histograms"`
}

// StatsResult holds the result of GetStats().
type StatsResult struct {
	WindowMs    int64                      `json:"windowMs"`
	TotalEvents int                        `json:"totalEvents"`
	Counters    map[string]float64         `json:"counters"`
	Gauges      map[string]float64         `json:"gauges"`
	Histograms  map[string]HistogramStats  `json:"histograms"`
}

// DownsampledBucket is a single bucket in a downsampled time series.
type DownsampledBucket struct {
	Time     int64   `json:"time"`
	Count    int     `json:"count"`
	ValueAvg float64 `json:"value_avg"`
}

type labeledValue struct {
	value  float64
	labels map[string]string
}

// MetricsService provides typed metrics with counters, gauges, histograms.
type MetricsService struct {
	mu sync.RWMutex

	events    []MetricEvent
	maxEvents int

	counters   map[string]*labeledValue
	gauges     map[string]*labeledValue
	histograms map[string][]float64

	monitorStop chan struct{}
	running     int32

	onMetric func(name string, mtype string, value float64, labels map[string]string)
}

var globalInstance atomic.Pointer[MetricsService]

// GetMetricsService returns the global singleton.
func GetMetricsService() *MetricsService {
	if inst := globalInstance.Load(); inst != nil {
		return inst
	}
	inst := NewMetricsService()
	if globalInstance.CompareAndSwap(nil, inst) {
		return inst
	}
	return globalInstance.Load()
}

// NewMetricsService creates a fresh MetricsService.
func NewMetricsService() *MetricsService {
	return &MetricsService{
		events:      make([]MetricEvent, 0, 10000),
		maxEvents:   10000,
		counters:    make(map[string]*labeledValue),
		gauges:      make(map[string]*labeledValue),
		histograms:  make(map[string][]float64),
		monitorStop: make(chan struct{}),
	}
}

// --- Counters ---

func (ms *MetricsService) IncCounter(name string, value float64, labels map[string]string) {
	key := compositeKey(name, labels)
	ms.mu.Lock()
	lv := ms.counters[key]
	if lv == nil {
		lv = &labeledValue{labels: labels}
		ms.counters[key] = lv
	}
	lv.value += value
	ms.mu.Unlock()
	ms.emitCallback(name, "counter", lv.value, labels)
}

func (ms *MetricsService) GetCounter(name string, labels map[string]string) float64 {
	key := compositeKey(name, labels)
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	if lv, ok := ms.counters[key]; ok {
		return lv.value
	}
	return 0
}

// --- Gauges ---

func (ms *MetricsService) SetGauge(name string, value float64, labels map[string]string) {
	key := compositeKey(name, labels)
	ms.mu.Lock()
	lv := ms.gauges[key]
	if lv == nil {
		lv = &labeledValue{labels: labels}
		ms.gauges[key] = lv
	}
	lv.value = value
	ms.mu.Unlock()
	ms.emitCallback(name, "gauge", value, labels)
}

func (ms *MetricsService) GetGauge(name string, labels map[string]string) float64 {
	key := compositeKey(name, labels)
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	if lv, ok := ms.gauges[key]; ok {
		return lv.value
	}
	return 0
}

func (ms *MetricsService) IncGauge(name string, value float64, labels map[string]string) {
	key := compositeKey(name, labels)
	ms.mu.Lock()
	lv := ms.gauges[key]
	if lv == nil {
		lv = &labeledValue{labels: labels}
		ms.gauges[key] = lv
	}
	lv.value += value
	ms.mu.Unlock()
}

func (ms *MetricsService) DecGauge(name string, value float64, labels map[string]string) {
	ms.IncGauge(name, -value, labels)
}

// --- Histograms ---

func (ms *MetricsService) ObserveHistogram(name string, value float64) {
	ms.mu.Lock()
	ms.histograms[name] = append(ms.histograms[name], value)
	ms.mu.Unlock()
}

func (ms *MetricsService) GetHistogramStats(name string) *HistogramStats {
	ms.mu.RLock()
	values, ok := ms.histograms[name]
	ms.mu.RUnlock()
	if !ok || len(values) == 0 {
		return nil
	}
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)

	count := len(sorted)
	sum := 0.0
	for _, v := range sorted {
		sum += v
	}
	return &HistogramStats{
		Count: count, Sum: sum, Avg: sum / float64(count),
		Min: sorted[0], Max: sorted[count-1],
		P50: sorted[int(float64(count)*0.5)],
		P95: sorted[minInt(int(float64(count)*0.95), count-1)],
		P99: sorted[minInt(int(float64(count)*0.99), count-1)],
	}
}

// Timer returns a function that records elapsed duration when called.
func (ms *MetricsService) Timer(name string) func() {
	start := time.Now()
	return func() {
		ms.ObserveHistogram(name, float64(time.Since(start).Milliseconds()))
	}
}

// --- Export ---

// ExportPrometheus returns metrics in Prometheus exposition format.
func (ms *MetricsService) ExportPrometheus() string {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	var lines []string
	for _, name := range baseNames(ms.counters) {
		lines = append(lines, fmt.Sprintf("# TYPE %s counter", name))
		for key, lv := range ms.counters {
			if keyBase(key) != name {
				continue
			}
			if ls := formatLabels(lv.labels); ls != "" {
				lines = append(lines, fmt.Sprintf("%s{%s} %.2f", name, ls, lv.value))
			} else {
				lines = append(lines, fmt.Sprintf("%s %.2f", name, lv.value))
			}
		}
	}
	for _, name := range baseNames(ms.gauges) {
		lines = append(lines, fmt.Sprintf("# TYPE %s gauge", name))
		for key, lv := range ms.gauges {
			if keyBase(key) != name {
				continue
			}
			if ls := formatLabels(lv.labels); ls != "" {
				lines = append(lines, fmt.Sprintf("%s{%s} %.2f", name, ls, lv.value))
			} else {
				lines = append(lines, fmt.Sprintf("%s %.2f", name, lv.value))
			}
		}
	}
	return strings.Join(lines, "\n") + "\n"
}

// GetAll returns a snapshot of all metrics.
func (ms *MetricsService) GetAll() *AllMetrics {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	counters := make(map[string]float64)
	for key, lv := range ms.counters {
		counters[keyBase(key)] += lv.value
	}
	gauges := make(map[string]float64)
	for key, lv := range ms.gauges {
		gauges[keyBase(key)] = lv.value
	}
	histograms := make(map[string]HistogramStats)
	for name := range ms.histograms {
		if stats := ms.GetHistogramStats(name); stats != nil {
			histograms[name] = *stats
		}
	}
	return &AllMetrics{Counters: counters, Gauges: gauges, Histograms: histograms}
}

// Track records a metric event (legacy API).
func (ms *MetricsService) Track(mtype string, value float64, tags map[string]string) {
	ms.mu.Lock()
	ms.events = append(ms.events, MetricEvent{
		Timestamp: time.Now().UnixMilli(), Type: mtype, Value: value, Tags: tags,
	})
	if len(ms.events) > ms.maxEvents {
		ms.events = ms.events[len(ms.events)-ms.maxEvents/2:]
	}
	ms.mu.Unlock()

	switch mtype {
	case "duration", "tool_execution":
		name := mtype
		if tags != nil {
			if t, ok := tags["tool"]; ok { name = t }
		}
		ms.ObserveHistogram("duration_"+name, value)
	case "memory_heap", "memory_rss":
		ms.SetGauge(mtype, value, tags)
	default:
		ms.IncCounter(mtype, value, tags)
	}
}

// StartMonitoring begins periodic system metrics collection.
func (ms *MetricsService) StartMonitoring(intervalMs int64) {
	if intervalMs <= 0 {
		intervalMs = 5000
	}
	if !atomic.CompareAndSwapInt32(&ms.running, 0, 1) {
		return
	}
	go func() {
		ticker := time.NewTicker(time.Duration(intervalMs) * time.Millisecond)
		defer ticker.Stop()
		var memStats runtime.MemStats
		for {
			select {
			case <-ticker.C:
				runtime.ReadMemStats(&memStats)
				ms.Track("memory_heap", float64(memStats.HeapAlloc), nil)
				ms.Track("memory_rss", float64(memStats.Sys), nil)
			case <-ms.monitorStop:
				return
			}
		}
	}()
}

// StopMonitoring stops the monitoring goroutine.
func (ms *MetricsService) StopMonitoring() {
	if atomic.CompareAndSwapInt32(&ms.running, 1, 0) {
		close(ms.monitorStop)
		ms.monitorStop = make(chan struct{})
	}
}

// Reset clears all metrics.
func (ms *MetricsService) Reset() {
	ms.mu.Lock()
	ms.counters = make(map[string]*labeledValue)
	ms.gauges = make(map[string]*labeledValue)
	ms.histograms = make(map[string][]float64)
	ms.events = nil
	ms.mu.Unlock()
}

// OnMetric registers a callback for metric events.
func (ms *MetricsService) OnMetric(fn func(name string, mtype string, value float64, labels map[string]string)) {
	ms.onMetric = fn
}

func (ms *MetricsService) emitCallback(name, mtype string, value float64, labels map[string]string) {
	if ms.onMetric != nil {
		ms.onMetric(name, mtype, value, labels)
	}
}

// --- internal helpers ---

func compositeKey(name string, labels map[string]string) string {
	if len(labels) == 0 {
		return name
	}
	keys := make([]string, 0, len(labels))
	for k := range labels { keys = append(keys, k) }
	sort.Strings(keys)
	parts := make([]string, len(keys))
	for i, k := range keys { parts[i] = k + "=" + labels[k] }
	return name + "{" + strings.Join(parts, ",") + "}"
}

func keyBase(key string) string {
	if idx := strings.IndexByte(key, '{'); idx >= 0 {
		return key[:idx]
	}
	return key
}

func baseNames(m map[string]*labeledValue) []string {
	seen := make(map[string]bool)
	var names []string
	for key := range m {
		base := keyBase(key)
		if !seen[base] {
			seen[base] = true
			names = append(names, base)
		}
	}
	sort.Strings(names)
	return names
}

func formatLabels(labels map[string]string) string {
	if len(labels) == 0 {
		return ""
	}
	keys := make([]string, 0, len(labels))
	for k := range labels { keys = append(keys, k) }
	sort.Strings(keys)
	parts := make([]string, len(keys))
	for i, k := range keys { parts[i] = fmt.Sprintf(`%s="%s"`, k, labels[k]) }
	return strings.Join(parts, ",")
}

func minInt(a, b int) int {
	if a < b { return a }
	return b
}

// Ensure math and time are used
var _ = math.Pi
var _ = time.Second
