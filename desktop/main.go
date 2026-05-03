package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	// Inicjalizacja rygorystycznych parametrów okna
	err := wails.Run(&options.App{
		Title:  "VEXTRO",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		// Transparentne tło wymagane dla warstwy szkła[cite: 1]
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		Frameless:        true, // Tryb Frameless Window[cite: 1]
		Bind: []interface{}{
			app,
		},
		Mac: &mac.Options{
			WebviewIsTransparent: true, // Aktywacja Vibrancy macOS[cite: 1]
			WindowIsTranslucent:  true,
			TitleBar:             mac.TitleBarHidden(),
		},
		Windows: &windows.Options{
			WebviewIsTransparent: true, // Aktywacja Acrylic Windows[cite: 1]
			WindowIsTranslucent:  true,
			BackdropType:         windows.Acrylic,
			DisableWindowIcon:    true,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
