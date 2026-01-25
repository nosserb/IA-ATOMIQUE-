# IA-ATOMIQUE - Interface Web

Interface web simple pour resumer des textes avec votre IA neuronale.

##  Demarrage

### Option 1: Avec le script (recommande)
```bash
./start-web.sh
```

### Option 2: Directement avec Go
```bash
go run main.go interaction.go web.go web_utils.go database/*.go web
```

### Option 3: Avec le binaire compile
```bash
./programme web
```

##  Acc�s

Une fois le serveur demarre, ouvrez votre navigateur et accedez �:
```
http://localhost:8080
```

##  Utilisation

1. **Collez votre texte** dans la zone "Votre Texte"
2. **Cliquez sur "Resumer"** ou appuyez sur **Ctrl+Entree**
3. **Consultez le resume** dans la zone de droite

### Resume genere
Le resume affiche:
-  Les categories principales detectees
-  Les phrases cles du texte
- Statistiques d'analyse (phrases, confiance, categories)

##  Personnalisation

### Modifier les couleurs
�ditez `web/style.css`:
```css
/* Modifier le gradient principal */
background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
```

### Modifier la taille des elements
Recherchez les valeurs de `padding`, `font-size`, `height` dans `style.css`

##  Fichiers

- `web/index.html` - Structure HTML
- `web/style.css` - Styles et mise en page
- `web/script.js` - Logique frontend
- `web.go` - Serveur HTTP et API
- `start-web.sh` - Script de lancement

##  API

### Endpoint: POST /api/summarize

**Request:**
```json
{
  "text": "Votre texte � resumer..."
}
```

**Response:**
```json
{
  "summary": "Resume genere...",
  "stats": {
    "phrases": 5,
    "confidence": "85.3",
    "categories": 3
  }
}
```

##  Troubleshooting

### Port 8080 dej� utilise
Modifiez le port dans `web.go` ligne:
```go
go StartWebServer("9000")  // Changez 8080 en 9000
```

### Fichiers web non trouves
Assurez-vous que le repertoire `web/` existe avec les fichiers:
- `index.html`
- `style.css`
- `script.js`

##  Ameliorations possibles

- [ ] Support des fichiers upload
- [ ] Historique des resumes
- [ ] Export en PDF
- [ ] Mode sombre
- [ ] Traduction du resume
- [ ] Graphiques d'analyse

Bon resume! 
