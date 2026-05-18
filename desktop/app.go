package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
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

func (a *App) GetSystemStatus() string {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:53535", 2*time.Second)
	if err != nil {
		return "DAEMON OFFLINE"
	}
	defer conn.Close()

	conn.Write([]byte("IPC_GET_STATUS"))
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return "BŁĄD ODCZYTU IPC"
	}
	return string(buffer[:n])
}

func (a *App) GetActiveNodes() string {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:53535", 2*time.Second)
	if err != nil {
		return "{}"
	}
	defer conn.Close()

	conn.Write([]byte("IPC_GET_NODES"))
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buffer := make([]byte, 8192)
	n, err := conn.Read(buffer)
	if err != nil {
		return "{}"
	}
	return string(buffer[:n])
}

func (a *App) GetChatHistory() string {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:53535", 2*time.Second)
	if err != nil {
		return "[]"
	}
	defer conn.Close()

	conn.Write([]byte("IPC_GET_CHAT"))
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buffer := make([]byte, 65536)
	n, err := conn.Read(buffer)
	if err != nil {
		return "[]"
	}
	return string(buffer[:n])
}

func (a *App) SendChatMessage(message string) string {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:53535", 2*time.Second)
	if err != nil {
		return "ERROR_DIAL"
	}
	defer conn.Close()

	payload := fmt.Sprintf("IPC_SEND_MSG:%s", message)
	conn.Write([]byte(payload))

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return "ERROR_READ"
	}
	return string(buffer[:n])
}

// [NOWE] SelectAndSendFile - Otwiera natywne okno, pobiera ścieżkę i zleca transfer
func (a *App) SelectAndSendFile(targetID string) string {
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Wybierz plik do wysłania przez LAN",
	})

	if err != nil || filePath == "" {
		return "CANCELLED"
	}

	conn, err := net.DialTimeout("tcp", "127.0.0.1:53535", 2*time.Second)
	if err != nil {
		return "ERROR_DIAL"
	}
	defer conn.Close()

	payload := fmt.Sprintf("IPC_SEND_FILE:%s|%s", targetID, filePath)
	conn.Write([]byte(payload))

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return "ERROR_READ"
	}
	return string(buffer[:n])
}
