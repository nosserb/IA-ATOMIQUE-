#  R√©sum√© Complet - Installation Windows

## ¶ Qu'est-ce qui a √©t√© cr√©√© pour Windows?

### Scripts d'Installation (Batch)
```
install.bat       Installation en 1 clic
run-web.bat       Lancer le serveur en 1 clic
build.bat         Recompiler en 1 clic
verify.bat        V√©rifier tout en 1 clic
```

### Scripts PowerShell (Avanc√©)
```
install.ps1       Installation (PowerShell)
run-web.ps1       Lancer le serveur (PowerShell)
build.ps1         Recompiler (PowerShell)
```

### Documentation Windows
```
WINDOWS_INSTALL.md   Guide d√©taill√© (tu es ici!)
WINDOWS_SETUP.md     Quick start
README.md            Mis √ jour avec section Windows
```

---

##  Installation Super Simple (3 √©tapes)

### √tape 1: T√©l√©charger Go
- Va sur: https://golang.org/dl/
- T√©l√©charge: `Windows x86-64.msi`
- Installe: Clique "Next" partout
- Red√©marre: Windows (recommand√©)

### √tape 2: Clone le Repo
```
git clone https://github.com/TON-USER/ia-atomique.git
cd ia-atomique
```

### √tape 3: Lancer
**Option A - 1 Clic (Recommand√©)**
```
Double-clique: install.bat
Puis double-clique: run-web.bat
Fini! 
```

**Option B - Command Prompt**
```
cmd
C:\> install.bat
C:\> run-web.bat
```

**Option C - PowerShell**
```
powershell
PS> .\install.ps1
PS> .\run-web.ps1
```

---

##  Qu'est-ce que chaque script fait?

### install.bat / install.ps1
-  V√©rifie que Go est install√©
-  Compile le projet
-  V√©rifie les fichiers web
-  Message final "Pr√t!"

### run-web.bat / run-web.ps1
-  V√©rifie que le binaire existe
-  V√©rifie les fichiers web
-  D√©marre le serveur sur port 8080
-  Affiche l'URL √ ouvrir

### build.bat / build.ps1
-  Compile le projet
-  Cr√©e `programme.exe`

### verify.bat
-  V√©rifie la compilation
-  V√©rifie les fichiers web
-  V√©rifie la configuration

---

##  Astuces Pratiques

### Cr√©er un Raccourci Bureau
1. Clique droit sur `run-web.bat`
2. "Envoyer vers > Bureau (cr√©er un raccourci)"
3. Double-clique le raccourci: Serveur lance!

### Lancer Automatiquement au D√©marrage
1. Win+R  tape: `shell:startup`  OK
2. Copie un raccourci de `run-web.bat` dans le dossier
3. √ chaque d√©marrage: Serveur lanc√©!

### Lancer sans Voir la Fen√tre CMD
Cr√©e fichier `run-web-hidden.vbs`:
```vbs
Set objShell = CreateObject("WScript.Shell")
objShell.Run "cmd /c run-web.bat", 0
WScript.Quit
```
Puis double-clique ce fichier!

---

##  Personnaliser le Port

Si le port 8080 est occup√©:

1. Ouvre `web.go` avec un √©diteur (Notepad++, VSCode, etc.)
2. Cherche la ligne: `go StartWebServer("8080")`
3. Change 8080 en autre port: `go StartWebServer("9000")`
4. Lance: `build.bat`
5. Lance: `run-web.bat`
6. Acc√de √: `http://localhost:9000`

---

##  D√©pannage

### Probl√me: "Go is not recognized"
```
Solution:
1. T√©l√©charge Go: https://golang.org/dl/
2. Installe le .msi
3. Red√©marre Windows compl√tement
4. R√©essaye
```

### Probl√me: Script PowerShell ne marche pas
```
Solution:
1. Clique droit PowerShell
2. "Run as Administrator"
3. Tape: Set-ExecutionPolicy -ExecutionPolicy RemoteSigned
4. Tape: Y et Enter
5. R√©essaye le script
```

### Probl√me: Port 8080 d√©j√ utilis√©
```
Solution A: Ferme l'autre application
Solution B: Change le port dans web.go (voir section Personnaliser)
```

### Probl√me: Fichiers web manquants
```
Solution:
1. Assure-toi que dossier "web" existe
2. Avec fichiers: index.html, style.css, script.js
3. R√©essaye
```

---

##  Structure pour Windows

```
ia-atomique/
 install.bat           √ double-cliquer en 1er
 run-web.bat           √ double-cliquer ensuite
 build.bat
 verify.bat

 install.ps1           Alternative PowerShell
 run-web.ps1
 build.ps1

 WINDOWS_INSTALL.md    Tu es ici!
 WINDOWS_SETUP.md

 programme.exe         Cr√©√© apr√s install.bat
 web/
    index.html
    style.css
    script.js

 ... autres fichiers
```

---

##  Workflow Typique Windows

```
1£ Double-clique install.bat
    Affiche: " Installation R√©ussie!"

2£ Double-clique run-web.bat
    Lance le serveur

3£ Une fen√tre CMD s'ouvre avec:
   - "Serveur Web D√©marr√©"
   - "URL: http://localhost:8080"

4£ Ouvre navigateur: http://localhost:8080

5£ Utilise l'interface web! 

6£ Pour arr√ter: Ferme la fen√tre CMD
```

---

##  Points Importants

 **Aucune ligne de commande requise** - Juste double-clique  
 **Automatique** - Install.bat fait tout tout seul  
 **Erreurs claires** - Sait dire exactement ce qui manque  
 **PowerShell aussi** - Pour ceux qui pr√©f√rent  
 **Portable** - Pas de d√©pendances externes  

---

##  Prochaines √tapes

### Pour toi (d√©veloppeur):
1. Teste les scripts sur Windows
2. Modifie le code si besoin
3. Lance `build.bat` pour recompiler
4. Lance `run-web.bat` pour tester
5. Pousse sur GitHub!

### Pour les utilisateurs:
1. Clone le repo
2. Double-clique `install.bat`
3. Double-clique `run-web.bat`
4. Ouvre navigateur
5. C'est tout! 

---

## û Support Windows

Si quelque chose marche pas:
1. Relance Windows compl√tement
2. R√©installe Go proprement
3. Double-clique `verify.bat` pour v√©rifier
4. Essaye √ nouveau

---

**Voil√! Installation super simple pour Windows! **
