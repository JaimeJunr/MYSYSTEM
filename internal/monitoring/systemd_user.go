package monitoring

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// FailedUserUnit is one failed user unit from systemctl.
type FailedUserUnit struct {
	Unit        string
	Load        string
	Active      string
	Sub         string
	Description string
}

// SystemdUserFailedSnapshot is failed --user units from systemctl.
type SystemdUserFailedSnapshot struct {
	Units []FailedUserUnit
}

// ReadSystemdUserFailed runs systemctl --user list-units --state=failed.
func ReadSystemdUserFailed() (*SystemdUserFailedSnapshot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "systemctl", "--user", "list-units",
		"--state=failed", "--no-legend", "--no-pager")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg != "" {
			return nil, fmt.Errorf("monitoring/systemd-user: %w: %s", err, msg)
		}
		return nil, fmt.Errorf("monitoring/systemd-user: %w", err)
	}
	lines := strings.Split(string(out), "\n")
	var units []FailedUserUnit
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		u := FailedUserUnit{
			Unit:   fields[0],
			Load:   fields[1],
			Active: fields[2],
			Sub:    fields[3],
		}
		if len(fields) > 4 {
			u.Description = strings.Join(fields[4:], " ")
		}
		units = append(units, u)
	}
	return &SystemdUserFailedSnapshot{Units: units}, nil
}
