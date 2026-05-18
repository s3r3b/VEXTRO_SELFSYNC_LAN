@echo off
color 0B
echo ===================================================
echo   [APEX OMNI-BUILDER] VEXTRO LAN - BUILD AUTOMATION
echo ===================================================
echo.

:: Tworzenie folderu wyjściowego
if not exist "release" mkdir release

echo [1/3] Budowanie VEXTRO Daemon (Backend Core)...
cd daemon
:: Flaga -H windowsgui ukrywa czarne okno konsoli w gotowym pliku .exe
:: Flagi -s -w zmniejszaja wage pliku (usuwaja tablice symboli)
go build -ldflags="-H windowsgui -s -w" -o ../release/vextro_daemon.exe
cd ..

echo.
echo [2/3] Budowanie VEXTRO Desktop (Wails UI)...
cd desktop
:: Uruchamiamy kompilacje produkcyjna Wails dla Windows
call wails build -clean -platform windows/amd64 -o vextro_desktop.exe
:: Wails domyslnie wrzuca plik do build\bin, przenosimy go do naszego release
move build\bin\vextro_desktop.exe ..\release\ >nul
cd ..

echo.
echo [3/3] Generowanie pliku startowego (Launcher)...
:: Tworzymy prosty skrypt startowy dla uzytkownika wewnatrz folderu release
(
echo @echo off
echo echo Uruchamianie VEXTRO LAN...
echo :: Odpalamy Daemona cicho w tle
echo start "" "vextro_daemon.exe"
echo :: Dajemy mu 1 sekunde na otwarcie portow TCP i mDNS
echo timeout /t 1 /nobreak ^>nul
echo :: Odpalamy warstwe graficzna
echo start "" "vextro_desktop.exe"
echo exit
) > release\start_vextro.bat

echo.
echo ===================================================
echo [SUKCES] Kompilacja zakonczona! 
echo Pelna aplikacja znajduje sie w folderze: /release
echo Uruchamiaj ja za pomoca: release/start_vextro.bat
echo ===================================================
pause