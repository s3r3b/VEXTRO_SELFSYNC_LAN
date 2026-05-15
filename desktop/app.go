package main

import (
	"context"
	"fmt"
	"net"
	"time"
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

// GetSystemStatus łączy się z lokalnym portem Daemona (53535) i pobiera DeviceID
func (a *App) GetSystemStatus() string {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:53535", 2*time.Second)
	if err != nil {
		return "DAEMON OFFLINE (BRAK POŁĄCZENIA)"
	}
	defer conn.Close()

	_, err = conn.Write([]byte("IPC_GET_STATUS"))
	if err != nil {
		return "BŁĄD ZAPISU IPC"
	}

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return "BŁĄD ODCZYTU IPC"
	}

	return string(buffer[:n])
}

// GetActiveNodes odpytuje Daemona o zrzucony do JSONa stan ActiveNodes (Radar mDNS)
func (a *App) GetActiveNodes() string {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:53535", 2*time.Second)
	if err != nil {
		return "[]"
	}
	defer conn.Close()

	_, err = conn.Write([]byte("IPC_GET_NODES"))
	if err != nil {
		return "[]"
	}

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buffer := make([]byte, 8192) // Większy bufor dla tablicy wielu węzłów JSON
	n, err := conn.Read(buffer)
	if err != nil {
		return "[]"
	}

	return string(buffer[:n])
}
