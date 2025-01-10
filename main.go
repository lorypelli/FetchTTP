package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)
//go:embed frontend/dist
var assets embed.FS

func main() {
	app := NewApp()
	wails.Run(&options.App{
		Title:  "FetchTTP",
		MinWidth:  630,
		MinHeight: 780,
		EnableDefaultContextMenu: false,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
	})
}
