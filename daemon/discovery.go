package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"time"

	badger "github.com/dgraph-io/badger/v4"
)

const mDNSIPv4 = "224.0.0.251:5353"

type NodeInfo struct {
	DeviceID string `json:"device_id"`
	Status   string `json:"status"`
	Port     string `json:"port"`
}

var LocalDeviceID string

// InitDiscovery inicjalizuje identyfikację urządzenia i uruchamia usługi mDNS w tle
func InitDiscovery() {
	LocalDeviceID = getOrGenerateDeviceID()
	fmt.Printf("[VEXTRO DISCOVERY] Local DeviceID: %s\n", LocalDeviceID)

	go startMDNSBroadcaster()
	go startMDNSListener()
}

func getOrGenerateDeviceID() string {
	var devID string
	err := DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("system_device_id"))
		if err != nil {
			return err
		}
		valCopy, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		devID = string(valCopy)
		return nil
	})

	// Jeśli klucz nie istnieje (pierwsze uruchomienie), generujemy nowe ID
	if err == badger.ErrKeyNotFound || devID == "" {
		bytes := make([]byte, 4)
		rand.Read(bytes)
		devID = "VXT-" + hex.EncodeToString(bytes)

		err = DB.Update(func(txn *badger.Txn) error {
			return txn.Set([]byte("system_device_id"), []byte(devID))
		})

		if err != nil {
			fmt.Printf("[FATAL] Błąd zapisu DeviceID do BadgerDB: %v\n", err)
		} else {
			fmt.Println("[VEXTRO DISCOVERY] Wygenerowano i zapisano nowy DeviceID.")
		}
	}
	return devID
}

func startMDNSBroadcaster() {
	addr, err := net.ResolveUDPAddr("udp4", mDNSIPv4)
	if err != nil {
		fmt.Printf("[FATAL] Błąd rozwiązywania adresu mDNS: %v\n", err)
		return
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		fmt.Printf("[FATAL] Błąd otwarcia gniazda nadawczego mDNS: %v\n", err)
		return
	}
	defer conn.Close()

	info := NodeInfo{
		DeviceID: LocalDeviceID,
		Status:   "online",
		Port:     DefaultPort,
	}

	payload, _ := json.Marshal(info)

	for {
		_, err := conn.Write(payload)
		if err != nil {
			fmt.Printf("[VEXTRO DISCOVERY] Błąd nadawania mDNS: %v\n", err)
		}
		// Wysyłamy puls (Keep-Alive sieciowy) co 5 sekund
		time.Sleep(5 * time.Second)
	}
}

func startMDNSListener() {
	addr, err := net.ResolveUDPAddr("udp4", mDNSIPv4)
	if err != nil {
		fmt.Printf("[FATAL] Błąd rozwiązywania adresu mDNS dla nasłuchu: %v\n", err)
		return
	}

	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		fmt.Printf("[FATAL] Błąd nasłuchu Multicast UDP: %v\n", err)
		return
	}
	defer conn.Close()

	conn.SetReadBuffer(1024)
	buffer := make([]byte, 1024)

	for {
		n, src, err := conn.ReadFromUDP(buffer)
		if err != nil {
			continue
		}

		var receivedNode NodeInfo
		if err := json.Unmarshal(buffer[:n], &receivedNode); err == nil {
			// Ignorujemy pakiety od samych siebie
			if receivedNode.DeviceID != LocalDeviceID {
				fmt.Printf("[VEXTRO RADAR] Wykryto aktywny węzeł: %s (IP: %s)\n", receivedNode.DeviceID, src.IP.String())
			}
		}
	}
}
