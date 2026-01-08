@echo off
REM Vérification rapide avant push GitHub

cls

echo.
echo ╔════════════════════════════════════════╗
echo ║  IA-ATOMIQUE - Verification            ║
echo ╚════════════════════════════════════════╝
echo.

REM Vérifier la compilation
echo 📦 Vérification de la compilation...
go build -o programme.exe
if errorlevel 1 (
    echo ❌ ERREUR: Compilation échouée
    pause
    exit /b 1
)
echo ✓ Compilation OK
echo.

REM Vérifier les fichiers web
echo 🌐 Vérification des fichiers web...
if not exist "web\index.html" (
    echo ❌ web\index.html manquant
    pause
    exit /b 1
)
echo ✓ web\index.html existe
if not exist "web\style.css" (
    echo ❌ web\style.css manquant
    pause
    exit /b 1
)
echo ✓ web\style.css existe
if not exist "web\script.js" (
    echo ❌ web\script.js manquant
    pause
    exit /b 1
)
echo ✓ web\script.js existe
echo.

REM Vérifier les fichiers config
echo ⚙️  Vérification de la configuration...
if not exist "Dockerfile" (
    echo ❌ Dockerfile manquant
    pause
    exit /b 1
)
echo ✓ Dockerfile existe
if not exist ".gitignore" (
    echo ❌ .gitignore manquant
    pause
    exit /b 1
)
echo ✓ .gitignore existe
echo.

echo.
echo ╔════════════════════════════════════════╗
echo ║  ✅ Vérification OK!                   ║
echo ╚════════════════════════════════════════╝
echo.
pause
