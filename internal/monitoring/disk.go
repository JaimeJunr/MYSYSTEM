package monitoring

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"syscall"
)

// DiskMount is usage for one mount point (statfs).
type DiskMount struct {
	Mountpoint string
	Fstype     string
	TotalBytes uint64
	UsedBytes  uint64
	AvailBytes uint64
	UsePercent float64
}

var skipDiskFstype = map[string]bool{
	"proc":         true,
	"sysfs":        true,
	"devtmpfs":     true,
	"devpts":       true,
	"cgroup":       true,
	"cgroup2":      true,
	"pstore":       true,
	"bpf":          true,
	"tracefs":      true,
	"fusectl":      true,
	"mqueue":       true,
	"securityfs":   true,
	"configfs":     true,
	"ramfs":        true,
	"autofs":       true,
	"binfmt_misc":  true,
	"rpc_pipefs":   true,
	"nfsd":         true,
	"debugfs":      true,
	"hugetlbfs":    true,
}

// ReadDiskMounts reads /proc/mounts and statfs per mount point.
func ReadDiskMounts() ([]DiskMount, error) {
	f, err := os.Open("/proc/mounts")
	if err != nil {
		return nil, fmt.Errorf("monitoring/disk: %w", err)
	}
	defer f.Close()

	seen := make(map[string]string)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) < 3 {
			continue
		}
		mp := fields[1]
		fst := fields[2]
		if skipDiskFstype[fst] {
			continue
		}
		seen[mp] = fst
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("monitoring/disk: %w", err)
	}

	mps := make([]string, 0, len(seen))
	for mp := range seen {
		mps = append(mps, mp)
	}
	sort.Strings(mps)

	out := make([]DiskMount, 0, len(mps))
	for _, mp := range mps {
		fst := seen[mp]
		var st syscall.Statfs_t
		if err := syscall.Statfs(mp, &st); err != nil {
			continue
		}
		bs := uint64(st.Bsize)
		if bs == 0 {
			continue
		}
		total := st.Blocks * bs
		avail := st.Bavail * bs
		var used uint64
		if total > avail {
			used = total - avail
		}
		pct := 0.0
		if total > 0 {
			pct = float64(used) * 100 / float64(total)
		}
		out = append(out, DiskMount{
			Mountpoint: mp,
			Fstype:     fst,
			TotalBytes: total,
			UsedBytes:  used,
			AvailBytes: avail,
			UsePercent: pct,
		})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("monitoring/disk: nenhum mount acessível encontrado")
	}
	return out, nil
}
