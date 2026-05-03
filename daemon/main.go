package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	DefaultPort = "53535"
)

func main() {
	fmt.Println("[VEXTRO CORE] Uruchamianie sekwencji startowej...")

	// 1. Inicjalizacja bazy danych w dokumentach usera
	InitDB()
	defer CloseDB()

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
	fmt.Printf("[VEXTRO CORE] [TX/RX] Odrzucono/Przyjęto sygnał od: %s\n", conn.RemoteAddr().String())
}
