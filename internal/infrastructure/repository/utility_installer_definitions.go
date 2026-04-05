package repository

// utilityInstallerDef describes a Homestead utility script exposed as PackageCategoryUtilities.
type utilityInstallerDef struct {
	ID, Name, Description, Path string
	RequiresSudo                bool
	ProjectURL                  string
}

// utilityInstallerDefinitions is the single source for utilitários: scripts no repositório + pacotes em Instaladores.
func utilityInstallerDefinitions() []utilityInstallerDef {
	return []utilityInstallerDef{
		{ID: "util-arch-update", Name: "Arch-Update", Description: "Atualizações Arch na bandeja (Chaotic-AUR + pacotes)", Path: "scripts/utilities/archupd.sh", RequiresSudo: true, ProjectURL: "https://github.com/Antiz96/arch-update"},
		{ID: "util-bitwarden", Name: "Bitwarden", Description: "Gestor de senhas", Path: "scripts/utilities/fp-user/com.bitwarden.desktop.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/com.bitwarden.desktop"},
		{ID: "util-bottles", Name: "Bottles", Description: "Executar apps Windows com Wine", Path: "scripts/utilities/fp-user/com.usebottles.bottles.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/com.usebottles.bottles"},
		{ID: "util-boxadv", Name: "Distrobox-Adv BR", Description: "Distrobox para certificado digital (advogados) + DistroShelf", Path: "scripts/utilities/boxadv.sh", RequiresSudo: true, ProjectURL: "https://github.com/89luca89/distrobox"},
		{ID: "util-brave", Name: "Brave", Description: "Navegador Brave", Path: "scripts/utilities/fp-user/com.brave.Browser.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/com.brave.Browser"},
		{ID: "util-bazaar", Name: "Bazaar", Description: "Loja de apps e jogos · instalação sistema", Path: "scripts/utilities/fp-system/io.github.kolunmi.Bazaar.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/io.github.kolunmi.Bazaar"},
		{ID: "util-cloudflare-warp", Name: "Cloudflare WARP", Description: "Cliente WARP (repositório oficial)", Path: "scripts/utilities/cloudflare-warp.sh", RequiresSudo: true, ProjectURL: "https://developers.cloudflare.com/cloudflare-one/connections/connect-devices/warp/download-warp/"},
		{ID: "util-cryptomator", Name: "Cryptomator", Description: "Criptografia de arquivos", Path: "scripts/utilities/fp-user/org.cryptomator.Cryptomator.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/org.cryptomator.Cryptomator"},
		{ID: "util-distroshelf", Name: "DistroShelf", Description: "Podman, Distrobox e DistroShelf", Path: "scripts/utilities/distroshelf.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/com.ranfdev.DistroShelf"},
		{ID: "util-easyeffects", Name: "Easy Effects", Description: "Efeitos de áudio PipeWire · instalação sistema", Path: "scripts/utilities/fp-system/com.github.wwmm.easyeffects.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/com.github.wwmm.easyeffects"},
		{ID: "util-extension-manager", Name: "Extension Manager", Description: "Gestor de extensões GNOME", Path: "scripts/utilities/fp-user/com.mattjakeman.ExtensionManager.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/com.mattjakeman.ExtensionManager"},
		{ID: "util-expressvpn", Name: "ExpressVPN", Description: "Instalador oficial ExpressVPN (Linux universal)", Path: "scripts/utilities/express-vpn.sh", RequiresSudo: true, ProjectURL: "https://www.expressvpn.com/support/latest/linux/"},
		{ID: "util-f3", Name: "f3", Description: "Testar cartões SD/USB contra fraude de capacidade", Path: "scripts/utilities/f3.sh", RequiresSudo: true, ProjectURL: "https://github.com/AltraMayor/f3"},
		{ID: "util-flatseal", Name: "Flatseal", Description: "Permissões de apps Flatpak", Path: "scripts/utilities/fp-user/com.github.tchx84.Flatseal.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/com.github.tchx84.Flatseal"},
		{ID: "util-gearlever", Name: "Gear Lever", Description: "Gestor de AppImages", Path: "scripts/utilities/fp-user/it.mijorus.gearlever.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/it.mijorus.gearlever"},
		{ID: "util-gnome-tweaks", Name: "GNOME Tweaks", Description: "Ajustes avançados do GNOME", Path: "scripts/utilities/gnome-tweaks.sh", RequiresSudo: true, ProjectURL: "https://wiki.gnome.org/Apps/Tweaks"},
		{ID: "util-handbrake", Name: "HandBrake", Description: "Codificador de vídeo", Path: "scripts/utilities/fp-user/fr.handbrake.ghb.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/fr.handbrake.ghb"},
		{ID: "util-haguichi", Name: "Haguichi", Description: "Cliente gráfico Hamachi", Path: "scripts/utilities/fp-user/com.github.ztefn.haguichi.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/com.github.ztefn.haguichi"},
		{ID: "util-hiddify", Name: "Hiddify", Description: "Cliente Hiddify (AppImage em ~/.local)", Path: "scripts/utilities/hiddify.sh", RequiresSudo: false, ProjectURL: "https://github.com/hiddify/hiddify-next"},
		{ID: "util-input-remapper", Name: "Input Remapper", Description: "Remapear teclado e mouse (input-remapper)", Path: "scripts/utilities/input-remapper.sh", RequiresSudo: true, ProjectURL: "https://github.com/sezanzeb/input-remapper"},
		{ID: "util-ivpn", Name: "IVPN", Description: "Cliente IVPN (repositório oficial)", Path: "scripts/utilities/ivpn.sh", RequiresSudo: true, ProjectURL: "https://www.ivpn.net/"},
		{ID: "util-keepassxc", Name: "KeePassXC", Description: "Gestor de senhas KeePassXC", Path: "scripts/utilities/fp-user/org.keepassxc.KeePassXC.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/org.keepassxc.KeePassXC"},
		{ID: "util-lact", Name: "LACT", Description: "Controle AMD GPU · instalação sistema", Path: "scripts/utilities/fp-system/io.github.ilya_zlobintsev.LACT.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/io.github.ilya_zlobintsev.LACT"},
		{ID: "util-librewolf", Name: "LibreWolf", Description: "Navegador LibreWolf", Path: "scripts/utilities/fp-user/io.gitlab.librewolf-community.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/io.gitlab.librewolf-community"},
		{ID: "util-logseq", Name: "Logseq", Description: "Bloco de notas em grafos", Path: "scripts/utilities/fp-user/com.logseq.Logseq.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/com.logseq.Logseq"},
		{ID: "util-mission-center", Name: "Mission Center", Description: "Monitor de sistema", Path: "scripts/utilities/fp-user/io.missioncenter.MissionCenter.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/io.missioncenter.MissionCenter"},
		{ID: "util-mullvad-browser", Name: "Mullvad Browser", Description: "Navegador Mullvad", Path: "scripts/utilities/fp-user/net.mullvad.MullvadBrowser.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/net.mullvad.MullvadBrowser"},
		{ID: "util-mullvad-vpn", Name: "Mullvad VPN", Description: "Cliente Mullvad VPN (repositório oficial)", Path: "scripts/utilities/mullvad-vpn.sh", RequiresSudo: true, ProjectURL: "https://mullvad.net/download/linux/"},
		{ID: "util-nerd-fonts", Name: "Nerd Fonts (JetBrains)", Description: "JetBrains Mono em ~/.local/share/fonts", Path: "scripts/utilities/nerd-fonts.sh", RequiresSudo: false, ProjectURL: "https://github.com/ryanoasis/nerd-fonts"},
		{ID: "util-nordvpn", Name: "NordVPN", Description: "Cliente NordVPN (script oficial ou AUR)", Path: "scripts/utilities/nord-vpn.sh", RequiresSudo: true, ProjectURL: "https://nordvpn.com/download/linux/"},
		{ID: "util-obs-studio", Name: "OBS Studio", Description: "OBS + plugin PipeWire e overrides", Path: "scripts/utilities/obs-studio.sh", RequiresSudo: true, ProjectURL: "https://obsproject.com/"},
		{ID: "util-openlinkhub", Name: "OpenLinkHub", Description: "Periféricos Corsair (repositório / COPR / AUR)", Path: "scripts/utilities/openlinkhub.sh", RequiresSudo: true, ProjectURL: "https://github.com/EvanMulawski/OpenLinkHub"},
		{ID: "util-openrgb", Name: "OpenRGB", Description: "Iluminação RGB", Path: "scripts/utilities/fp-user/org.openrgb.OpenRGB.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/org.openrgb.OpenRGB"},
		{ID: "util-openrazer", Name: "OpenRazer", Description: "Drivers Razer + Polychromatic · instalação sistema", Path: "scripts/utilities/openrazer.sh", RequiresSudo: true, ProjectURL: "https://openrazer.github.io/"},
		{ID: "util-openrazer-ublue", Name: "OpenRazer (ujust)", Description: "OpenRazer em imagens Universal Blue (ujust)", Path: "scripts/utilities/openrazer-ublue.sh", RequiresSudo: true, ProjectURL: "https://openrazer.github.io/"},
		{ID: "util-oversteer", Name: "Oversteer", Description: "Ajuste de volantes", Path: "scripts/utilities/fp-user/io.github.berarma.Oversteer.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/io.github.berarma.Oversteer"},
		{ID: "util-peazip", Name: "PeaZip", Description: "Arquivador", Path: "scripts/utilities/fp-user/io.github.peazip.PeaZip.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/io.github.peazip.PeaZip"},
		{ID: "util-pika-backup", Name: "Pika Backup", Description: "Backups com Borg", Path: "scripts/utilities/fp-user/org.gnome.World.PikaBackup.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/org.gnome.World.PikaBackup"},
		{ID: "util-piper", Name: "Piper", Description: "Configurar mouse (ratbagd + Flatpak sistema)", Path: "scripts/utilities/piper.sh", RequiresSudo: true, ProjectURL: "https://github.com/libratbag/piper"},
		{ID: "util-proton-vpn", Name: "Proton VPN", Description: "Cliente Proton VPN", Path: "scripts/utilities/fp-user/com.protonvpn.www.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/com.protonvpn.www"},
		{ID: "util-qpwgraph", Name: "qpwgraph", Description: "Patchbay PipeWire", Path: "scripts/utilities/fp-user/org.rncbc.qpwgraph.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/org.rncbc.qpwgraph"},
		{ID: "util-rclone-ui", Name: "Rclone UI", Description: "Interface para rclone", Path: "scripts/utilities/fp-user/com.rcloneui.RcloneUI.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/com.rcloneui.RcloneUI"},
		{ID: "util-s3drive", Name: "S3Drive", Description: "Cliente S3 / armazenamento", Path: "scripts/utilities/fp-user/io.kapsa.drive.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/io.kapsa.drive"},
		{ID: "util-sirikali", Name: "SiriKali", Description: "Volumes cifrados", Path: "scripts/utilities/fp-user/io.github.mhogomchungu.sirikali.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/io.github.mhogomchungu.sirikali"},
		{ID: "util-solaar", Name: "Solaar", Description: "Unifying / Bolt Logitech", Path: "scripts/utilities/solaar.sh", RequiresSudo: true, ProjectURL: "https://github.com/pwr-solaar/Solaar"},
		{ID: "util-stream-controller", Name: "Stream Controller", Description: "Controles para streaming", Path: "scripts/utilities/fp-user/com.core447.StreamController.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/com.core447.StreamController"},
		{ID: "util-surfshark", Name: "Surfshark VPN", Description: "Cliente Surfshark", Path: "scripts/utilities/fp-user/com.surfshark.Surfshark.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/com.surfshark.Surfshark"},
		{ID: "util-udev-hidraw", Name: "udev HID raw", Description: "Regra udev 0666 para USB (menu no terminal)", Path: "scripts/utilities/udev-hidraw.sh", RequiresSudo: true, ProjectURL: "https://github.com/JaimeJunr/Homestead"},
		{ID: "util-ungoogled-chromium", Name: "Ungoogled Chromium", Description: "Chromium sem Google", Path: "scripts/utilities/fp-user/io.github.ungoogled_software.ungoogled_chromium.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/io.github.ungoogled_software.ungoogled_chromium"},
		{ID: "util-vlc", Name: "VLC", Description: "Multimedia (repositórios nativos + codecs quando aplicável)", Path: "scripts/utilities/vlc.sh", RequiresSudo: true, ProjectURL: "https://www.videolan.org/vlc/"},
		{ID: "util-warehouse", Name: "Warehouse", Description: "Gestor de apps Flatpak", Path: "scripts/utilities/fp-user/io.github.flattool.Warehouse.sh", RequiresSudo: true, ProjectURL: "https://flathub.org/apps/io.github.flattool.Warehouse"},
		{ID: "util-waydroid", Name: "Waydroid", Description: "Android no Linux (Wayland)", Path: "scripts/utilities/waydroid.sh", RequiresSudo: true, ProjectURL: "https://waydro.id/"},
		{ID: "util-windscribe", Name: "Windscribe VPN", Description: "Cliente Windscribe (pacote oficial)", Path: "scripts/utilities/windscribe-vpn.sh", RequiresSudo: true, ProjectURL: "https://windscribe.com/unlimited/linux/"},
		{ID: "util-wireguard", Name: "WireGuard", Description: "Ferramentas WireGuard (pacotes da distro)", Path: "scripts/utilities/wireguard.sh", RequiresSudo: true, ProjectURL: "https://www.wireguard.com/"},
	}
}
