package main
import (
	"fmt"
	"os"
	"path/filepath"

	badger "github.com/dgraph-io/badger/v4"
)

var DB *badger.DB

func InitDB() {
	// Pobranie ścieżki do katalogu domowego
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("[FATAL] Nie można zlokalizować katalogu domowego: %v\n", err)
		os.Exit(1)
	}

	// Konstrukcja ścieżki: ~/Documents/VEXTRO_SELFsynclan/db
	dbPath := filepath.Join(homeDir, "Documents", "VEXTRO_SELFsynclan", "db")

	if err := os.MkdirAll(dbPath, 0755); err != nil {
		fmt.Printf("[FATAL] Nie można utworzyć katalogu bazy danych: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("[VEXTRO CORE] Inicjalizacja BadgerDB w: %s\n", dbPath)

	opts := badger.DefaultOptions(dbPath)
	opts.Logger = nil 

	db, err := badger.Open(opts)
	if err != nil {
		fmt.Printf("[FATAL] Błąd otwierania BadgerDB: %v\n", err)
		os.Exit(1)
	}

	DB = db
	fmt.Println("[VEXTRO CORE] Silnik BadgerDB I/O w pełni operacyjny.")
}

func CloseDB() {
	if DB != nil {
		fmt.Println("[VEXTRO CORE] Zamykanie strumieni BadgerDB...")
		DB.Close()
	}
}
