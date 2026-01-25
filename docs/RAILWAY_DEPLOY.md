#  Deploiement sur Railway.app

Guide complet pour heberger IA-ATOMIQUE sur Railway.

##  Setup Rapide (5 minutes)

### 1� Creer un compte Railway
- Va sur: https://railway.app
- Clique "Sign Up"
- Connecte-toi avec GitHub (recommande)

### 2� Creer un nouveau projet
- Clique "New Project"
- Choisis "Deploy from GitHub repo"
- Selectionne ton repo `IA-ATOMIQUE`

### 3� Railway detecte automatiquement
- Railway voit `go.mod` et `Procfile`
- Compile automatiquement
- Lance le serveur

### 4� Attendre le deploiement
- Quelques secondes/minutes
- Tu vois l'URL en haut:
  ```
  https://ton-projet-xxx.railway.app
  ```

### 5� Ouvrir l'URL
```
https://ton-projet-xxx.railway.app
```

C'est tout! 

---

## � Fichiers Necessaires

 `Procfile` - Comment lancer l'app  
 `railway.json` - Configuration Railway  
 `go.mod` - Dependances Go  
 `web/` - Fichiers statiques  

**Tous sont dans le repo!**

---

##  Configuration Avancee

### Voir les logs
```bash
railway logs
```

### Variables d'environnement
Railway ajoute automatiquement `PORT`
(ton code l'utilise dej�!)

### Redeployer
```bash
git push origin main
# Railway detecte et redeploie auto
```

---

##  URL Finale

Une fois deployee, ton app sera accessible sur:
```
https://ia-atomique-xxx.railway.app
```

Remplace `xxx` par le nom que Railway gen�re.

---

##  Troubleshooting

### Erreur: "Build failed"
- Verifier que `go.mod` existe
- Verifier que `Procfile` existe
- Verifier les logs: `railway logs`

### Port error
- Railway assigne le port automatiquement
- Notre code utilise `os.Getenv("PORT")`
- C'est bon!

### Fichiers web manquants
- Assurer que `web/` existe avec les fichiers
- Verifier que c'est commit sur GitHub

---

##  Performance

- Go est ultra-rapide
- Railway gratuit peut gerer beaucoup de traffic
- Le serveur tourne 24/7

---

##  Co�t

- Gratuit les premiers credits ($5/mois)
- Puis payant (tr�s peu cher)

---

##  Avantages Railway

 Deploiement automatique (git push)  
 Zero configuration requise  
 Logs en temps reel  
 Domaine inclus  
 SSL HTTPS automatique  
 Redemarrage auto si crash  

---

##  Prochaines �tapes

1. Push ton code sur GitHub
2. Va sur railway.app
3. Connecte ton repo
4. Attends 2 minutes
5. Ton serveur est live!

---

**C'est simpler que �a en a l'air!** 

Besoin d'aide? Dis-moi!
