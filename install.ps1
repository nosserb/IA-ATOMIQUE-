#!/usr/bin/env pwsh

# Script PowerShell d'installation pour Windows
# IA-ATOMIQUE Windows Installer - PowerShell

Clear-Host

Write-Host ""
Write-Host "╔════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║  IA-ATOMIQUE - Installation Windows    ║" -ForegroundColor Cyan
Write-Host "╚════════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host ""

# Étape 1: Vérifier Go
Write-Host "📋 Étape 1/3: Vérification de Go..." -ForegroundColor Yellow
$goCheck = try { go version } catch { $false }
if (-not $goCheck) {
    Write-Host "❌ Go n'est pas installé!" -ForegroundColor Red
    Write-Host ""
    Write-Host "Télécharge Go depuis:" -ForegroundColor Yellow
    Write-Host "https://golang.org/dl/" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "1. Télécharge le fichier .msi pour Windows" -ForegroundColor White
    Write-Host "2. Exécute l'installateur" -ForegroundColor White
    Write-Host "3. Redémarre PowerShell" -ForegroundColor White
    Write-Host "4. Relance ce script" -ForegroundColor White
    Write-Host ""
    Read-Host "Appuie sur Entrée pour fermer"
    exit 1
}
Write-Host "$(go version)" -ForegroundColor Green
Write-Host "✓ Go OK" -ForegroundColor Green
Write-Host ""

# Étape 2: Compiler
Write-Host "📋 Étape 2/3: Compilation..." -ForegroundColor Yellow
Write-Host "📦 Compilation du projet..." -ForegroundColor Cyan
go build -o programme.exe
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ ERREUR: La compilation a échoué" -ForegroundColor Red
    Read-Host "Appuie sur Entrée pour fermer"
    exit 1
}
Write-Host "✓ Compilation OK" -ForegroundColor Green
Write-Host ""

# Étape 3: Vérification des fichiers web
Write-Host "📋 Étape 3/3: Vérification de l'interface web..." -ForegroundColor Yellow
if (-not (Test-Path "web\index.html")) {
    Write-Host "❌ ERREUR: Fichiers web manquants" -ForegroundColor Red
    Read-Host "Appuie sur Entrée pour fermer"
    exit 1
}
Write-Host "✓ Fichiers OK" -ForegroundColor Green
Write-Host ""

Write-Host ""
Write-Host "╔════════════════════════════════════════╗" -ForegroundColor Green
Write-Host "║  ✅ Installation Réussie!              ║" -ForegroundColor Green
Write-Host "╚════════════════════════════════════════╝" -ForegroundColor Green
Write-Host ""
Write-Host "🚀 Prêt à démarrer!" -ForegroundColor Green
Write-Host ""
Write-Host "Commandes:" -ForegroundColor Cyan
Write-Host "  .\run-web.ps1     → Lancer l'interface web" -ForegroundColor White
Write-Host "  .\build.ps1       → Recompiler" -ForegroundColor White
Write-Host ""
Read-Host "Appuie sur Entrée pour fermer"
