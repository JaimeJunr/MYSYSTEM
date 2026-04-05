package monitoring

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ThermalReading is one sysfs temperature (°C).
type ThermalReading struct {
	Label string
	TempC float64
}

// ThermalSnapshot groups sysfs thermal readings.
type ThermalSnapshot struct {
	Readings []ThermalReading
}

// ReadThermal reads thermal_zone* and hwmon temp*_input (millidegrees → °C).
func ReadThermal() (*ThermalSnapshot, error) {
	var readings []ThermalReading

	zones, _ := filepath.Glob("/sys/class/thermal/thermal_zone*")
	sort.Strings(zones)
	for _, zdir := range zones {
		tPath := filepath.Join(zdir, "type")
		tempPath := filepath.Join(zdir, "temp")
		if _, err := os.Stat(tempPath); err != nil {
			continue
		}
		label := readTrimFile(tPath)
		if label == "" {
			label = filepath.Base(zdir)
		}
		md, ok := readIntFile(tempPath)
		if !ok {
			continue
		}
		readings = append(readings, ThermalReading{
			Label: label,
			TempC: float64(md) / 1000.0,
		})
	}

	hwmonDirs, _ := filepath.Glob("/sys/class/hwmon/hwmon*")
	sort.Strings(hwmonDirs)
	for _, hdir := range hwmonDirs {
		chip := readTrimFile(filepath.Join(hdir, "name"))
		if chip == "" {
			chip = filepath.Base(hdir)
		}
		ents, err := os.ReadDir(hdir)
		if err != nil {
			continue
		}
		for _, e := range ents {
			name := e.Name()
			if !strings.HasPrefix(name, "temp") || !strings.HasSuffix(name, "_input") {
				continue
			}
			base := strings.TrimSuffix(name, "_input")
			lbl := readTrimFile(filepath.Join(hdir, base+"_label"))
			if lbl == "" {
				lbl = chip + " " + base
			}
			p := filepath.Join(hdir, name)
			v, ok := readIntFile(p)
			if !ok {
				continue
			}
			readings = append(readings, ThermalReading{
				Label: lbl,
				TempC: float64(v) / 1000.0,
			})
		}
	}

	return &ThermalSnapshot{Readings: readings}, nil
}
