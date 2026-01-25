#  Windows Setup - IA-ATOMIQUE

Installation ultra-simple pour Windows.

##  3 àtapes Seulement

### 1à Installer Go (une fois pour toutes)

- **Télécharge:** https://golang.org/dl/
- **Choisis:** `Windows x86-64` (la version 64-bit pour la plupart)
- **Installe:** Clique Next partout, c'est facile
- **Redémarre** (optionnel mais recommandé)

**Vérification:**
- Ouvre Command Prompt
- Tape: `go version`
- Tu devrais voir: `go version go1.22+ windows/amd64`

---

### 2à Clone le Repo

```bash
git clone https://github.com/TON-USER/ia-atomique.git
cd ia-atomique
```

---

### 3à Lance l'Installation

**Option A: Batch Files (Command Prompt)**
```bash
install.bat
run-web.bat
```

**Option B: PowerShell Scripts**
```powershell
.\install.ps1
.\run-web.ps1
```

**Option C: Double-clique Simple**
- Double-clique `install.bat`
- Double-clique `run-web.bat`
- C'est tout!

---

##  Fichiers Disponibles

### Batch Files (.bat) - Command Prompt
| Fichier | Action |
|---------|--------|
| `install.bat` | Installation initiale |
| `run-web.bat` | Lancer le serveur |
| `build.bat` | Recompiler |
| `verify.bat` | Vérifier tout |

### PowerShell Scripts (.ps1)
| Fichier | Action |
|---------|--------|
| `install.ps1` | Installation initiale |
| `run-web.ps1` | Lancer le serveur |
| `build.ps1` | Recompiler |

---

##  Utilisation

### Premiàre Fois
```bash
# 1. Double-clique install.bat (ou .\install.ps1)
# Ça compile tout automatiquement

# 2. Double-clique run-web.bat (ou .\run-web.ps1)
# Le serveur démarre

# 3. Ouvre navigateur: http://localhost:8080
```

### Prochaines Fois
```bash
# Juste: Double-clique run-web.bat
# Et voilà!
```

---

##  Troubleshooting

###  "Go is not recognized"
- Go n'est pas installé
- Télécharge depuis https://golang.org/dl/
- Redémarre Command Prompt apràs installation

###  "Programme.exe n'existe pas"
- Lancer `install.bat` ou `build.bat` d'abord

###  Port 8080 occupé
- Ferme l'autre application
- Ou modifie `web.go` ligne 29 pour utiliser port 9000

###  PowerShell dit "Script cannot be run"
- Ouvre PowerShell en tant qu'admin
- Tape: `Set-ExecutionPolicy -ExecutionPolicy RemoteSigned`
- Tape: `Y` et Enter
- Réessaye le script

---

##  Astuces Windows

### Créer un Raccourci Bureau
1. Clique droit sur `run-web.bat`
2. "Envoyer vers > Bureau (créer un raccourci)"
3. Double-clique le raccourci pour lancer!

### Lancer au Démarrage Windows
1. Win+R  `shell:startup`  Enter
2. Copie un raccourci de `run-web.bat` là
3. à chaque redémarrage, le serveur se lance!

### Créer un Fichier VBS (Cache la fenàtre)
Crée `run-web-silent.vbs`:
```vbs
Set objShell = CreateObject("WScript.Shell")
objShell.Run "cmd /c run-web.bat", 0
WScript.Quit
```

---

##  Voir Aussi

- [WINDOWS_INSTALL.md](WINDOWS_INSTALL.md) - Guide détaillé Windows
- [INSTALL.md](INSTALL.md) - Installation générale
- [WEB_README.md](WEB_README.md) - Configuration web

---

**C'est tout! Rien de compliqué! **
