package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	DefaultPort = "53535"
)

func main() {
	fmt.Println("[VEXTRO CORE] Uruchamianie sekwencji startowej...")

	InitDB()
	defer CloseDB()
	InitDiscovery()

	listener, err := net.Listen("tcp", ":"+DefaultPort)
	if err != nil {
		fmt.Printf("[FATAL] Nie można uruchomić nasłuchu na porcie %s: %v\n", DefaultPort, err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("[VEXTRO CORE] Nasłuch TCP aktywny na porcie %s. Daemon gotowy.\n", DefaultPort)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go handleConnection(conn)
		}
	}()

	<-sigChan
	fmt.Println("\n[VEXTRO CORE] Otrzymano sygnał zamknięcia...")
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 8192)
	n, err := conn.Read(buffer)
	if err != nil {
		return
	}

	cmd := strings.TrimSpace(string(buffer[:n]))

	if cmd == "IPC_GET_STATUS" {
		conn.Write([]byte(LocalDeviceID))
		return
	}

	if cmd == "IPC_GET_NODES" {
		conn.Write([]byte(GetActiveNodesJSON()))
		return
	}

	if cmd == "IPC_GET_CHAT" {
		conn.Write([]byte(GetChatHistory()))
		return
	}

	// OBSŁUGA WYCHODZĄCEJ WIADOMOŚCI (Z UI)
	if strings.HasPrefix(cmd, "IPC_SEND_MSG:") {
		msgContent := strings.TrimPrefix(cmd, "IPC_SEND_MSG:")
		// 1. Zapisz lokalnie
		AppendChatMessage(LocalDeviceID, msgContent)
		// 2. Roześlij do innych węzłów w LAN [NOWE]
		go BroadcastMessage(LocalDeviceID, msgContent)
		conn.Write([]byte("OK"))
		return
	}

	// OBSŁUGA PRZYCHODZĄCEJ WIADOMOŚCI (OD INNEGO DAEMONA) [NOWE]
	if strings.HasPrefix(cmd, "P2P_RELAY_MSG:") {
		// Format: P2P_RELAY_MSG:SENDER_ID|CONTENT
		data := strings.TrimPrefix(cmd, "P2P_RELAY_MSG:")
		parts := strings.SplitN(data, "|", 2)
		if len(parts) == 2 {
			AppendChatMessage(parts[0], parts[1])
		}
		return
	}

	fmt.Printf("[VEXTRO CORE] [TX/RX] Nierozpoznany sygnał: %s\n", cmd)
}
