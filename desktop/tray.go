package main

import (
	_ "embed"
	"net"
	"os"

	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Wails automatycznie generuje ten plik; użyjemy go jako ikony w trayu
//
//go:embed build/appicon.png
var iconData []byte

// RunTray odpala logikę zasobnika systemowego
func RunTray(app *App) {
	systray.Run(func() { onReady(app) }, onExit)
}

func onReady(app *App) {
	systray.SetIcon(iconData)
	systray.SetTitle("VEXTRO")
	systray.SetTooltip("VEXTRO LAN - Nasłuch aktywny")

	mShow := systray.AddMenuItem("Pokaż Panel LAN", "Przywraca okno czatu i transferu")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Zakończ całkowicie", "Zamyka UI oraz zatrzymuje Daemona w tle")

	go func() {
		for {
			select {
			case <-mShow.ClickedCh:
				// Przywrócenie i skupienie okna UI na ekranie
				if app.ctx != nil {
					runtime.WindowShow(app.ctx)
				}
			case <-mQuit.ClickedCh:
				// 1. ZABIJ DAEMONA (Wysłanie komendy na TCP)
				conn, err := net.Dial("tcp", "127.0.0.1:53535")
				if err == nil {
					conn.Write([]byte("IPC_SHUTDOWN"))
					conn.Close()
				}

				// 2. WYŁĄCZ ZASOBNIK
				systray.Quit()

				// 3. ZABIJ PROCES GŁÓWNY WAILSA
				if app.ctx != nil {
					runtime.Quit(app.ctx)
				}
				os.Exit(0)
			}
		}
	}()
}

func onExit() {
	// Czyszczenie cyklu życia
}
