cat << 'EOF' > .blueprints/01_ROADMAP_INIT.md
# 🗺️ [VEXTRO SELFsynclan] - ROADMAP: FAZA 1 (Inicjalizacja Infrastruktury)

**CEL:** Ustanowienie surowego środowiska developerskiego, struktury monorepo oraz fundamentów dla demona Go i aplikacji okienkowych.

- [ ] **KROK 3:** Konfiguracja przestrzeni roboczej Go Daemon (`go mod init vextro-daemon`) w nowym katalogu `daemon`.
- [ ] **KROK 4:** Implementacja absolutnego szkieletu Daemona Go (`main.go` z obsługą zasobnika systemowego / tła).
- [ ] **KROK 5:** Inicjalizacja lokalnej instancji BadgerDB (zapis/odczyt plików w systemie).
- [ ] **KROK 6:** Inicjalizacja środowiska Wails dla Desktop UI (`desktop`).
- [ ] **KROK 7:** Inicjalizacja środowiska React Native CLI dla Mobile UI (`mobile`).
EOF