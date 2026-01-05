# 🚂 Déploiement sur Railway.app

Guide complet pour héberger IA-ATOMIQUE sur Railway.

## 🚀 Setup Rapide (5 minutes)

### 1️⃣ Créer un compte Railway
- Va sur: https://railway.app
- Clique "Sign Up"
- Connecte-toi avec GitHub (recommandé)

### 2️⃣ Créer un nouveau projet
- Clique "New Project"
- Choisis "Deploy from GitHub repo"
- Sélectionne ton repo `IA-ATOMIQUE`

### 3️⃣ Railway détecte automatiquement
- Railway voit `go.mod` et `Procfile`
- Compile automatiquement
- Lance le serveur

### 4️⃣ Attendre le déploiement
- Quelques secondes/minutes
- Tu vois l'URL en haut:
  ```
  https://ton-projet-xxx.railway.app
  ```

### 5️⃣ Ouvrir l'URL
```
https://ton-projet-xxx.railway.app
```

C'est tout! 🎉

---

## 📋 Fichiers Nécessaires

✅ `Procfile` - Comment lancer l'app  
✅ `railway.json` - Configuration Railway  
✅ `go.mod` - Dépendances Go  
✅ `web/` - Fichiers statiques  

**Tous sont dans le repo!**

---

## 🔧 Configuration Avancée

### Voir les logs
```bash
railway logs
```

### Variables d'environnement
Railway ajoute automatiquement `PORT`
(ton code l'utilise déjà!)

### Redéployer
```bash
git push origin main
# Railway détecte et redéploie auto
```

---

## 🌐 URL Finale

Une fois déployée, ton app sera accessible sur:
```
https://ia-atomique-xxx.railway.app
```

Remplace `xxx` par le nom que Railway génère.

---

## 🐛 Troubleshooting

### Erreur: "Build failed"
- Vérifier que `go.mod` existe
- Vérifier que `Procfile` existe
- Vérifier les logs: `railway logs`

### Port error
- Railway assigne le port automatiquement
- Notre code utilise `os.Getenv("PORT")`
- C'est bon!

### Fichiers web manquants
- Assurer que `web/` existe avec les fichiers
- Vérifier que c'est commit sur GitHub

---

## 📊 Performance

- Go est ultra-rapide
- Railway gratuit peut gérer beaucoup de traffic
- Le serveur tourne 24/7

---

## 💰 Coût

- Gratuit les premiers crédits ($5/mois)
- Puis payant (très peu cher)

---

## ✨ Avantages Railway

✅ Déploiement automatique (git push)  
✅ Zéro configuration requise  
✅ Logs en temps réel  
✅ Domaine inclus  
✅ SSL HTTPS automatique  
✅ Redémarrage auto si crash  

---

## 🎯 Prochaines Étapes

1. Push ton code sur GitHub
2. Va sur railway.app
3. Connecte ton repo
4. Attends 2 minutes
5. Ton serveur est live!

---

**C'est simpler que ça en a l'air!** 🚀

Besoin d'aide? Dis-moi!
