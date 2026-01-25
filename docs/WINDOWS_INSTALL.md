# IA-ATOMIQUE sur Windows

Guide d'installation simple pour Windows.

## Installation Rapide (2 étapes)

### πtape 1: Installer Go (une seule fois)

1. **Télécharge Go** depuis https://golang.org/dl/
2. **Choisir la version Windows** (msi):
   - Windows x86-64 (64-bit) = La plupart des gens
   - Windows x86 (32-bit) = PC ancien
3. **Exécute l'installateur** et clique "Next" partout
4. **Redémarre** (optionnel mais recommandé)

**Vérifier que Go est installé:**
- Ouvre PowerShell ou Command Prompt
- Tape: `go version`
- Tu dois voir quelque chose comme: `go version go1.22.2 windows/amd64`

Si πa marche pas, redémarre ton PC et réessaye.

---

### πtape 2: Lancer le projet

1. **Clone le repo GitHub**
   ```
   git clone https://github.com/TON-USER/ia-atomique.git
   cd ia-atomique
   ```

2. **Double-clique sur `install.bat`**
   - Ça compile automatiquement
   - Ça vérifie les fichiers
   - C'est tout!

3. **Double-clique sur `run-web.bat`**
   - Le serveur démarre
   - Une fenπtre s'ouvre

4. **Ouvre ton navigateur**
   - Accπde π: http://localhost:8080
   - Et voilπ! 

---

## Fichiers Disponibles

### Scripts Windows

| Fichier | Fonction |
|---------|----------|
| `install.bat` | Installation initiale (π faire une seule fois) |
| `run-web.bat` | Lancer le serveur web |
| `build.bat` | Recompiler le projet |
| `verify.bat` | Vérifier que tout est OK |

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

## Mode Command Prompt/PowerShell

Si tu préfπres la ligne de commande:

### Premiπre utilisation
```bash
# Ouvre Command Prompt (Cmd.exe)
# Navigue vers le répertoire:
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

## Résolution des Problπmes

### "Go is not recognized"

**Problπme:** Go n'est pas installé ou trouvable

**Solution:**
1. Installe Go depuis https://golang.org/dl/
2. Redémarre ton PC complπtement
3. Réessaye

### "web\index.html n'existe pas"

**Problπme:** Fichiers web manquants

**Solution:**
- Assure-toi que le dossier `web/` existe
- Avec les fichiers: `index.html`, `style.css`, `script.js`

### Port 8080 déjπ utilisé

**Problπme:** Un autre programme utilise le port 8080

**Solution 1:**
- Ferme l'autre application

**Solution 2:**
- Modifie `web.go` ligne 29:
  ```go
  go StartWebServer("9000")  // Change 8080 en 9000
  ```
- Recompile: `build.bat`
- Lance avec: `programme.exe web`
- Accπde π: http://localhost:9000

### "programme.exe n'existe pas"

**Problπme:** Pas compilé

**Solution:**
- Double-clique `build.bat` ou `install.bat`

---

## Astuces Windows

### Créer un raccourci sur le Bureau

1. **Clique droit sur `run-web.bat`**
2. **Envoyer vers > Bureau (créer un raccourci)**
3. **Double-clique le raccourci pour lancer!**

### Lancer dans une vraie fenπtre (pas CMD)

Pour une expérience plus fluide:

1. **Crée `run-web.vbs`:**

```vbs
Set objShell = CreateObject("WScript.Shell")
objShell.Run "cmd /c run-web.bat", 0
WScript.Quit
```

2. **Sauvegarde dans le mπme dossier**
3. **Double-clique pour lancer sans fenπtre CMD**

### Ajouter au menu Démarrer

1. **Appuie sur Win+R**
2. **Tape: `shell:startup`**
3. **Mets un raccourci de `run-web.bat` lπ**
4. **π chaque démarrage, le serveur se lance!**

---

## Workflow Typique

```
1. Double-clique install.bat (premiπre fois seulement)
   
2. Double-clique run-web.bat
   
3. Une fenπtre s'ouvre avec le serveur
   
4. Ouvre navigateur: http://localhost:8080
   
5. Utilise l'interface!
   
6. Pour arrπter: Ferme la fenπtre CMD
```

---

## Pour Développeurs

Si tu modifies le code:

```bash
# Compile
build.bat

# Teste
run-web.bat
```

---

## Avant de Pusher sur GitHub

```bash
# Vérifier que tout compile
verify.bat

# Si tout est OK, push!
git add .
git commit -m "feat: Add web interface"
git push origin main
```

---

## π Support

Si quelque chose ne marche pas:

1. **Vérifie que Go est installé:** `go version`
2. **Recompile:** `build.bat`
3. **Vérification:** `verify.bat`
4. **Réessaye:** `run-web.bat`

Si πa marche toujours pas:
- Redémarre ton PC
- Réinstalle Go proprement
- Essaye π nouveau

---

**C'est simple! Bon résumé! **
