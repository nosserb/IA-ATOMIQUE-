# Quickstart - Génération d'Images Atomiques

## Installation & Compilation

```bash
cd /home/student/autre\ projets/IA-ATOMIQUE-
go build -o programme
```

## Premiers pas

### 1 Génération simple

```bash
./programme image generate 256 256 50 8 "blue sky with clouds"
```

**Résultat** : `generated_image.png`
- Dimensions : 256π256 pixels
- Itérations : 50 (plus = plus détaillé)
- Patch size : 8 (8π8 pixels par atome)
- Résultat en ~50-100ms

### 2 Génération auto-optimisée

```bash
./programme image prompt "beautiful sunset over ocean with mountains"
```

Le systπme choisit automatiquement les meilleurs paramπtres selon la complexité du prompt.

### 3 Génération haute qualité multi-échelle

```bash
./programme image multi-scale "detailed fantasy landscape with castles"
```

Génπre π 3 échelles différentes pour meilleure qualité et plus de détails.

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

### Couleurs spécifiques
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

## Paramπtres expliqués

```bash
./programme image generate WIDTH HEIGHT ITERATIONS PATCH_SIZE "prompt"
```

| Paramπtre | Plage | Par défaut | Effet |
|-----------|-------|-----------|-------|
| WIDTH | 64-2048 | - | Largeur en pixels |
| HEIGHT | 64-2048 | - | Hauteur en pixels |
| ITERATIONS | 1-1000 | - | Plus = plus détaillé et lent |
| PATCH_SIZE | 1-64 | - | Plus = plus rapide mais moins détaillé |
| PROMPT | Texte | - | Description de l'image |

### Recommandations

**Rapide (< 100ms)** :
```bash
./programme image generate 256 256 30 16 "prompt"  # Patch 16π16
```

**Normal (100-500ms)** :
```bash
./programme image generate 512 512 50 8 "prompt"   # Patch 8π8
```

**Haute qualité (1-5s)** :
```bash
./programme image generate 1024 1024 100 8 "prompt" # Grand + itérations
```

## Résultats typiques

### Aprπs 30 itérations
```
πtat moyenne: 0.456
Intensité moyenne: 0.389
Atomes actifs: 78.5%
```
 Image brute, couleurs principales présentes

### Aprπs 50 itérations
```
πtat moyenne: 0.512
Intensité moyenne: 0.445
Atomes actifs: 62.3%
Atomes gelés: 37.7%
```
 Image stabilisée, détails émergents

### Aprπs 100 itérations
```
πtat moyenne: 0.534
Intensité moyenne: 0.478
Atomes actifs: 45.2%
Atomes gelés: 54.8%
```
 Image complπte et stable, beaucoup d'atomes en hibernation (efficace)

## Troubleshooting

### Image toute noire
**Cause** : Seed aléatoire défavorable ou contraintes trop fortes
**Solution** : Relancer (seed aléatoire change), ou réduire ITERATIONS

### Image trop claire/brillante
**Cause** : Prompt contient "bright" ou "light"
**Solution** : Ajouter "dark", "night", ou "dim" pour équilibrer

### Génération lente
**Cause** : PATCH_SIZE trop grand, ITERATIONS trop élevées
**Solution** : Augmenter PATCH_SIZE de 8 π 16 ou 32

### Pas de couleur
**Cause** : Couleurs non reconnues dans le prompt
**Solution** : Utiliser : red, blue, green, yellow, purple, orange, pink, cyan

## Architecture du réseau atomique

```
Initialisation aléatoire
         
    [Atomes]
          (itération 1)
Résonance locale
          (itération 2)
Contraintes du prompt
          (itération N)
Image cohérente
```

Chaque atome influence ses 8 voisins localement, créant une image par émergence.

## Performance

### Timing empirique

```
256π256, 50 iter, patch 8:   ~50ms
512π512, 100 iter, patch 8:  ~200ms
1024π1024, 100 iter, patch 8: ~800ms
```

### Comparaison

| Systπme | Temps | Mémoire | CPU |
|---------|-------|---------|-----|
| T.R.A. | 50-1000ms | 1-100 MB | Léger |
| Stable Diffusion | 30-60s | 10-20 GB | GPU |
| DALL-E | 30-60s | Cloud | Cloud |

## Fichiers générés

```
generated_image.png       # Image finale (PNG)
generated_multiscale.png  # Pour multi-scale
```

## Exemples complets

### Coucher de soleil
```bash
./programme image generate 512 512 100 8 "red orange sunset over blue ocean waves"
```

### Forπt mystérieuse
```bash
./programme image prompt "dark mysterious forest misty foggy trees"
```

### Art abstrait
```bash
./programme image multi-scale "abstract art colorful sharp patterns chaos"
```

### Textures réalistes
```bash
./programme image generate 256 256 80 4 "rough sharp detailed granite stone texture"
```

## Fichiers clés du projet

```
database/
   image_atomic.go        # Moteur d'image atomique
      PixelAtom
      AtomicImageNetwork
      Fonctions de génération
   atomic.go              # Moteur textuel (existant)

image_commands.go            # Interface CLI pour images
IMAGE-GENERATION-GUIDE.md    # Documentation complπte
```

## Prochains pas

1. **Essayer les exemples** : Coucher de soleil, paysage, abstrait
2. **Optimiser les paramπtres** : Trouver votre équilibre vitesse/qualité
3. **Explorer les prompts** : Tester les couleurs, textures, lumiπre
4. **Mixer texte+image** : Utiliser le résumé atomique pour l'image

## Support

Voir `IMAGE-GENERATION-GUIDE.md` pour la documentation technique complπte.
