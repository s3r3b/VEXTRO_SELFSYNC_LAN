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
	InitFileTransfer()

	listener, err := net.Listen("tcp", ":"+DefaultPort)
	if err != nil {
		fmt.Printf("[FATAL] Nie można uruchomić nasłuchu na porcie %s: %v\n", DefaultPort, err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("[VEXTRO CORE] Nasłuch TCP (IPC/CHAT) aktywny na porcie %s. Daemon gotowy.\n", DefaultPort)

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

	if strings.HasPrefix(cmd, "IPC_SEND_MSG:") {
		msgContent := strings.TrimPrefix(cmd, "IPC_SEND_MSG:")
		AppendChatMessage(LocalDeviceID, msgContent)
		go BroadcastMessage(LocalDeviceID, msgContent)
		conn.Write([]byte("OK"))
		return
	}

	if strings.HasPrefix(cmd, "IPC_SEND_FILE:") {
		data := strings.TrimPrefix(cmd, "IPC_SEND_FILE:")
		parts := strings.SplitN(data, "|", 2)
		if len(parts) == 2 {
			go func() {
				err := SendFileToNode(parts[0], parts[1])
				if err != nil {
					fmt.Printf("[VEXTRO CORE] Błąd transferu: %v\n", err)
				}
			}()
			conn.Write([]byte("TRANSFER_STARTED"))
		} else {
			conn.Write([]byte("ERROR_MALFORMED"))
		}
		return
	}

	// [ZMODYFIKOWANE] Odbieranie P2P - dodano powiadomienie
	if strings.HasPrefix(cmd, "P2P_RELAY_MSG:") {
		data := strings.TrimPrefix(cmd, "P2P_RELAY_MSG:")
		parts := strings.SplitN(data, "|", 2)
		if len(parts) == 2 {
			AppendChatMessage(parts[0], parts[1])
			// Wyświetlamy powiadomienie o nowej wiadomości
			NotifyUser("VEXTRO: Nowa wiadomość w LAN", fmt.Sprintf("Od %s: %s", parts[0], parts[1]))
		}
		return
	}

	if cmd == "IPC_SHUTDOWN" {
		fmt.Println("[VEXTRO CORE] Otrzymano krytyczny sygnał IPC_SHUTDOWN z zasobnika. Zamykanie...")
		conn.Write([]byte("SHUTTING_DOWN"))
		os.Exit(0)
		return
	}

	fmt.Printf("[VEXTRO CORE] [TX/RX] Nierozpoznany sygnał: %s\n", cmd)
}
