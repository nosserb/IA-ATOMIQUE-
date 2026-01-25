#  Quickstart - Generation d'Images Atomiques

## Installation & Compilation

```bash
cd /home/student/autre\ projets/IA-ATOMIQUE-
go build -o programme
```

## Premiers pas

### 1� Generation simple

```bash
./programme image generate 256 256 50 8 "blue sky with clouds"
```

**Resultat** : `generated_image.png`
- Dimensions : 256�256 pixels
- Iterations : 50 (plus = plus detaille)
- Patch size : 8 (8�8 pixels par atome)
- Resultat en ~50-100ms

### 2� Generation auto-optimisee

```bash
./programme image prompt "beautiful sunset over ocean with mountains"
```

Le syst�me choisit automatiquement les meilleurs param�tres selon la complexite du prompt.

### 3� Generation haute qualite multi-echelle

```bash
./programme image multi-scale "detailed fantasy landscape with castles"
```

Gen�re � 3 echelles differentes pour meilleure qualite et plus de details.

## Exemples de prompts

### Paysages
```bash
./programme image prompt "mountain landscape with snow and pine trees"
./programme image prompt "desert sunset with sand dunes"
./programme image prompt "tropical beach with palm trees and clear ocean"
```

### Abstraits
```bash
./programme image prompt "abstract colorful art waves and energy"
./programme image prompt "geometric patterns bright and sharp"
./programme image prompt "rough dark textured stone surface"
```

### Couleurs specifiques
```bash
./programme image prompt "red roses on dark background"
./programme image prompt "purple mountains under pink sky"
./programme image prompt "blue ocean green grass yellow sun"
```

### Textures
```bash
./programme image prompt "smooth silk texture light and reflective"
./programme image prompt "rough stone wall detailed and sharp"
./programme image prompt "blurry misty foggy forest"
```

## Param�tres expliques

```bash
./programme image generate WIDTH HEIGHT ITERATIONS PATCH_SIZE "prompt"
```

| Param�tre | Plage | Par defaut | Effet |
|-----------|-------|-----------|-------|
| WIDTH | 64-2048 | - | Largeur en pixels |
| HEIGHT | 64-2048 | - | Hauteur en pixels |
| ITERATIONS | 1-1000 | - | Plus = plus detaille et lent |
| PATCH_SIZE | 1-64 | - | Plus = plus rapide mais moins detaille |
| PROMPT | Texte | - | Description de l'image |

### Recommandations

**Rapide (< 100ms)** :
```bash
./programme image generate 256 256 30 16 "prompt"  # Patch 16�16
```

**Normal (100-500ms)** :
```bash
./programme image generate 512 512 50 8 "prompt"   # Patch 8�8
```

**Haute qualite (1-5s)** :
```bash
./programme image generate 1024 1024 100 8 "prompt" # Grand + iterations
```

## Resultats typiques

### Apr�s 30 iterations
```
�tat moyenne: 0.456
Intensite moyenne: 0.389
Atomes actifs: 78.5%
```
 Image brute, couleurs principales presentes

### Apr�s 50 iterations
```
�tat moyenne: 0.512
Intensite moyenne: 0.445
Atomes actifs: 62.3%
Atomes geles: 37.7%
```
 Image stabilisee, details emergents

### Apr�s 100 iterations
```
�tat moyenne: 0.534
Intensite moyenne: 0.478
Atomes actifs: 45.2%
Atomes geles: 54.8%
```
 Image compl�te et stable, beaucoup d'atomes en hibernation (efficace)

## Troubleshooting

### Image toute noire
**Cause** : Seed aleatoire defavorable ou contraintes trop fortes
**Solution** : Relancer (seed aleatoire change), ou reduire ITERATIONS

### Image trop claire/brillante
**Cause** : Prompt contient "bright" ou "light"
**Solution** : Ajouter "dark", "night", ou "dim" pour equilibrer

### Generation lente
**Cause** : PATCH_SIZE trop grand, ITERATIONS trop elevees
**Solution** : Augmenter PATCH_SIZE de 8 � 16 ou 32

### Pas de couleur
**Cause** : Couleurs non reconnues dans le prompt
**Solution** : Utiliser : red, blue, green, yellow, purple, orange, pink, cyan

## Architecture du reseau atomique

```
Initialisation aleatoire
         
    [Atomes]
          (iteration 1)
Resonance locale
          (iteration 2)
Contraintes du prompt
          (iteration N)
Image coherente
```

Chaque atome influence ses 8 voisins localement, creant une image par emergence.

## Performance

### Timing empirique

```
256�256, 50 iter, patch 8:   ~50ms
512�512, 100 iter, patch 8:  ~200ms
1024�1024, 100 iter, patch 8: ~800ms
```

### Comparaison

| Syst�me | Temps | Memoire | CPU |
|---------|-------|---------|-----|
| T.R.A. | 50-1000ms | 1-100 MB | Leger |
| Stable Diffusion | 30-60s | 10-20 GB | GPU |
| DALL-E | 30-60s | Cloud | Cloud |

## Fichiers generes

```
generated_image.png       # Image finale (PNG)
generated_multiscale.png  # Pour multi-scale
```

## Exemples complets

### Coucher de soleil
```bash
./programme image generate 512 512 100 8 "red orange sunset over blue ocean waves"
```

### For�t mysterieuse
```bash
./programme image prompt "dark mysterious forest misty foggy trees"
```

### Art abstrait
```bash
./programme image multi-scale "abstract art colorful sharp patterns chaos"
```

### Textures realistes
```bash
./programme image generate 256 256 80 4 "rough sharp detailed granite stone texture"
```

## Fichiers cles du projet

```
database/
   image_atomic.go        # Moteur d'image atomique
      PixelAtom
      AtomicImageNetwork
      Fonctions de generation
   atomic.go              # Moteur textuel (existant)

image_commands.go            # Interface CLI pour images
IMAGE-GENERATION-GUIDE.md    # Documentation compl�te
```

## Prochains pas

1. **Essayer les exemples** : Coucher de soleil, paysage, abstrait
2. **Optimiser les param�tres** : Trouver votre equilibre vitesse/qualite
3. **Explorer les prompts** : Tester les couleurs, textures, lumi�re
4. **Mixer texte+image** : Utiliser le resume atomique pour l'image

## Support

Voir `IMAGE-GENERATION-GUIDE.md` pour la documentation technique compl�te.
