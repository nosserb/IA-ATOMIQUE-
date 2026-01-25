#  Resume Complet - Installation Windows

## � Qu'est-ce qui a ete cree pour Windows?

### Scripts d'Installation (Batch)
```
install.bat       Installation en 1 clic
run-web.bat       Lancer le serveur en 1 clic
build.bat         Recompiler en 1 clic
verify.bat        Verifier tout en 1 clic
```

### Scripts PowerShell (Avance)
```
install.ps1       Installation (PowerShell)
run-web.ps1       Lancer le serveur (PowerShell)
build.ps1         Recompiler (PowerShell)
```

### Documentation Windows
```
WINDOWS_INSTALL.md   Guide detaille (tu es ici!)
WINDOWS_SETUP.md     Quick start
README.md            Mis � jour avec section Windows
```

---

##  Installation Super Simple (3 etapes)

### �tape 1: Telecharger Go
- Va sur: https://golang.org/dl/
- Telecharge: `Windows x86-64.msi`
- Installe: Clique "Next" partout
- Redemarre: Windows (recommande)

### �tape 2: Clone le Repo
```
git clone https://github.com/TON-USER/ia-atomique.git
cd ia-atomique
```

### �tape 3: Lancer
**Option A - 1 Clic (Recommande)**
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
-  Verifie que Go est installe
-  Compile le projet
-  Verifie les fichiers web
-  Message final "Pr�t!"

### run-web.bat / run-web.ps1
-  Verifie que le binaire existe
-  Verifie les fichiers web
-  Demarre le serveur sur port 8080
-  Affiche l'URL � ouvrir

### build.bat / build.ps1
-  Compile le projet
-  Cree `programme.exe`

### verify.bat
-  Verifie la compilation
-  Verifie les fichiers web
-  Verifie la configuration

---

##  Astuces Pratiques

### Creer un Raccourci Bureau
1. Clique droit sur `run-web.bat`
2. "Envoyer vers > Bureau (creer un raccourci)"
3. Double-clique le raccourci: Serveur lance!

### Lancer Automatiquement au Demarrage
1. Win+R  tape: `shell:startup`  OK
2. Copie un raccourci de `run-web.bat` dans le dossier
3. � chaque demarrage: Serveur lance!

### Lancer sans Voir la Fen�tre CMD
Cree fichier `run-web-hidden.vbs`:
```vbs
Set objShell = CreateObject("WScript.Shell")
objShell.Run "cmd /c run-web.bat", 0
WScript.Quit
```
Puis double-clique ce fichier!

---

##  Personnaliser le Port

Si le port 8080 est occupe:

1. Ouvre `web.go` avec un editeur (Notepad++, VSCode, etc.)
2. Cherche la ligne: `go StartWebServer("8080")`
3. Change 8080 en autre port: `go StartWebServer("9000")`
4. Lance: `build.bat`
5. Lance: `run-web.bat`
6. Acc�de �: `http://localhost:9000`

---

##  Depannage

### Probl�me: "Go is not recognized"
```
Solution:
1. Telecharge Go: https://golang.org/dl/
2. Installe le .msi
3. Redemarre Windows compl�tement
4. Reessaye
```

### Probl�me: Script PowerShell ne marche pas
```
Solution:
1. Clique droit PowerShell
2. "Run as Administrator"
3. Tape: Set-ExecutionPolicy -ExecutionPolicy RemoteSigned
4. Tape: Y et Enter
5. Reessaye le script
```

### Probl�me: Port 8080 dej� utilise
```
Solution A: Ferme l'autre application
Solution B: Change le port dans web.go (voir section Personnaliser)
```

### Probl�me: Fichiers web manquants
```
Solution:
1. Assure-toi que dossier "web" existe
2. Avec fichiers: index.html, style.css, script.js
3. Reessaye
```

---

##  Structure pour Windows

```
ia-atomique/
 install.bat           � double-cliquer en 1er
 run-web.bat           � double-cliquer ensuite
 build.bat
 verify.bat

 install.ps1           Alternative PowerShell
 run-web.ps1
 build.ps1

 WINDOWS_INSTALL.md    Tu es ici!
 WINDOWS_SETUP.md

 programme.exe         Cree apr�s install.bat
 web/
    index.html
    style.css
    script.js

 ... autres fichiers
```

---

##  Workflow Typique Windows

```
1� Double-clique install.bat
    Affiche: " Installation Reussie!"

2� Double-clique run-web.bat
    Lance le serveur

3� Une fen�tre CMD s'ouvre avec:
   - "Serveur Web Demarre"
   - "URL: http://localhost:8080"

4� Ouvre navigateur: http://localhost:8080

5� Utilise l'interface web! 

6� Pour arr�ter: Ferme la fen�tre CMD
```

---

##  Points Importants

 **Aucune ligne de commande requise** - Juste double-clique  
 **Automatique** - Install.bat fait tout tout seul  
 **Erreurs claires** - Sait dire exactement ce qui manque  
 **PowerShell aussi** - Pour ceux qui pref�rent  
 **Portable** - Pas de dependances externes  

---

##  Prochaines �tapes

### Pour toi (developpeur):
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

## � Support Windows

Si quelque chose marche pas:
1. Relance Windows compl�tement
2. Reinstalle Go proprement
3. Double-clique `verify.bat` pour verifier
4. Essaye � nouveau

---

**Voil�! Installation super simple pour Windows! **
