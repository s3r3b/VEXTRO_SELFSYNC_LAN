package main

import (
	"context"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	// [NOWE] Odpalamy zasobnik systemowy w niezależnej gorutynie
	go RunTray(app)

	err := wails.Run(&options.App{
		Title:  "VEXTRO LAN",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 15, G: 15, B: 17, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
		},
		Bind: []interface{}{
			app,
		},
		// [KRYTYCZNE] Zapobiega natychmiastowej śmierci procesu po zamknięciu okna
		HideWindowOnClose: true,
	})

	if err != nil {
		println("Błąd uruchamiania Wails:", err.Error())
	}
}
