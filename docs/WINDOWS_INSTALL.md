#  IA-ATOMIQUE sur Windows

Guide d'installation simple pour Windows.

##  Installation Rapide (2 etapes)

### �tape 1: Installer Go (une seule fois)

1. **Telecharge Go** depuis https://golang.org/dl/
2. **Choisir la version Windows** (msi):
   - Windows x86-64 (64-bit) = La plupart des gens
   - Windows x86 (32-bit) = PC ancien
3. **Execute l'installateur** et clique "Next" partout
4. **Redemarre** (optionnel mais recommande)

**Verifier que Go est installe:**
- Ouvre PowerShell ou Command Prompt
- Tape: `go version`
- Tu dois voir quelque chose comme: `go version go1.22.2 windows/amd64`

Si �a marche pas, redemarre ton PC et reessaye.

---

### �tape 2: Lancer le projet

1. **Clone le repo GitHub**
   ```
   git clone https://github.com/TON-USER/ia-atomique.git
   cd ia-atomique
   ```

2. **Double-clique sur `install.bat`**
   - Ça compile automatiquement
   - Ça verifie les fichiers
   - C'est tout!

3. **Double-clique sur `run-web.bat`**
   - Le serveur demarre
   - Une fen�tre s'ouvre

4. **Ouvre ton navigateur**
   - Acc�de �: http://localhost:8080
   - Et voil�! 

---

##  Fichiers Disponibles

### Scripts Windows

| Fichier | Fonction |
|---------|----------|
| `install.bat` | Installation initiale (� faire une seule fois) |
| `run-web.bat` | Lancer le serveur web |
| `build.bat` | Recompiler le projet |
| `verify.bat` | Verifier que tout est OK |

### Comment les utiliser

**Double-clique simplement sur le fichier!**

Ou ouvre Command Prompt/PowerShell et tape:
```bash
install.bat
run-web.bat
build.bat
verify.bat
```

---

##  Mode Command Prompt/PowerShell

Si tu pref�res la ligne de commande:

### Premi�re utilisation
```bash
# Ouvre Command Prompt (Cmd.exe)
# Navigue vers le repertoire:
cd "C:\Users\TON-USER\Downloads\ia-atomique"

# Compile:
go build -o programme.exe

# Lancer:
programme.exe web
```

### Prochaines utilisations
```bash
# Ouvrir Command Prompt dans le dossier
# Puis:
programme.exe web
```

---

##  Resolution des Probl�mes

###  "Go is not recognized"

**Probl�me:** Go n'est pas installe ou trouvable

**Solution:**
1. Installe Go depuis https://golang.org/dl/
2. Redemarre ton PC compl�tement
3. Reessaye

###  "web\index.html n'existe pas"

**Probl�me:** Fichiers web manquants

**Solution:**
- Assure-toi que le dossier `web/` existe
- Avec les fichiers: `index.html`, `style.css`, `script.js`

###  Port 8080 dej� utilise

**Probl�me:** Un autre programme utilise le port 8080

**Solution 1:**
- Ferme l'autre application

**Solution 2:**
- Modifie `web.go` ligne 29:
  ```go
  go StartWebServer("9000")  // Change 8080 en 9000
  ```
- Recompile: `build.bat`
- Lance avec: `programme.exe web`
- Acc�de �: http://localhost:9000

###  "programme.exe n'existe pas"

**Probl�me:** Pas compile

**Solution:**
- Double-clique `build.bat` ou `install.bat`

---

##  Astuces Windows

### Creer un raccourci sur le Bureau

1. **Clique droit sur `run-web.bat`**
2. **Envoyer vers > Bureau (creer un raccourci)**
3. **Double-clique le raccourci pour lancer!**

### Lancer dans une vraie fen�tre (pas CMD)

Pour une experience plus fluide:

1. **Cree `run-web.vbs`:**

```vbs
Set objShell = CreateObject("WScript.Shell")
objShell.Run "cmd /c run-web.bat", 0
WScript.Quit
```

2. **Sauvegarde dans le m�me dossier**
3. **Double-clique pour lancer sans fen�tre CMD**

### Ajouter au menu Demarrer

1. **Appuie sur Win+R**
2. **Tape: `shell:startup`**
3. **Mets un raccourci de `run-web.bat` l�**
4. **� chaque demarrage, le serveur se lance!**

---

##  Workflow Typique

```
1. Double-clique install.bat (premi�re fois seulement)
   
2. Double-clique run-web.bat
   
3. Une fen�tre s'ouvre avec le serveur
   
4. Ouvre navigateur: http://localhost:8080
   
5. Utilise l'interface!
   
6. Pour arr�ter: Ferme la fen�tre CMD
```

---

##  Pour Developpeurs

Si tu modifies le code:

```bash
# Compile
build.bat

# Teste
run-web.bat
```

---

##  Avant de Pusher sur GitHub

```bash
# Verifier que tout compile
verify.bat

# Si tout est OK, push!
git add .
git commit -m "feat: Add web interface"
git push origin main
```

---

## � Support

Si quelque chose ne marche pas:

1. **Verifie que Go est installe:** `go version`
2. **Recompile:** `build.bat`
3. **Verification:** `verify.bat`
4. **Reessaye:** `run-web.bat`

Si �a marche toujours pas:
- Redemarre ton PC
- Reinstalle Go proprement
- Essaye � nouveau

---

**C'est simple! Bon resume! **
