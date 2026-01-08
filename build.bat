@echo off
REM Script pour compiler le projet sur Windows
REM IA-ATOMIQUE Build Script

setlocal enabledelayexpansion

echo.
echo ╔════════════════════════════════════════╗
echo ║  IA-ATOMIQUE - Build Windows           ║
echo ╚════════════════════════════════════════╝
echo.

REM Vérifier si Go est installé
go version >nul 2>&1
if errorlevel 1 (
    echo ❌ ERREUR: Go n'est pas installé ou pas dans le PATH
    echo.
    echo Télécharge Go depuis: https://golang.org/dl/
    echo.
    pause
    exit /b 1
)

echo ✓ Go détecté
echo.

REM Compiler le projet
echo 📦 Compilation du projet...
go build -o programme.exe
if errorlevel 1 (
    echo ❌ ERREUR: La compilation a échoué
    pause
    exit /b 1
)

echo ✓ Compilation réussie!
echo.
echo 🎉 programme.exe créé avec succès
echo.
echo Commandes disponibles:
echo   run-web.bat    - Lancer l'interface web
echo   programme.exe web  - Lancer directement
echo.
pause
