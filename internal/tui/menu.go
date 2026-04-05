package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/JaimeJunr/Homestead/internal/tui/items"
)

func getMainMenuItems(zshCoreInstalled bool) []list.Item {
	out := []list.Item{
		items.MenuItem{Label: "🧹 Limpeza do Sistema", Desc: "Scripts de limpeza e manutenção", Action: menuActionCleanup},
		items.MenuItem{Label: "📊 Monitoramento", Desc: "Painéis nativos: bateria, RAM, disco, CPU, rede, temperatura, systemd --user", Action: menuActionMonitoring},
		items.MenuItem{Label: "🩺 Manutenção / Check-up", Desc: "Só leitura: resumo de alertas + monitores nativos", Action: menuActionCheckup},
		items.MenuItem{Label: "📦 Instaladores", Desc: "IDEs, apps, terminais, utilitários e componentes de sistema", Action: menuActionInstallers},
	}
	if zshCoreInstalled {
		out = append(out, items.MenuItem{Label: "🔧 Plugins e temas Zsh", Desc: "Plugins, temas e .zshrc local", Action: menuActionZshPlugins})
	}
	out = append(out,
		items.MenuItem{Label: "⚙️  Configurar Zsh", Desc: "Repositório de config: backup e migração entre máquinas", Action: menuActionZshRepo},
		items.MenuItem{Label: "⚙️  Configurações", Desc: "Catálogo, tema, caminhos e confirmações", Action: menuActionSettings},
		items.MenuItem{Label: "❌ Sair", Desc: "Fechar Homestead", Action: menuActionQuit},
	)
	return out
}
