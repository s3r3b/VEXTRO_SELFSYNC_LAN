package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
)

const FileTransferPort = "53536"
const DownloadDir = "vextro_downloads"

// InitFileTransfer tworzy folder pobierania i uruchamia niezależny nasłuch TCP
func InitFileTransfer() {
	os.MkdirAll(DownloadDir, 0755)
	go startFileListener()
}

func startFileListener() {
	listener, err := net.Listen("tcp", ":"+FileTransferPort)
	if err != nil {
		fmt.Printf("[FATAL] Błąd nasłuchu File Transfer na porcie %s: %v\n", FileTransferPort, err)
		return
	}
	defer listener.Close()

	fmt.Printf("[VEXTRO CORE] Nasłuch TRANSFERU PLIKÓW aktywny na porcie %s.\n", FileTransferPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleIncomingFile(conn)
	}
}

func handleIncomingFile(conn net.Conn) {
	defer conn.Close()

	// 1. Oczekujemy nagłówka w formacie: FILE_OFFER:nazwa_pliku.ext
	headerBuf := make([]byte, 1024)
	n, err := conn.Read(headerBuf)
	if err != nil {
		return
	}

	header := strings.TrimSpace(string(headerBuf[:n]))
	if !strings.HasPrefix(header, "FILE_OFFER:") {
		conn.Write([]byte("ERROR: Invalid header"))
		return
	}

	fileName := strings.TrimPrefix(header, "FILE_OFFER:")
	safeFileName := filepath.Base(fileName) // Zabezpieczenie przed Path Traversal
	savePath := filepath.Join(DownloadDir, safeFileName)

	outFile, err := os.Create(savePath)
	if err != nil {
		conn.Write([]byte("ERROR: Cannot create file"))
		return
	}
	defer outFile.Close()

	// 2. Wysyłamy zgodę na transfer
	conn.Write([]byte("ACCEPT"))

	// 3. Odbieramy binarne dane pliku bezpośrednio na dysk
	_, err = io.Copy(outFile, conn)
	if err == nil {
		fmt.Printf("[VEXTRO TRANSFER] Zapisano plik pomyślnie: %s\n", savePath)
	} else {
		fmt.Printf("[VEXTRO TRANSFER] Błąd strumieniowania: %v\n", err)
	}
}

// SendFileToNode łączy się z celem i przesyła wybrany plik z dysku
func SendFileToNode(targetID, filePath string) error {
	nodes := GetActiveNodes()
	targetAddr, exists := nodes[targetID]
	if !exists {
		return fmt.Errorf("węzeł offline lub nieznany: %s", targetID)
	}

	// Zamieniamy port z mDNS (53535) na port transferowy (53536)
	host, _, _ := net.SplitHostPort(targetAddr)
	targetHost := net.JoinHostPort(host, FileTransferPort)

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	conn, err := net.Dial("tcp", targetHost)
	if err != nil {
		return err
	}
	defer conn.Close()

	// 1. Inicjacja Handshake'a
	fileName := filepath.Base(filePath)
	header := fmt.Sprintf("FILE_OFFER:%s", fileName)
	conn.Write([]byte(header))

	// 2. Czekamy na zezwolenie z węzła docelowego
	ackBuf := make([]byte, 32)
	n, err := conn.Read(ackBuf)
	if err != nil || string(ackBuf[:n]) != "ACCEPT" {
		return fmt.Errorf("transfer odrzucony lub błąd handshaku")
	}

	// 3. Strumieniowanie zawartości
	_, err = io.Copy(conn, file)
	return err
}
