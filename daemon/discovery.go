package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	badger "github.com/dgraph-io/badger/v4"
)

const mDNSIPv4 = "224.0.0.251:5353"

type NodeInfo struct {
	DeviceID string `json:"device_id"`
	Status   string `json:"status"`
	Port     string `json:"port"`
	IP       string `json:"ip"`
	LastSeen int64  `json:"last_seen"`
}

var (
	LocalDeviceID string
	ActiveNodes   = make(map[string]NodeInfo)
	nodesMutex    sync.RWMutex
)

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

	for {
		// Aktualizujemy payload za każdym razem, aby mieć świeży timestamp (choć nie używamy go w nadajniku)
		payload, _ := json.Marshal(info)
		_, err := conn.Write(payload)
		if err != nil {
			fmt.Printf("[VEXTRO DISCOVERY] Błąd nadawania mDNS: %v\n", err)
		}
		// Puls mDNS co 5 sekund
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
			if receivedNode.DeviceID != LocalDeviceID {
				receivedNode.IP = src.IP.String()
				receivedNode.LastSeen = time.Now().Unix()

				nodesMutex.Lock()
				_, exists := ActiveNodes[receivedNode.DeviceID]
				ActiveNodes[receivedNode.DeviceID] = receivedNode
				nodesMutex.Unlock()

				if !exists {
					fmt.Printf("[VEXTRO RADAR] Węzeł '%s' wszedł do strefy LAN (IP: %s)\n", receivedNode.DeviceID, receivedNode.IP)
				}
			}
		}
	}
}

// GetActiveNodesJSON zwraca listę aktywnych węzłów. Odcina te, od których nie było pulsu od ponad 15 sekund.
func GetActiveNodesJSON() string {
	nodesMutex.RLock()
	defer nodesMutex.RUnlock()

	var nodeList []NodeInfo
	now := time.Now().Unix()

	for _, node := range ActiveNodes {
		// Timeout: 15 sekund. Jeśli węzeł zniknie z sieci, znika z radaru.
		if now-node.LastSeen <= 15 {
			nodeList = append(nodeList, node)
		}
	}

	// Jeśli tablica jest pusta, upewniamy się, że zwrócimy [] zamiast null
	if nodeList == nil {
		nodeList = make([]NodeInfo, 0)
	}

	data, _ := json.Marshal(nodeList)
	return string(data)
}
