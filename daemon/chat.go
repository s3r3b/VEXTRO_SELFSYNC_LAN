package main

import (
	"encoding/json"
	"os"
	"strings"
	"sync"
	"time"
)

const ChatFileName = "shared_chat.txt"

var chatMutex sync.Mutex

// ChatMessage definiuje strukturę pojedynczej wiadomości w historii
type ChatMessage struct {
	Timestamp string `json:"timestamp"`
	SenderID  string `json:"senderId"`
	Content   string `json:"content"`
}

// AppendChatMessage dodaje nową wiadomość do płaskiego pliku txt
func AppendChatMessage(senderID, message string) error {
	chatMutex.Lock()
	defer chatMutex.Unlock()

	file, err := os.OpenFile(ChatFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	timestamp := time.Now().Format(time.RFC3339)
	// Format zapisu: TIMESTAMP|SENDER_ID|MESSAGE
	line := timestamp + "|" + senderID + "|" + message + "\n"
	_, err = file.WriteString(line)
	return err
}

// GetChatHistory odczytuje plik i zwraca ustandaryzowaną listę w formacie JSON
func GetChatHistory() string {
	chatMutex.Lock()
	defer chatMutex.Unlock()

	data, err := os.ReadFile(ChatFileName)
	if err != nil {
		return "[]" // Jeśli plik nie istnieje, zwracamy pustą tablicę JSON
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	var messages []ChatMessage

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 3)
		if len(parts) == 3 {
			messages = append(messages, ChatMessage{
				Timestamp: parts[0],
				SenderID:  parts[1],
				Content:   parts[2],
			})
		}
	}

// BroadcastMessage przesyła wiadomość do wszystkich aktywnych węzłów w sieci LAN
func BroadcastMessage(senderID, content string) {
	nodes := GetActiveNodes() // Pobiera mapę węzłów z discovery.go
	payload := fmt.Sprintf("P2P_RELAY_MSG:%s|%s", senderID, content)

	for id, addr := range nodes {
		if id == LocalDeviceID {
			continue
		}

		go func(address string) {
			// Próba połączenia z innym Daemonem (port 53535)
			// Uwaga: adres z mDNS zawiera port discovery, musimy użyć portu Daemona
			host, _, _ := net.SplitHostPort(address)
			target := net.JoinHostPort(host, DefaultPort)

			conn, err := net.DialTimeout("tcp", target, 2*time.Second)
			if err != nil {
				return
			}
			defer conn.Close()
			conn.Write([]byte(payload))
		}(addr)
	}
}

	// Pakujemy przetworzoną historię w JSON, by ułatwić życie Frontendowi (R2)
	jsonBytes, err := json.Marshal(messages)
	if err != nil {
		return "[]"
	}

	return string(jsonBytes)
}
