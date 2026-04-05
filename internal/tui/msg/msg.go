package msg

import (
	"github.com/JaimeJunr/Homestead/internal/domain/interfaces"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/preferences"
	"github.com/JaimeJunr/Homestead/internal/monitoring"
)

type Progress interfaces.InstallProgress

type InstallComplete struct {
	Err error
}

type ZshCoreInstalled struct {
	Installed bool
}

type ZshApplyResult struct {
	Err error
}

type ZshApplyReturnToMenu struct{}

type ScriptCaptured struct {
	Output string
	Err    error
}

type ScriptExecFinished struct {
	Err error
}

type URLActionDone struct {
	Err  error
	Verb string
}

type ClearKeyboardToast struct{}

type CatalogFetched struct {
	Err error
	Ok  bool
}

type NativeMonitorReload struct {
	Kind        string
	Battery     *monitoring.BatterySnapshot
	Memory      *monitoring.MemorySnapshot
	Disk        []monitoring.DiskMount
	Load        *monitoring.LoadSnapshot
	Network     *monitoring.NetworkSnapshot
	Thermal     *monitoring.ThermalSnapshot
	SystemdUser *monitoring.SystemdUserFailedSnapshot
	Err         error
}

type NativeMonitorTick struct{}

type SettingsSaved struct {
	Prefs preferences.Preferences
}

type SettingsCancelled struct{}
