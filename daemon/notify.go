package main

import (
	"fmt"

	"github.com/gen2brain/beeep"
)

// NotifyUser wysyła natywne powiadomienie systemowe (np. Windows Toast)
func NotifyUser(title, message string) {
	// Trzeci argument to ścieżka do ikony - w MVP zostawiamy puste dla domyślnej ikony
	err := beeep.Notify(title, message, "")
	if err != nil {
		fmt.Printf("[VEXTRO NOTIFY] Błąd wyświetlania powiadomienia: %v\n", err)
	}
}
