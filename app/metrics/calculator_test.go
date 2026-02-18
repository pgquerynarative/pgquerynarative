package metrics

import (
	"testing"
)

func TestCalculateMetrics_TimeSeries_AdvancedFields(t *testing.T) {
	columns := []string{"month", "revenue"}
	profiles := []ColumnProfile{
		{Name: "month", Type: ColumnTypeDate, IsTimeSeries: true},
		{Name: "revenue", Type: ColumnTypeNumeric, IsMeasure: true},
	}
	// 8 periods so we get trend (6), moving avg (3), and enough for potential anomalies
	rows := [][]interface{}{
		{"2025-01", 100.0},
		{"2025-02", 105.0},
		{"2025-03", 98.0},
		{"2025-04", 110.0},
		{"2025-05", 115.0},
		{"2025-06", 120.0},
		{"2025-07", 118.0},
		{"2025-08", 125.0},
	}
	m := CalculateMetrics(columns, rows, profiles, 0.5)
	if m == nil {
		t.Fatal("expected non-nil metrics")
	}
	if len(m.TimeSeries) == 0 {
		t.Fatal("expected time series for revenue")
	}
	ts, ok := m.TimeSeries["revenue"]
	if !ok {
		t.Fatal("expected revenue in time series")
	}
	if ts.CurrentPeriod != 125.0 {
		t.Errorf("current period = %v, want 125", ts.CurrentPeriod)
	}
	if len(ts.Periods) == 0 {
		t.Error("expected periods to be populated")
	}
	if ts.MovingAverage == nil {
		t.Error("expected moving average (8 periods >= 3)")
	}
	if ts.TrendSummary == nil {
		t.Error("expected trend summary (>= 2 periods)")
	}
	if ts.TrendSummary != nil && ts.TrendSummary.Summary == "" {
		t.Error("expected non-empty trend summary")
	}
}

func TestCalculateMetrics_TimeSeries_AnomalyDetection(t *testing.T) {
	columns := []string{"month", "value"}
	profiles := []ColumnProfile{
		{Name: "month", Type: ColumnTypeDate, IsTimeSeries: true},
		{Name: "value", Type: ColumnTypeNumeric, IsMeasure: true},
	}
	// One value far from mean to trigger anomaly (e.g. 1000 when rest ~10)
	rows := [][]interface{}{
		{"2025-01", 10.0},
		{"2025-02", 11.0},
		{"2025-03", 10.5},
		{"2025-04", 1000.0}, // outlier
		{"2025-05", 10.0},
		{"2025-06", 11.0},
	}
	m := CalculateMetrics(columns, rows, profiles, 0.5)
	ts, ok := m.TimeSeries["value"]
	if !ok {
		t.Fatal("expected time series for value")
	}
	if len(ts.Anomalies) == 0 {
		t.Error("expected at least one anomaly (1000 is outlier)")
	}
}

func TestLinearRegression(t *testing.T) {
	// y = 2*x + 1 for x=0,1,2,3,4 -> slope 2, intercept 1
	y := []float64{1, 3, 5, 7, 9}
	slope, intercept := linearRegression(y)
	if slope < 1.99 || slope > 2.01 {
		t.Errorf("slope = %v, want ~2", slope)
	}
	if intercept < 0.99 || intercept > 1.01 {
		t.Errorf("intercept = %v, want ~1", intercept)
	}
}

func TestMeanAndStd(t *testing.T) {
	vals := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	mean, std := meanAndStd(vals)
	if mean < 4.99 || mean > 5.01 {
		t.Errorf("mean = %v, want 5", mean)
	}
	if std < 2.0 || std > 2.1 {
		t.Errorf("std = %v, want ~2", std)
	}
}

func TestCalculateMetrics_NoTimeSeries(t *testing.T) {
	columns := []string{"category", "total"}
	profiles := []ColumnProfile{
		{Name: "category", Type: ColumnTypeText, IsDimension: true},
		{Name: "total", Type: ColumnTypeNumeric, IsMeasure: true},
	}
	rows := [][]interface{}{
		{"A", 10.0},
		{"B", 20.0},
	}
	m := CalculateMetrics(columns, rows, profiles, 0.5)
	if len(m.TimeSeries) != 0 {
		t.Errorf("expected no time series, got %d", len(m.TimeSeries))
	}
}
