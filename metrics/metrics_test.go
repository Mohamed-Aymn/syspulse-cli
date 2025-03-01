package metrics

import (
	"runtime"
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

func TestCpuIntensiveProcess(t *testing.T) {
	initialMetrics := GetMetrics()

	cpuIntensiveTask()

	// Wait for CPU usage to reflect the change
	// time.Sleep(5 * time.Second)

	secondaryMetrics := GetMetrics()

	initialCPU := initialMetrics["CPU Usage"]
	secondaryCPU := secondaryMetrics["CPU Usage"]

	t.Logf("Initial CPU Usage: %s", initialCPU)
	t.Logf("Secondary CPU Usage: %s", secondaryCPU)

	if initialCPU >= secondaryCPU {
		t.Errorf("CPU usage did not increase as expected")
	}
}

var globalData []*[]byte

func memoryIntensiveTask() {
	// Disable GC completely for the test duration
	runtime.GC()
	runtime.MemProfileRate = 0

	// Allocate a large amount of memory (e.g., 1GB) and write to force commitment
	for i := 0; i < 100; i++ { // 100 * 10MB = ~1GB
		data := make([]byte, 10*1024*1024) // 10MB per slice
		for j := range data {
			data[j] = 1 // Write to force physical allocation
		}
		globalData = append(globalData, &data)
	}

	// Ensure memory is actively held
	time.Sleep(5 * time.Second)
}

func TestMemoryIntensiveProcess(t *testing.T) {
	initialMetrics := GetMetrics()

	memoryIntensiveTask()

	secondaryMetrics := GetMetrics()

	initialMem := initialMetrics["Memory Usage"]
	secondaryMem := secondaryMetrics["Memory Usage"]

	t.Logf("Initial Memory Usage: %s", initialMem)
	t.Logf("Secondary Memory Usage: %s", secondaryMem)

	if initialMem >= secondaryMem {
		t.Errorf("Memory usage did not increase as expected")
	}
}

func cpuIntensiveTask() {
	for i := 0; i < 1000; i++ {
		go func() {
			a, b := 0, 1
			for j := 0; j < 1e7; j++ {
				a, b = b, a+b
			}
		}()
	}
}
