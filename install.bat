@echo off
REM Script d'installation pour Windows
REM IA-ATOMIQUE Windows Installer

setlocal enabledelayexpansion

cls

echo.
echo ╔════════════════════════════════════════╗
echo ║  IA-ATOMIQUE - Installation Windows    ║
echo ╚════════════════════════════════════════╝
echo.

REM Étape 1: Vérifier Go
echo 📋 Étape 1/3: Vérification de Go...
go version >nul 2>&1
if errorlevel 1 (
    echo ❌ Go n'est pas installé!
    echo.
    echo Télécharge Go depuis:
    echo https://golang.org/dl/
    echo.
    echo 1. Télécharge le fichier .msi pour Windows
    echo 2. Exécute l'installateur
    echo 3. Relance ce script
    echo.
    pause
    exit /b 1
)
go version
echo ✓ Go OK
echo.

REM Étape 2: Compiler
echo 📋 Étape 2/3: Compilation...
echo 📦 Compilation du projet...
go build -o programme.exe
if errorlevel 1 (
    echo ❌ ERREUR: La compilation a échoué
    pause
    exit /b 1
)
echo ✓ Compilation OK
echo.

REM Étape 3: Vérification des fichiers web
echo 📋 Étape 3/3: Vérification de l'interface web...
if not exist "web\index.html" (
    echo ❌ ERREUR: Fichiers web manquants
    pause
    exit /b 1
)
echo ✓ Fichiers OK
echo.

echo.
echo ╔════════════════════════════════════════╗
echo ║  ✅ Installation Réussie!              ║
echo ╚════════════════════════════════════════╝
echo.
echo 🚀 Prêt à démarrer!
echo.
echo Commandes:
echo   run-web.bat     → Lancer l'interface web
echo   build.bat       → Recompiler
echo.
pause
