package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/monitoring"
	btmsg "github.com/JaimeJunr/Homestead/internal/tui/msg"
	"github.com/JaimeJunr/Homestead/internal/tui/theme"
)

func (m Model) withClearedNativeMonitors() Model {
	m.nativeBattery, m.nativeMemory = nil, nil
	m.nativeBatteryErr, m.nativeMemoryErr = nil, nil
	m.nativeDisk, m.nativeDiskErr = nil, nil
	m.nativeLoad, m.nativeLoadErr = nil, nil
	m.nativeNetwork, m.nativeNetworkErr = nil, nil
	m.nativeNetworkAt = time.Time{}
	m.nativeNetRates = nil
	m.nativeThermal, m.nativeThermalErr = nil, nil
	m.nativeSystemdUser, m.nativeSystemdErr = nil, nil
	return m
}

const nativeMonitorRefreshInterval = 3 * time.Second

func (m Model) nativeMonitorScheduleTick() tea.Cmd {
	d := nativeMonitorRefreshInterval
	if m.prefs.ReduceMotion {
		d = 10 * time.Second
	}
	return tea.Tick(d, func(time.Time) tea.Msg {
		return btmsg.NativeMonitorTick{}
	})
}

func (m Model) nativeMonitorLoadCmd() tea.Cmd {
	kind := m.nativeMonitorKind
	return func() tea.Msg {
		out := btmsg.NativeMonitorReload{Kind: kind}
		switch kind {
		case entities.NativeMonitorBattery:
			out.Battery, out.Err = monitoring.ReadBattery()
		case entities.NativeMonitorMemory:
			out.Memory, out.Err = monitoring.ReadMemory()
		case entities.NativeMonitorDisk:
			out.Disk, out.Err = monitoring.ReadDiskMounts()
		case entities.NativeMonitorLoad:
			out.Load, out.Err = monitoring.ReadLoadAvg()
		case entities.NativeMonitorNetwork:
			out.Network, out.Err = monitoring.ReadNetwork()
		case entities.NativeMonitorThermal:
			out.Thermal, out.Err = monitoring.ReadThermal()
		case entities.NativeMonitorSystemdUser:
			out.SystemdUser, out.Err = monitoring.ReadSystemdUserFailed()
		default:
			out.Err = fmt.Errorf("monitor desconhecido: %q", kind)
		}
		return out
	}
}

func (m Model) renderNativeMonitorView() string {
	boxW := scriptOutputCardWidth(m.width)
	head := theme.Title.Render("Homestead") + "\n" +
		theme.Help.Render("Gerenciador de Sistema") + "\n" +
		scriptOutputDivider(boxW) + "\n"

	var body string
	switch m.nativeMonitorKind {
	case entities.NativeMonitorBattery:
		body = renderNativeBatteryPanel(m)
	case entities.NativeMonitorMemory:
		body = renderNativeMemoryPanel(m)
	case entities.NativeMonitorDisk:
		body = renderNativeDiskPanel(m)
	case entities.NativeMonitorLoad:
		body = renderNativeLoadPanel(m)
	case entities.NativeMonitorNetwork:
		body = renderNativeNetworkPanel(m)
	case entities.NativeMonitorThermal:
		body = renderNativeThermalPanel(m)
	case entities.NativeMonitorSystemdUser:
		body = renderNativeSystemdUserPanel(m)
	default:
		body = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("Monitor inválido.")
	}

	sec := 3
	if m.prefs.ReduceMotion {
		sec = 10
	}
	footer := theme.ScriptScreenFooterBar.Width(max(12, boxW-8)).Render(
		fmt.Sprintf("r: atualizar agora · Enter / Esc / q: voltar · ? ajuda · atualiza a cada %ds", sec),
	)
	content := head + body + "\n" + footer
	box := theme.ScriptScreenOuter.Width(boxW)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box.Render(content))
}

func renderNativeBatteryPanel(m Model) string {
	title := theme.ScriptScreenAccent.Render("🔋 Monitor de bateria")

	if m.nativeBatteryErr != nil {
		return title + "\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(m.nativeBatteryErr.Error())
	}
	b := m.nativeBattery
	if b == nil {
		return title + "\n\n" + theme.Help.Render("Carregando…")
	}

	var sb strings.Builder
	sb.WriteString(title)
	sb.WriteString("\n\n")

	kv := func(k, v string) {
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Width(22).Render(k))
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render(v))
		sb.WriteString("\n")
	}

	status := b.Status
	if status == "Charging" {
		status = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render(status)
	} else if status == "Discharging" {
		status = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render(status)
	} else {
		status = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render(status)
	}
	kv("Status", status)

	if b.Capacity >= 0 {
		kv("Capacidade", fmt.Sprintf("%d %%", b.Capacity))
	}
	if b.CapacityLevel != "" && b.CapacityLevel != "Unknown" {
		kv("Nível", b.CapacityLevel)
	}

	if b.ACOnline != nil {
		var acLine string
		if *b.ACOnline {
			acLine = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("conectado")
		} else {
			acLine = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("desconectado")
		}
		if b.ACName != "" {
			acLine += theme.Help.Render("  (" + b.ACName + ")")
		}
		kv("Carregador", acLine)
	}

	sb.WriteString("\n")
	sb.WriteString(theme.ScriptScreenAccent.Render("Detalhes") + "\n")

	if b.EnergyNowUWh > 0 {
		kv("Energia agora", fmt.Sprintf("%.2f Wh", float64(b.EnergyNowUWh)/1e6))
	}
	if b.EnergyFullUWh > 0 {
		kv("Energia cheia", fmt.Sprintf("%.2f Wh", float64(b.EnergyFullUWh)/1e6))
	}
	if b.EnergyDesignUWh > 0 {
		kv("Design (cheia)", fmt.Sprintf("%.2f Wh", float64(b.EnergyDesignUWh)/1e6))
	}
	if h, ok := b.HealthPercent(); ok {
		kv("Saúde (est.)", fmt.Sprintf("%.1f %%", h))
	}
	if w, ok := b.PowerWatts(); ok {
		prefix := ""
		if b.Status == "Charging" {
			prefix = "+"
		} else if b.Status == "Discharging" {
			prefix = "−"
		}
		kv("Potência", fmt.Sprintf("%s%.2f W", prefix, w))
	}
	if b.VoltageNowUV > 0 {
		kv("Tensão", fmt.Sprintf("%.2f V", float64(b.VoltageNowUV)/1e6))
	}
	if b.CycleCount > 0 {
		kv("Ciclos", fmt.Sprintf("%d", b.CycleCount))
	}
	if b.Technology != "" {
		kv("Tecnologia", b.Technology)
	}
	if b.Manufacturer != "" {
		kv("Fabricante", b.Manufacturer)
	}
	if b.ModelName != "" {
		kv("Modelo", b.ModelName)
	}

	sb.WriteString("\n")
	sb.WriteString(theme.ScriptScreenAccent.Render("Resumo") + "\n")
	sb.WriteString(batteryStatusLine(b))
	if est := batteryChargeETA(b); est != "" {
		sb.WriteString(theme.Help.Render(est) + "\n")
	}

	return sb.String()
}

func batteryStatusLine(b *monitoring.BatterySnapshot) string {
	if b == nil {
		return ""
	}
	acOn := b.ACOnline != nil && *b.ACOnline
	switch {
	case b.Status == "Charging" && acOn:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("✓ Carregando com o adaptador conectado.")
	case b.Status == "Full" && acOn:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("✓ Bateria cheia.")
	case b.Status == "Discharging" && !acOn:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render("⚠ Descarregando (sem AC).")
	case b.Status == "Discharging" && acOn:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("⚠ AC conectado, mas status é descarga — verifique cabo/adaptador.")
	default:
		return theme.Help.Render("Status: " + b.Status)
	}
}

func batteryChargeETA(b *monitoring.BatterySnapshot) string {
	if b == nil || b.Status != "Charging" || b.PowerNowUW <= 0 {
		return ""
	}
	if b.EnergyFullUWh <= 0 || b.EnergyNowUWh < 0 {
		return ""
	}
	rem := b.EnergyFullUWh - b.EnergyNowUWh
	if rem <= 0 {
		return ""
	}
	h := float64(rem) / float64(b.PowerNowUW)
	if h <= 0 || h > 48 {
		return ""
	}
	hh := int(h)
	mm := int((h - float64(hh)) * 60)
	return fmt.Sprintf("Tempo estimado até cheio: ~%dh %dm", hh, mm)
}

func renderNativeMemoryPanel(m Model) string {
	title := theme.ScriptScreenAccent.Render("🧠 Uso de memória")

	if m.nativeMemoryErr != nil {
		return title + "\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(m.nativeMemoryErr.Error())
	}
	s := m.nativeMemory
	if s == nil {
		return title + "\n\n" + theme.Help.Render("Carregando…")
	}

	mb := func(kb uint64) string {
		return fmt.Sprintf("%.0f MiB", float64(kb)/1024)
	}

	var sb strings.Builder
	sb.WriteString(title)
	sb.WriteString("\n\n")

	kv := func(k, v string) {
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Width(14).Render(k))
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render(v))
		sb.WriteString("\n")
	}

	sb.WriteString(theme.ScriptScreenAccent.Render("RAM") + "\n")
	kv("Total", mb(s.MemTotalKB))
	kv("Usado*", mb(s.UsedApproxKB()))
	kv("Livre", mb(s.MemFreeKB))
	if s.MemAvailableKB > 0 {
		kv("Disponível", mb(s.MemAvailableKB))
	}
	kv("Compart.", mb(s.ShmemKB))

	sb.WriteString("\n")
	sb.WriteString(theme.ScriptScreenAccent.Render("Swap") + "\n")
	kv("Total", mb(s.SwapTotalKB))
	kv("Livre", mb(s.SwapFreeKB))
	if s.SwapTotalKB > s.SwapFreeKB {
		kv("Usado", mb(s.SwapTotalKB-s.SwapFreeKB))
	}

	sb.WriteString("\n")
	sb.WriteString(theme.Help.Render("* “Usado” é uma estimativa."))

	return sb.String()
}

func humanBytesIEC(n uint64) string {
	f := float64(n)
	switch {
	case n >= 1<<40:
		return fmt.Sprintf("%.2f TiB", f/(1<<40))
	case n >= 1<<30:
		return fmt.Sprintf("%.2f GiB", f/(1<<30))
	case n >= 1<<20:
		return fmt.Sprintf("%.1f MiB", f/(1<<20))
	case n >= 1<<10:
		return fmt.Sprintf("%.1f KiB", f/(1<<10))
	default:
		return fmt.Sprintf("%d B", n)
	}
}

func humanBps(bps float64) string {
	if bps < 1024 {
		return fmt.Sprintf("%.0f B/s", bps)
	}
	if bps < 1024*1024 {
		return fmt.Sprintf("%.1f KiB/s", bps/1024)
	}
	return fmt.Sprintf("%.2f MiB/s", bps/(1024*1024))
}

func renderNativeDiskPanel(m Model) string {
	title := theme.ScriptScreenAccent.Render("💾 Espaço em disco (por mount)")
	if m.nativeDiskErr != nil {
		return title + "\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(m.nativeDiskErr.Error())
	}
	if len(m.nativeDisk) == 0 {
		return title + "\n\n" + theme.Help.Render("Carregando…")
	}
	var sb strings.Builder
	sb.WriteString(title)
	sb.WriteString("\n\n")
	kv := func(k, v string) {
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Width(12).Render(k))
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render(v))
		sb.WriteString("\n")
	}
	for _, d := range m.nativeDisk {
		sb.WriteString(theme.ScriptScreenAccent.Render(d.Mountpoint) + theme.Help.Render("  ("+d.Fstype+")") + "\n")
		kv("Total", humanBytesIEC(d.TotalBytes))
		kv("Usado", humanBytesIEC(d.UsedBytes))
		kv("Livre", humanBytesIEC(d.AvailBytes))
		kv("Uso", fmt.Sprintf("%.1f %%", d.UsePercent))
		sb.WriteString("\n")
	}
	sb.WriteString(theme.Help.Render("/proc/mounts + statfs"))
	return sb.String()
}

func renderNativeLoadPanel(m Model) string {
	title := theme.ScriptScreenAccent.Render("⚙ Carga da CPU (load average)")
	if m.nativeLoadErr != nil {
		return title + "\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(m.nativeLoadErr.Error())
	}
	s := m.nativeLoad
	if s == nil {
		return title + "\n\n" + theme.Help.Render("Carregando…")
	}
	var sb strings.Builder
	sb.WriteString(title)
	sb.WriteString("\n\n")
	kv := func(k, v string) {
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Width(14).Render(k))
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render(v))
		sb.WriteString("\n")
	}
	kv("1 min", fmt.Sprintf("%.2f", s.Load1))
	kv("5 min", fmt.Sprintf("%.2f", s.Load5))
	kv("15 min", fmt.Sprintf("%.2f", s.Load15))
	kv("CPUs (Go)", fmt.Sprintf("%d", s.CPUs))
	if s.CPUs > 0 {
		kv("Load / CPU", fmt.Sprintf("%.2f (1 min)", s.Load1/float64(s.CPUs)))
	}
	kv("Processos", s.Procs)
	sb.WriteString("\n")
	sb.WriteString(theme.Help.Render("Load próximo ao nº de CPUs costuma indicar saturação."))
	return sb.String()
}

func renderNativeNetworkPanel(m Model) string {
	title := theme.ScriptScreenAccent.Render("🌐 Rede (contadores /proc)")
	if m.nativeNetworkErr != nil {
		return title + "\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(m.nativeNetworkErr.Error())
	}
	s := m.nativeNetwork
	if s == nil {
		return title + "\n\n" + theme.Help.Render("Carregando…")
	}
	var sb strings.Builder
	sb.WriteString(title)
	sb.WriteString("\n\n")
	kv := func(k, v string) {
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Width(10).Render(k))
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render(v))
		sb.WriteString("\n")
	}
	for _, iface := range s.Ifaces {
		sb.WriteString(theme.ScriptScreenAccent.Render(iface.Name) + "\n")
		kv("RX", humanBytesIEC(iface.RxBytes))
		kv("TX", humanBytesIEC(iface.TxBytes))
		if r, ok := m.nativeNetRates[iface.Name]; ok {
			kv("Δ RX", humanBps(r.RxBps))
			kv("Δ TX", humanBps(r.TxBps))
		}
		sb.WriteString("\n")
	}
	sb.WriteString(theme.Help.Render("Δ = entre atualizações (~3s); contadores desde o boot."))
	return sb.String()
}

func renderNativeThermalPanel(m Model) string {
	title := theme.ScriptScreenAccent.Render("🌡 Temperatura (sysfs)")
	if m.nativeThermalErr != nil {
		return title + "\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(m.nativeThermalErr.Error())
	}
	s := m.nativeThermal
	if s == nil {
		return title + "\n\n" + theme.Help.Render("Carregando…")
	}
	if len(s.Readings) == 0 {
		return title + "\n\n" + theme.Help.Render("Nenhum sensor em /sys/class/thermal ou hwmon.")
	}
	var sb strings.Builder
	sb.WriteString(title)
	sb.WriteString("\n\n")
	for _, r := range s.Readings {
		line := lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Width(28).Render(r.Label)
		line += lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render(fmt.Sprintf("%.1f °C", r.TempC))
		sb.WriteString(line + "\n")
	}
	sb.WriteString("\n")
	sb.WriteString(theme.Help.Render("thermal_zone*, hwmon"))
	return sb.String()
}

func renderNativeSystemdUserPanel(m Model) string {
	title := theme.ScriptScreenAccent.Render("📋 systemd --user (falhando)")
	if m.nativeSystemdErr != nil {
		return title + "\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(m.nativeSystemdErr.Error())
	}
	s := m.nativeSystemdUser
	if s == nil {
		return title + "\n\n" + theme.Help.Render("Carregando…")
	}
	var sb strings.Builder
	sb.WriteString(title)
	sb.WriteString("\n\n")
	if len(s.Units) == 0 {
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("Nenhuma unidade falhou.") + "\n")
	} else {
		for _, u := range s.Units {
			sb.WriteString(theme.ScriptScreenAccent.Render(u.Unit) + "\n")
			sb.WriteString(theme.Help.Render(u.Sub+" · "+u.Description) + "\n\n")
		}
	}
	sb.WriteString(theme.Help.Render("systemctl --user list-units --state=failed"))
	return sb.String()
}
