package main

import (
	"context"
	"fmt"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	fmt.Println("[VEXTRO WAILS BINDING] Proces graficzny zainicjalizowany.")
}

// GetSystemStatus to testowy most, który wywołamy z Reacta
func (a *App) GetSystemStatus() string {
	return "VEXTRO_DAEMON_WAITING"
}
