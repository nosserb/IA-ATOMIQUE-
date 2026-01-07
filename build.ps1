#!/usr/bin/env pwsh

# Script PowerShell pour compiler le projet
# IA-ATOMIQUE Build Script - Windows PowerShell

Write-Host ""
Write-Host "╔════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║  IA-ATOMIQUE - Build Windows           ║" -ForegroundColor Cyan
Write-Host "╚════════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host ""

# Vérifier si Go est installé
$goCheck = try { go version } catch { $false }
if (-not $goCheck) {
    Write-Host "❌ ERREUR: Go n'est pas installé ou pas dans le PATH" -ForegroundColor Red
    Write-Host ""
    Write-Host "Télécharge Go depuis: https://golang.org/dl/" -ForegroundColor Yellow
    Write-Host ""
    Read-Host "Appuie sur Entrée pour fermer"
    exit 1
}

Write-Host "✓ Go détecté" -ForegroundColor Green
Write-Host ""

# Compiler le projet
Write-Host "📦 Compilation du projet..." -ForegroundColor Cyan
go build -o programme.exe
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ ERREUR: La compilation a échoué" -ForegroundColor Red
    Read-Host "Appuie sur Entrée pour fermer"
    exit 1
}

Write-Host "✓ Compilation réussie!" -ForegroundColor Green
Write-Host ""
Write-Host "🎉 programme.exe créé avec succès" -ForegroundColor Green
Write-Host ""
Write-Host "Commandes disponibles:" -ForegroundColor Cyan
Write-Host "  run-web.ps1    - Lancer l'interface web" -ForegroundColor White
Write-Host "  .\programme.exe web  - Lancer directement" -ForegroundColor White
Write-Host ""
Read-Host "Appuie sur Entrée pour fermer"
