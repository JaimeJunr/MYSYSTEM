package monitoring

import "testing"

func TestParseProcNetDev(t *testing.T) {
	sample := `Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
    lo: 123456      10    0    0    0     0          0         0    123456      10    0    0    0     0       0          0
  eth0: 1000 20 0 0 0 0 0 0 2000 15 0 0 0 0 0 0
`
	ifaces := parseProcNetDev(sample)
	if len(ifaces) != 1 {
		t.Fatalf("want 1 iface (no lo), got %d", len(ifaces))
	}
	if ifaces[0].Name != "eth0" || ifaces[0].RxBytes != 1000 || ifaces[0].TxBytes != 2000 {
		t.Fatalf("unexpected iface: %+v", ifaces[0])
	}
}

func TestComputeNetRates(t *testing.T) {
	prev := &NetworkSnapshot{Ifaces: []NetIface{{Name: "e", RxBytes: 100, TxBytes: 50}}}
	cur := &NetworkSnapshot{Ifaces: []NetIface{{Name: "e", RxBytes: 200, TxBytes: 80}}}
	r := ComputeNetRates(prev, cur, 2.0)
	if r["e"].RxBps != 50 || r["e"].TxBps != 15 {
		t.Fatalf("rates: %+v", r["e"])
	}
}
