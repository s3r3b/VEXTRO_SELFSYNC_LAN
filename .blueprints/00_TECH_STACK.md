[!!!WARNING!!!!]

[FORBIDDEN]
[IT IS TOTALLY PROHIBITED TO EDIT THIS FILE]



[ALLOWED]
[YOU CAN COME BACK HERE FOR INFORMATION AND READ IT AS MUCH AS YOU WANT]

# 🗂️ [TECH STACK - RELIGION GRADE] (Kwiecień 2026)

## 1. ZATWIERDZONY STACK & WERSJE
- **CORE/DAEMON**: Go v1.26.x | BadgerDB v4.x | Czyste TCP/UDP (Brak Szyfrowania)
- **ŚRODOWISKO**: Node.js v24.x LTS
- **DESKTOP**: Wails v2.12.x | React v18.3.x | TypeScript v5.4.x | Vite v5.x
- **MOBILE**: React Native v0.85.x (Hermes + New Arch) | Target: Android 16+ (API 36+) | React Navigation v6.x
- **UI/STYLING**: Tailwind CSS v4.2.x | NativeWind v5.x

## 2. GŁÓWNE FUNKCJE (Smoked Glass / LAN P2P)
- **Identyfikacja**: Wektorowe ikony, stałe DeviceID.
- **Topologia**: mDNS Discovery, brak serwera centralnego.
- **Czat**: Lokalny zapis (.txt dla wiadomości, natywny dla plików). Jeden współdzielony widok.
- **Transfer**: <5GB surowe TCP | >5GB chunking + hash + wznawianie.
- **Desktop w Tle**: Daemon zminimalizowany do Tray'a z nasłuchem TCP (Keep-Alive).
- **Zarządzanie Stanem**: Sync Request wysyłany po przebudzeniu z uśpienia.
EOF