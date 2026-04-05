package monitoring

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// LoadSnapshot is parsed /proc/loadavg plus NumCPU.
type LoadSnapshot struct {
	Load1  float64
	Load5  float64
	Load15 float64
	Procs  string
	CPUs   int
}

// ReadLoadAvg parses /proc/loadavg.
func ReadLoadAvg() (*LoadSnapshot, error) {
	b, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return nil, fmt.Errorf("monitoring/loadavg: %w", err)
	}
	line := strings.TrimSpace(string(b))
	fields := strings.Fields(line)
	if len(fields) < 5 {
		return nil, fmt.Errorf("monitoring/loadavg: formato inesperado")
	}
	l1, err1 := strconv.ParseFloat(fields[0], 64)
	l5, err5 := strconv.ParseFloat(fields[1], 64)
	l15, err15 := strconv.ParseFloat(fields[2], 64)
	if err1 != nil || err5 != nil || err15 != nil {
		return nil, fmt.Errorf("monitoring/loadavg: valores inválidos")
	}
	return &LoadSnapshot{
		Load1:  l1,
		Load5:  l5,
		Load15: l15,
		Procs:  fields[3],
		CPUs:   runtime.NumCPU(),
	}, nil
}
