#!/usr/bin/env pwsh

# Script PowerShell pour lancer le serveur web
# IA-ATOMIQUE Web Server - Windows PowerShell

Write-Host ""
Write-Host "╔════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║  IA-ATOMIQUE - Serveur Web             ║" -ForegroundColor Cyan
Write-Host "╚════════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host ""

# Vérifier si le binaire existe
if (-not (Test-Path "programme.exe")) {
    Write-Host "❌ ERREUR: programme.exe n'existe pas" -ForegroundColor Red
    Write-Host ""
    Write-Host "Exécute d'abord: build.bat" -ForegroundColor Yellow
    Write-Host ""
    Read-Host "Appuie sur Entrée pour fermer"
    exit 1
}

# Vérifier que les fichiers web existent
if (-not (Test-Path "web\index.html")) {
    Write-Host "❌ ERREUR: web\index.html n'existe pas" -ForegroundColor Red
    Read-Host "Appuie sur Entrée pour fermer"
    exit 1
}

Write-Host "✓ Vérifications OK" -ForegroundColor Green
Write-Host ""
Write-Host "🚀 Démarrage du serveur..." -ForegroundColor Green
Write-Host ""
Write-Host "📍 URL: http://localhost:8080" -ForegroundColor Cyan
Write-Host "🌐 Ouvre cette URL dans ton navigateur" -ForegroundColor Cyan
Write-Host ""
Write-Host "💡 Pour arrêter le serveur: Ctrl+C" -ForegroundColor Yellow
Write-Host ""

# Lancer le serveur
& .\programme.exe web

if ($LASTEXITCODE -ne 0) {
    Write-Host ""
    Write-Host "❌ ERREUR au démarrage du serveur" -ForegroundColor Red
    Read-Host "Appuie sur Entrée pour fermer"
}
