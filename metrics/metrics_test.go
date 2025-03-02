package metrics

import (
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

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

func parseMetric(metric string) float64 {
	value, err := strconv.ParseFloat(strings.TrimSuffix(metric, "%"), 64)
	if err != nil {
		return 0.0
	}
	return value
}

func TestCpuIntensiveProcess(t *testing.T) {
	initialMetrics := GetMetrics("test-device")

	// Run CPU-intensive task
	cpuIntensiveTask()

	// Wait for CPU usage to reflect the change
	time.Sleep(2 * time.Second)

	secondaryMetrics := GetMetrics("test-device")

	initialCPU := parseMetric(initialMetrics["CPU Usage"])
	secondaryCPU := parseMetric(secondaryMetrics["CPU Usage"])

	t.Logf("Initial CPU Usage: %.2f%%", initialCPU)
	t.Logf("Secondary CPU Usage: %.2f%%", secondaryCPU)

	if secondaryCPU <= initialCPU {
		t.Errorf("Expected CPU usage to increase, but it did not")
	}
}

var globalData [][]byte
var mu sync.Mutex

func memoryIntensiveTask() {
	// Disable GC for the test duration
	runtime.GC()
	runtime.MemProfileRate = 0

	mu.Lock()
	defer mu.Unlock()

	// Allocate a large amount of memory (~1GB)
	for i := 0; i < 100; i++ { // 100 * 10MB = ~1GB
		data := make([]byte, 10*1024*1024) // 10MB per slice
		for j := range data {
			data[j] = 1 // Write to force physical allocation
		}
		globalData = append(globalData, data)
	}

	// Hold memory for some time
	time.Sleep(3 * time.Second)
}

func TestMemoryIntensiveProcess(t *testing.T) {
	initialMetrics := GetMetrics("test-device")

	// Run memory-intensive task
	memoryIntensiveTask()

	secondaryMetrics := GetMetrics("test-device")

	initialMem := parseMetric(initialMetrics["Memory Usage"])
	secondaryMem := parseMetric(secondaryMetrics["Memory Usage"])

	t.Logf("Initial Memory Usage: %.2f%%", initialMem)
	t.Logf("Secondary Memory Usage: %.2f%%", secondaryMem)

	if secondaryMem <= initialMem {
		t.Errorf("Expected Memory usage to increase, but it did not")
	}

	// Free allocated memory
	mu.Lock()
	globalData = nil
	mu.Unlock()
	runtime.GC() // Force garbage collection
}

func cpuIntensiveTask() {
	var wg sync.WaitGroup
	numGoroutines := 100

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			a, b := 0, 1
			for j := 0; j < 1e6; j++ { // Reduced to 1e6 for better test stability
				a, b = b, a+b
			}
		}()
	}

	wg.Wait()
}
