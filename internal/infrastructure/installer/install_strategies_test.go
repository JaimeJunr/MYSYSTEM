package installer

import (
	"testing"

	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

func TestExpandInstallCmd(t *testing.T) {
	got := expandInstallCmd("sudo dpkg -i {{download_path}}", "/tmp/pkg.deb")
	want := "sudo dpkg -i /tmp/pkg.deb"
	if got != want {
		t.Fatalf("expandInstallCmd = %q, want %q", got, want)
	}
}

func TestNewStrategyMapCoversAllKinds(t *testing.T) {
	m := newStrategyMap(defaultStrategies())
	for _, k := range types.AllPackageInstallKinds() {
		if m[k] == nil {
			t.Fatalf("missing strategy for kind %q", k)
		}
	}
}
