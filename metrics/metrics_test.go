package metrics

import "testing"

func TestHasSignificantChange(t *testing.T) {
	oldMetrics := map[string]string{"CPU Usage": "20%", "Memory Usage": "30%"}
	newMetrics := map[string]string{"CPU Usage": "25%", "Memory Usage": "30%"}

	if !HasSignificantChange(oldMetrics, newMetrics, 4.0) {
		t.Errorf("Expected significant change but found none")
	}

	if HasSignificantChange(oldMetrics, newMetrics, 10.0) {
		t.Errorf("Expected no significant change but found one")
	}
}
