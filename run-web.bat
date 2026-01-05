@echo off
REM Script pour lancer l'interface web sur Windows
REM IA-ATOMIQUE Web Server

setlocal enabledelayexpansion

echo.
echo ╔════════════════════════════════════════╗
echo ║  IA-ATOMIQUE - Serveur Web             ║
echo ╚════════════════════════════════════════╝
echo.

REM Vérifier si le binaire existe
if not exist "programme.exe" (
    echo ❌ ERREUR: programme.exe n'existe pas
    echo.
    echo Exécute d'abord: build.bat
    echo.
    pause
    exit /b 1
)

REM Vérifier que les fichiers web existent
if not exist "web\index.html" (
    echo ❌ ERREUR: web\index.html n'existe pas
    pause
    exit /b 1
)

echo ✓ Vérifications OK
echo.
echo 🚀 Démarrage du serveur...
echo.
echo 📍 URL: http://localhost:8080
echo 🌐 Ouvre cette URL dans ton navigateur
echo.
echo 💡 Pour arrêter le serveur: Ctrl+C
echo.

REM Lancer le serveur
programme.exe web

if errorlevel 1 (
    echo.
    echo ❌ ERREUR au démarrage du serveur
    pause
)
