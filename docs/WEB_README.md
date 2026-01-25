# IA-ATOMIQUE - Interface Web

Interface web simple pour résumer des textes avec votre IA neuronale.

## Démarrage

### Option 1: Avec le script (recommandé)
```bash
./start-web.sh
```

### Option 2: Directement avec Go
```bash
go run main.go interaction.go web.go web_utils.go database/*.go web
```

### Option 3: Avec le binaire compilé
```bash
./programme web
```

## Accàs

Une fois le serveur démarré, ouvrez votre navigateur et accédez à:
```
http://localhost:8080
```

## Utilisation

1. **Collez votre texte** dans la zone "Votre Texte"
2. **Cliquez sur "Résumer"** ou appuyez sur **Ctrl+Entrée**
3. **Consultez le résumé** dans la zone de droite

### Résumé généré
Le résumé affiche:
-  Les catégories principales détectées
-  Les phrases clés du texte
- Statistiques d'analyse (phrases, confiance, catégories)

## Personnalisation

### Modifier les couleurs
àditez `web/style.css`:
```css
/* Modifier le gradient principal */
background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
```

### Modifier la taille des éléments
Recherchez les valeurs de `padding`, `font-size`, `height` dans `style.css`

## Fichiers

- `web/index.html` - Structure HTML
- `web/style.css` - Styles et mise en page
- `web/script.js` - Logique frontend
- `web.go` - Serveur HTTP et API
- `start-web.sh` - Script de lancement

## API

### Endpoint: POST /api/summarize

**Request:**
```json
{
  "text": "Votre texte à résumer..."
}
```

**Response:**
```json
{
  "summary": "Résumé généré...",
  "stats": {
    "phrases": 5,
    "confidence": "85.3",
    "categories": 3
  }
}
```

## Troubleshooting

### Port 8080 déjà utilisé
Modifiez le port dans `web.go` ligne:
```go
go StartWebServer("9000")  // Changez 8080 en 9000
```

### Fichiers web non trouvés
Assurez-vous que le répertoire `web/` existe avec les fichiers:
- `index.html`
- `style.css`
- `script.js`

## Améliorations possibles

- [ ] Support des fichiers upload
- [ ] Historique des résumés
- [ ] Export en PDF
- [ ] Mode sombre
- [ ] Traduction du résumé
- [ ] Graphiques d'analyse

Bon résumé! 
