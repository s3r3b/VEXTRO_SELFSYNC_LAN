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

	// 1. Inicjalizacja bazy danych
	InitDB()
	defer CloseDB()

	// 1.5. Inicjalizacja mDNS i Identyfikacji
	InitDiscovery()

	// 2. Inicjalizacja gniazda TCP
	listener, err := net.Listen("tcp", ":"+DefaultPort)
	if err != nil {
		fmt.Printf("[FATAL] Nie można uruchomić nasłuchu na porcie %s: %v\n", DefaultPort, err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("[VEXTRO CORE] Nasłuch TCP aktywny na porcie %s. Daemon gotowy.\n", DefaultPort)

	// 3. Graceful Shutdown
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
	fmt.Println("\n[VEXTRO CORE] Otrzymano sygnał zamknięcia. Wykonywanie zrzutu pamięci...")
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return
	}

	cmd := strings.TrimSpace(string(buffer[:n]))

	// Mikro-router IPC dla połączeń lokalnych (Frontend Wails <-> Daemon)
	if cmd == "IPC_GET_STATUS" {
		response := fmt.Sprintf("%s", LocalDeviceID)
		conn.Write([]byte(response))
		return
	}

	fmt.Printf("[VEXTRO CORE] [TX/RX] Nierozpoznany sygnał lub payload TCP od: %s\n", conn.RemoteAddr().String())
}
