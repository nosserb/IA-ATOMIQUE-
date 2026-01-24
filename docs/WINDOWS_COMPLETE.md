# 🎯 Résumé Complet - Installation Windows

## 📦 Qu'est-ce qui a été créé pour Windows?

### Scripts d'Installation (Batch)
```
install.bat      ← Installation en 1 clic
run-web.bat      ← Lancer le serveur en 1 clic
build.bat        ← Recompiler en 1 clic
verify.bat       ← Vérifier tout en 1 clic
```

### Scripts PowerShell (Avancé)
```
install.ps1      ← Installation (PowerShell)
run-web.ps1      ← Lancer le serveur (PowerShell)
build.ps1        ← Recompiler (PowerShell)
```

### Documentation Windows
```
WINDOWS_INSTALL.md  ← Guide détaillé (tu es ici!)
WINDOWS_SETUP.md    ← Quick start
README.md           ← Mis à jour avec section Windows
```

---

## 🚀 Installation Super Simple (3 étapes)

### Étape 1: Télécharger Go
- Va sur: https://golang.org/dl/
- Télécharge: `Windows x86-64.msi`
- Installe: Clique "Next" partout
- Redémarre: Windows (recommandé)

### Étape 2: Clone le Repo
```
git clone https://github.com/TON-USER/ia-atomique.git
cd ia-atomique
```

### Étape 3: Lancer
**Option A - 1 Clic (Recommandé)**
```
Double-clique: install.bat
Puis double-clique: run-web.bat
Fini! 🎉
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

## 📝 Qu'est-ce que chaque script fait?

### install.bat / install.ps1
- ✅ Vérifie que Go est installé
- ✅ Compile le projet
- ✅ Vérifie les fichiers web
- ✅ Message final "Prêt!"

### run-web.bat / run-web.ps1
- ✅ Vérifie que le binaire existe
- ✅ Vérifie les fichiers web
- ✅ Démarre le serveur sur port 8080
- ✅ Affiche l'URL à ouvrir

### build.bat / build.ps1
- ✅ Compile le projet
- ✅ Crée `programme.exe`

### verify.bat
- ✅ Vérifie la compilation
- ✅ Vérifie les fichiers web
- ✅ Vérifie la configuration

---

## 💡 Astuces Pratiques

### Créer un Raccourci Bureau
1. Clique droit sur `run-web.bat`
2. "Envoyer vers > Bureau (créer un raccourci)"
3. Double-clique le raccourci: Serveur lance!

### Lancer Automatiquement au Démarrage
1. Win+R → tape: `shell:startup` → OK
2. Copie un raccourci de `run-web.bat` dans le dossier
3. À chaque démarrage: Serveur lancé!

### Lancer sans Voir la Fenêtre CMD
Crée fichier `run-web-hidden.vbs`:
```vbs
Set objShell = CreateObject("WScript.Shell")
objShell.Run "cmd /c run-web.bat", 0
WScript.Quit
```
Puis double-clique ce fichier!

---

## 🔧 Personnaliser le Port

Si le port 8080 est occupé:

1. Ouvre `web.go` avec un éditeur (Notepad++, VSCode, etc.)
2. Cherche la ligne: `go StartWebServer("8080")`
3. Change 8080 en autre port: `go StartWebServer("9000")`
4. Lance: `build.bat`
5. Lance: `run-web.bat`
6. Accède à: `http://localhost:9000`

---

## ❌ Dépannage

### Problème: "Go is not recognized"
```
Solution:
1. Télécharge Go: https://golang.org/dl/
2. Installe le .msi
3. Redémarre Windows complètement
4. Réessaye
```

### Problème: Script PowerShell ne marche pas
```
Solution:
1. Clique droit PowerShell
2. "Run as Administrator"
3. Tape: Set-ExecutionPolicy -ExecutionPolicy RemoteSigned
4. Tape: Y et Enter
5. Réessaye le script
```

### Problème: Port 8080 déjà utilisé
```
Solution A: Ferme l'autre application
Solution B: Change le port dans web.go (voir section Personnaliser)
```

### Problème: Fichiers web manquants
```
Solution:
1. Assure-toi que dossier "web" existe
2. Avec fichiers: index.html, style.css, script.js
3. Réessaye
```

---

## 📊 Structure pour Windows

```
ia-atomique/
├── install.bat          ← À double-cliquer en 1er
├── run-web.bat          ← À double-cliquer ensuite
├── build.bat
├── verify.bat
│
├── install.ps1          ← Alternative PowerShell
├── run-web.ps1
├── build.ps1
│
├── WINDOWS_INSTALL.md   ← Tu es ici!
├── WINDOWS_SETUP.md
│
├── programme.exe        ← Créé après install.bat
├── web/
│   ├── index.html
│   ├── style.css
│   └── script.js
│
└── ... autres fichiers
```

---

## 🎯 Workflow Typique Windows

```
1️⃣ Double-clique install.bat
   └─ Affiche: "✅ Installation Réussie!"

2️⃣ Double-clique run-web.bat
   └─ Lance le serveur

3️⃣ Une fenêtre CMD s'ouvre avec:
   - "Serveur Web Démarré"
   - "URL: http://localhost:8080"

4️⃣ Ouvre navigateur: http://localhost:8080

5️⃣ Utilise l'interface web! 🎉

6️⃣ Pour arrêter: Ferme la fenêtre CMD
```

---

## ✅ Points Importants

✅ **Aucune ligne de commande requise** - Juste double-clique  
✅ **Automatique** - Install.bat fait tout tout seul  
✅ **Erreurs claires** - Sait dire exactement ce qui manque  
✅ **PowerShell aussi** - Pour ceux qui préfèrent  
✅ **Portable** - Pas de dépendances externes  

---

## 🚀 Prochaines Étapes

### Pour toi (développeur):
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
5. C'est tout! 🎉

---

## 📞 Support Windows

Si quelque chose marche pas:
1. Relance Windows complètement
2. Réinstalle Go proprement
3. Double-clique `verify.bat` pour vérifier
4. Essaye à nouveau

---

**Voilà! Installation super simple pour Windows! 🎉**
