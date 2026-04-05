package monitoring

import (
	"fmt"
	"os"
	"strings"
)

// NetIface is RX/TX byte totals for one interface (/proc/net/dev).
type NetIface struct {
	Name    string
	RxBytes uint64
	TxBytes uint64
}

// NetworkSnapshot is RX/TX byte totals from /proc/net/dev (since boot).
type NetworkSnapshot struct {
	Ifaces []NetIface
}

// ReadNetwork reads /proc/net/dev (skips lo).
func ReadNetwork() (*NetworkSnapshot, error) {
	b, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		return nil, fmt.Errorf("monitoring/network: %w", err)
	}
	ifaces := parseProcNetDev(string(b))
	if len(ifaces) == 0 {
		return nil, fmt.Errorf("monitoring/network: nenhuma interface encontrada")
	}
	return &NetworkSnapshot{Ifaces: ifaces}, nil
}

func parseProcNetDev(content string) []NetIface {
	lines := strings.Split(content, "\n")
	var out []NetIface
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "|") {
			continue
		}
		colon := strings.IndexByte(line, ':')
		if colon < 0 {
			continue
		}
		name := strings.TrimSpace(line[:colon])
		if name == "lo" {
			continue
		}
		rest := strings.Fields(line[colon+1:])
		if len(rest) < 9 {
			continue
		}
		rx, ok1 := parseUintField(rest[0])
		tx, ok2 := parseUintField(rest[8])
		if !ok1 || !ok2 {
			continue
		}
		out = append(out, NetIface{Name: name, RxBytes: rx, TxBytes: tx})
	}
	return out
}

func parseUintField(s string) (uint64, bool) {
	var n uint64
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, false
		}
		n = n*10 + uint64(c-'0')
	}
	return n, true
}

// NetRates is per-interface throughput (bytes/s).
type NetRates struct {
	RxBps float64
	TxBps float64
}

// ComputeNetRates returns bytes/s per iface between two snapshots.
func ComputeNetRates(prev, cur *NetworkSnapshot, dtSec float64) map[string]NetRates {
	if prev == nil || cur == nil || dtSec <= 0 {
		return nil
	}
	prevMap := make(map[string]NetIface, len(prev.Ifaces))
	for _, iface := range prev.Ifaces {
		prevMap[iface.Name] = iface
	}
	out := make(map[string]NetRates)
	for _, c := range cur.Ifaces {
		p, ok := prevMap[c.Name]
		if !ok {
			continue
		}
		drx := int64(c.RxBytes) - int64(p.RxBytes)
		dtx := int64(c.TxBytes) - int64(p.TxBytes)
		if drx < 0 {
			drx = 0
		}
		if dtx < 0 {
			dtx = 0
		}
		out[c.Name] = NetRates{
			RxBps: float64(drx) / dtSec,
			TxBps: float64(dtx) / dtSec,
		}
	}
	return out
}
