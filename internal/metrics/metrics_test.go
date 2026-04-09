package metrics

import (
	"testing"
	"time"
)

func TestNewMetricsService(t *testing.T) {
	ms := NewMetricsService()
	if ms == nil {
		t.Fatal("should create metrics service")
	}
}

func TestCounter(t *testing.T) {
	ms := NewMetricsService()
	ms.IncCounter("requests", 1, nil)
	ms.IncCounter("requests", 2, nil)
	if v := ms.GetCounter("requests", nil); v != 3 {
		t.Errorf("expected 3, got %f", v)
	}
}

func TestCounterWithLabels(t *testing.T) {
	ms := NewMetricsService()
	ms.IncCounter("requests", 1, map[string]string{"method": "GET"})
	ms.IncCounter("requests", 1, map[string]string{"method": "POST"})
	ms.IncCounter("requests", 2, map[string]string{"method": "GET"})

	getCount := ms.GetCounter("requests", map[string]string{"method": "GET"})
	postCount := ms.GetCounter("requests", map[string]string{"method": "POST"})

	if getCount != 3 {
		t.Errorf("GET: expected 3, got %f", getCount)
	}
	if postCount != 1 {
		t.Errorf("POST: expected 1, got %f", postCount)
	}
}

func TestGauge(t *testing.T) {
	ms := NewMetricsService()
	ms.SetGauge("temperature", 72.5, nil)
	if v := ms.GetGauge("temperature", nil); v != 72.5 {
		t.Errorf("expected 72.5, got %f", v)
	}
	ms.SetGauge("temperature", 73.0, nil)
	if v := ms.GetGauge("temperature", nil); v != 73.0 {
		t.Errorf("expected 73.0, got %f", v)
	}
}

func TestGaugeIncDec(t *testing.T) {
	ms := NewMetricsService()
	ms.SetGauge("connections", 10, nil)
	ms.IncGauge("connections", 5, nil)
	if v := ms.GetGauge("connections", nil); v != 15 {
		t.Errorf("expected 15, got %f", v)
	}
	ms.DecGauge("connections", 3, nil)
	if v := ms.GetGauge("connections", nil); v != 12 {
		t.Errorf("expected 12, got %f", v)
	}
}

func TestHistogram(t *testing.T) {
	ms := NewMetricsService()
	for i := 0; i < 100; i++ {
		ms.ObserveHistogram("latency", float64(i))
	}
	stats := ms.GetHistogramStats("latency")
	if stats == nil {
		t.Fatal("should have stats")
	}
	if stats.Count != 100 {
		t.Errorf("count: %d", stats.Count)
	}
	if stats.Min != 0 {
		t.Errorf("min: %f", stats.Min)
	}
	if stats.Max != 99 {
		t.Errorf("max: %f", stats.Max)
	}
	if stats.Avg != 49.5 {
		t.Errorf("avg: %f", stats.Avg)
	}
}

func TestTimer(t *testing.T) {
	ms := NewMetricsService()
	done := ms.Timer("operation")
	time.Sleep(10 * time.Millisecond)
	done()

	stats := ms.GetHistogramStats("operation")
	if stats == nil {
		t.Fatal("should have timer stats")
	}
	if stats.Count != 1 {
		t.Errorf("count: %d", stats.Count)
	}
	if stats.Min < 8 {
		t.Errorf("too fast: %fms", stats.Min)
	}
}

func TestExportPrometheus(t *testing.T) {
	ms := NewMetricsService()
	ms.IncCounter("http_requests", 42, map[string]string{"method": "GET"})
	ms.SetGauge("temperature", 72.5, nil)

	output := ms.ExportPrometheus()
	if output == "" {
		t.Error("should produce output")
	}
	if !contains(output, "http_requests") {
		t.Error("should contain counter")
	}
	if !contains(output, "temperature") {
		t.Error("should contain gauge")
	}
}

func TestGetAll(t *testing.T) {
	ms := NewMetricsService()
	ms.IncCounter("test_counter", 5, nil)
	ms.SetGauge("test_gauge", 42.0, nil)
	ms.ObserveHistogram("test_hist", 100.0)

	all := ms.GetAll()
	if all.Counters["test_counter"] != 5 {
		t.Errorf("counter: %f", all.Counters["test_counter"])
	}
	if all.Gauges["test_gauge"] != 42.0 {
		t.Errorf("gauge: %f", all.Gauges["test_gauge"])
	}
	if all.Histograms["test_hist"].Count != 1 {
		t.Errorf("hist count: %d", all.Histograms["test_hist"].Count)
	}
}

func TestTrack(t *testing.T) {
	ms := NewMetricsService()
	ms.Track("custom_event", 1.0, map[string]string{"source": "test"})
	ms.Track("duration", 100.0, map[string]string{"tool": "bash"})
	ms.Track("memory_heap", 1024.0, nil)

	// Duration should go to histogram
	stats := ms.GetHistogramStats("duration_bash")
	if stats == nil {
		t.Error("should have duration histogram")
	}

	// Memory should be gauge
	if v := ms.GetGauge("memory_heap", nil); v != 1024.0 {
		t.Errorf("gauge: %f", v)
	}
}

func TestReset(t *testing.T) {
	ms := NewMetricsService()
	ms.IncCounter("test", 5, nil)
	ms.Reset()
	if v := ms.GetCounter("test", nil); v != 0 {
		t.Errorf("should be 0 after reset, got %f", v)
	}
}

func TestCallback(t *testing.T) {
	ms := NewMetricsService()
	var received string
	ms.OnMetric(func(name string, mtype string, value float64, labels map[string]string) {
		received = name
	})
	ms.IncCounter("test_cb", 1, nil)
	if received != "test_cb" {
		t.Errorf("callback received: %s", received)
	}
}

func TestGetHistogramStatsEmpty(t *testing.T) {
	ms := NewMetricsService()
	stats := ms.GetHistogramStats("nonexistent")
	if stats != nil {
		t.Error("should be nil for nonexistent histogram")
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(len(s) > 0 && len(sub) > 0 && findSubstring(s, sub)))
}

func findSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
