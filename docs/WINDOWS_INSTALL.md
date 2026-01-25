#  IA-ATOMIQUE sur Windows

Guide d'installation simple pour Windows.

##  Installation Rapide (2 √©tapes)

### √tape 1: Installer Go (une seule fois)

1. **T√©l√©charge Go** depuis https://golang.org/dl/
2. **Choisir la version Windows** (msi):
   - Windows x86-64 (64-bit) = La plupart des gens
   - Windows x86 (32-bit) = PC ancien
3. **Ex√©cute l'installateur** et clique "Next" partout
4. **Red√©marre** (optionnel mais recommand√©)

**V√©rifier que Go est install√©:**
- Ouvre PowerShell ou Command Prompt
- Tape: `go version`
- Tu dois voir quelque chose comme: `go version go1.22.2 windows/amd64`

Si √a marche pas, red√©marre ton PC et r√©essaye.

---

### √tape 2: Lancer le projet

1. **Clone le repo GitHub**
   ```
   git clone https://github.com/TON-USER/ia-atomique.git
   cd ia-atomique
   ```

2. **Double-clique sur `install.bat`**
   - √áa compile automatiquement
   - √áa v√©rifie les fichiers
   - C'est tout!

3. **Double-clique sur `run-web.bat`**
   - Le serveur d√©marre
   - Une fen√tre s'ouvre

4. **Ouvre ton navigateur**
   - Acc√de √: http://localhost:8080
   - Et voil√! 

---

##  Fichiers Disponibles

### Scripts Windows

| Fichier | Fonction |
|---------|----------|
| `install.bat` | Installation initiale (√ faire une seule fois) |
| `run-web.bat` | Lancer le serveur web |
| `build.bat` | Recompiler le projet |
| `verify.bat` | V√©rifier que tout est OK |

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

Si tu pr√©f√res la ligne de commande:

### Premi√re utilisation
```bash
# Ouvre Command Prompt (Cmd.exe)
# Navigue vers le r√©pertoire:
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

##  R√©solution des Probl√mes

###  "Go is not recognized"

**Probl√me:** Go n'est pas install√© ou trouvable

**Solution:**
1. Installe Go depuis https://golang.org/dl/
2. Red√©marre ton PC compl√tement
3. R√©essaye

###  "web\index.html n'existe pas"

**Probl√me:** Fichiers web manquants

**Solution:**
- Assure-toi que le dossier `web/` existe
- Avec les fichiers: `index.html`, `style.css`, `script.js`

###  Port 8080 d√©j√ utilis√©

**Probl√me:** Un autre programme utilise le port 8080

**Solution 1:**
- Ferme l'autre application

**Solution 2:**
- Modifie `web.go` ligne 29:
  ```go
  go StartWebServer("9000")  // Change 8080 en 9000
  ```
- Recompile: `build.bat`
- Lance avec: `programme.exe web`
- Acc√de √: http://localhost:9000

###  "programme.exe n'existe pas"

**Probl√me:** Pas compil√©

**Solution:**
- Double-clique `build.bat` ou `install.bat`

---

##  Astuces Windows

### Cr√©er un raccourci sur le Bureau

1. **Clique droit sur `run-web.bat`**
2. **Envoyer vers > Bureau (cr√©er un raccourci)**
3. **Double-clique le raccourci pour lancer!**

### Lancer dans une vraie fen√tre (pas CMD)

Pour une exp√©rience plus fluide:

1. **Cr√©e `run-web.vbs`:**

```vbs
Set objShell = CreateObject("WScript.Shell")
objShell.Run "cmd /c run-web.bat", 0
WScript.Quit
```

2. **Sauvegarde dans le m√me dossier**
3. **Double-clique pour lancer sans fen√tre CMD**

### Ajouter au menu D√©marrer

1. **Appuie sur Win+R**
2. **Tape: `shell:startup`**
3. **Mets un raccourci de `run-web.bat` l√**
4. **√ chaque d√©marrage, le serveur se lance!**

---

##  Workflow Typique

```
1. Double-clique install.bat (premi√re fois seulement)
   
2. Double-clique run-web.bat
   
3. Une fen√tre s'ouvre avec le serveur
   
4. Ouvre navigateur: http://localhost:8080
   
5. Utilise l'interface!
   
6. Pour arr√ter: Ferme la fen√tre CMD
```

---

##  Pour D√©veloppeurs

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
# V√©rifier que tout compile
verify.bat

# Si tout est OK, push!
git add .
git commit -m "feat: Add web interface"
git push origin main
```

---

## û Support

Si quelque chose ne marche pas:

1. **V√©rifie que Go est install√©:** `go version`
2. **Recompile:** `build.bat`
3. **V√©rification:** `verify.bat`
4. **R√©essaye:** `run-web.bat`

Si √a marche toujours pas:
- Red√©marre ton PC
- R√©installe Go proprement
- Essaye √ nouveau

---

**C'est simple! Bon r√©sum√©! **
