# 📖 GUIDE D'UTILISATION - SYSTÈME DE RÉSUMÉ ADAPTATIF

## 🎯 Concept Clé

Le système détecte automatiquement le **type de texte** et adapte le pipeline de résumé en conséquence. Vous obtenez des résumés optimisés sans configuration supplémentaire.

## 📊 3 Types de Textes Supportés

### 1. **ENCYCLOPÉDIQUE** 🔬
Textes factuels, scientifiques, informatifs

**Caractéristiques détectées:**
- Définitions ("est un", "est une", "est le")
- Processus ("photosynthèse", "réaction")
- Structures ("molécule", "chloroplaste", "organite")
- Équations et formules scientifiques

**Comportement du système:**
```
✅ Résume par PHRASES (pas par mots)
✅ Compression limité max 30%
✅ SKIP Phase X+1 (conserve faits concrets)
✅ SKIP enrichissement (garde structure pure)
✅ SKIP humanisation (pas de connecteurs parasites)
✅ Phase X+4 optionnelle si ultra-compression
```

**Qualité obtenue:**
- Cohérence maximale
- Lisibilité excellente
- Faits préservés
- Pas d'hallucinations IA

**Exemples:** Photosynthèse, biologie, chimie, physique, géographie

---

### 2. **NARRATIF** 📖
Histoires, récits, expériences, témoignages

**Caractéristiques détectées:**
- Personnages ("héros", "protagoniste")
- Dialogues ("dit", "demanda", "s'exclama")
- Actions ("aventure", "voyage", "conflit")
- Ambiance temporelle ("autrefois", "soudain", "alors")

**Comportement du système:**
```
✓ Résume par mots (plus flexible)
✓ Respecte ratio utilisateur (0.1 → 0.9)
✓ Phase X+1 appliquée si abstraction < 60%
✓ Enrichissement normal
✓ Humanisation normale (Phase X+3)
```

**Qualité obtenue:**
- Rythme préservé
- Tension narrative maintenue
- Émotions transcrites
- Cohérence acceptable

**Exemples:** Romans, contes, mémoires, témoignages

---

### 3. **CONCEPTUEL** 💡
Philosophie, analyse, théorie, argumentation

**Caractéristiques détectées:**
- Concepts ("idée", "principe", "théorie")
- Logique ("cause", "conséquence", "donc")
- Abstraction ("philosophie", "métaphysique")
- Argumentation ("par contre", "néanmoins", "cependant")

**Comportement du système:**
```
✓ Résume par mots (compression agressive OK)
✓ Accepte ratio élevé (0.7-0.9)
✓ Phase X+1 forçée si abstraction < 60%
✓ Humanisation Phase X+3 appliquée
✓ Convertit concret → abstrait
```

**Qualité obtenue:**
- Abstraction forcée
- Essence conservée
- Connecteurs naturels
- Argumention simplifiée

**Exemples:** Philosophie, essais, analyses politiques, théories

---

## 🛠️ UTILISATION RECOMMANDÉE

### Utilisation Simple (RECOMMANDÉE)
```bash
./programme resume votre_document.txt 0.5
```
Le système détecte le type et adapte automatiquement! 

**Résultat selon type:**
- Encyclopédique + ratio 0.1 → 30% gardé automatiquement (protection)
- Narratif + ratio 0.5 → 50% gardé (respecté)
- Conceptuel + ratio 0.7 → 30% gardé (agressif, acceptable)

### Utilisation Avancée (PRÉDICTION)

**Pour textes encyclopédiques:**
```bash
# Vous voullez 15% → demandez 0.15 → reçoit 30% (meilleur compromis)
./programme resume physique.txt 0.15

# Vous voullez 40% → demandez 0.4 → reçoit 30-40%
./programme resume chimie.txt 0.4
```

**Pour textes narratifs:**
```bash
# Respecte votre ratio exactement
./programme resume roman.txt 0.3     # 30% gardé
./programme resume histoire.txt 0.7  # 70% gardé
```

**Pour textes conceptuels:**
```bash
# Permet compression agressive
./programme resume philo.txt 0.5      # Compression normale
./programme resume essai.txt 0.9      # Ultra-compression acceptable
```

---

## 📈 INTERPRÉTATION DU RAPPORT

### Exemple 1: Texte Encyclopédique

```
[PHASE 15] Étape 0: Détection du type de texte...
  ℹ️  Type: Encyclopédique (compression: 0.10)
  ⚠️  Phase X+1 désactivée pour ce type

[PHASE 15] Étape 2: Résumé atomique (Phase 13+++)...
  ⚠️  Compression limitée à 30% pour texte encyclopédique (demandée: 10%)
  ✓ Résumé généré: 1847 caractères
```

**Interprétation:**
- ✅ Type correct identifié
- ℹ️ Phase X+1 skippée (fait correct, pas d'abstraction)
- ⚠️ Vous avez demandé 10% mais reçu 30% → Protection appliquée
- ✓ Résumé bien généré

**En bas du rapport:**
```
Text Type:         Encyclopédique
Abstraction:       SKIPPED (facts preservation)
Compression:       82.4%
```
✅ Tout est correct

---

### Exemple 2: Texte Narratif

```
[PHASE 15] Étape 0: Détection du type de texte...
  ℹ️  Type: Narratif (compression: 0.50)

[PHASE 15] Étape 2: Résumé atomique...
  ✓ Résumé généré: 12500 caractères
```

**Interprétation:**
- ✅ Type narratif détecté
- Pas de limite appliquée (ratio respecté)
- Abstraction appliquée si nécessaire

---

## ⚠️ MESSAGES SPÉCIAUX

### "⚠️ Phase X+1 désactivée pour ce type"
**Signifie:** Texte encyclopédique détecté → faits conservés → c'est normal et bon

### "⚠️ Compression limitée à 30%"
**Signifie:** Vous avez demandé < 30% sur encyclopédique → protection automatique appliquée

### "ℹ️ Texte encyclopédique: conservation des faits concrets"
**Signifie:** Phase X+1 skippée → pas de forçage d'abstraction

### "ℹ️ Texte encyclopédique: conservation structure originale"
**Signifie:** Phase X+3 skippée → pas de connecteurs parasites ajoutés

---

## 🎬 CAS D'USAGE PRATIQUES

### Cas 1: Résumer un article scientifique
```bash
./programme resume article_biologie.txt 0.2
```
- Détection: ENCYCLOPÉDIQUE
- Limite: 30% gardé (pas 20%)
- Résultat: Article cohérent et factuel

### Cas 2: Résumer un roman
```bash
./programme resume roman_250pages.txt 0.1
```
- Détection: NARRATIF
- Limite: Aucune (10% respecté)
- Résultat: Résumé court du rythme narratif

### Cas 3: Résumer un essai philosophique
```bash
./programme resume essai_politique.txt 0.5
```
- Détection: CONCEPTUEL
- Limite: Aucune (50% respecté)
- Résultat: Essence conceptuelle extraite

---

## 🚀 CONSEILS D'OPTIMISATION

### Pour Encyclopédique
- Demandez 30-40% pour meilleur résultat
- Ne demandez pas < 20% (perte importante)
- Le système corrigera automatiquement

### Pour Narratif
- Demandez 30-50% pour bonne cohérence
- Demandez 10-20% pour résumé très court
- Demandez 70-90% pour simple réduction

### Pour Conceptuel
- Demandez 50-70% pour abstraction modérée
- Demandez 80-90% pour ultra-dense
- Demandez 30-40% pour préserver plus de détails

---

## 📊 METRIQUES À SURVEILLER

### Compression
- < -10%: Expansion (texte enrichi, normal pour encyclopédique)
- 0-50%: Modéré (bon compromis)
- 50-80%: Agressif (synopsis)
- > 80%: Ultra-dense (une phrase)

### Coherence Score
- > 50%: Excellent
- 20-50%: Bon
- < 20%: À améliorer (trop agressif)

### Grammar Score
- > 75%: Excellent
- 50-75%: Bon
- < 50%: Problématique

---

## ❓ DÉPANNAGE

### "Le résumé est trop court"
→ Augmentez le ratio (ex: 0.3 → 0.5)
→ Ou laissez le système adapter pour encyclopédique

### "Le résumé garde trop de contenu"
→ Diminuez le ratio (ex: 0.7 → 0.5)
→ Conceptuel accepte 0.9, encyclopédique a limite 0.3

### "Beaucoup de connecteurs inutiles"
→ C'est Narratif/Conceptuel (normal)
→ Encyclopédique a protection contre ça

### "Perte de faits importants"
→ Vérifiez que c'est détecté ENCYCLOPÉDIQUE
→ Augmentez ratio si < 0.3

---

## ✅ CHECKLIST UTILISATION

- [ ] Vous avez un document à résumer
- [ ] Vous avez deviné le type (encyclopédique/narratif/conceptuel)
- [ ] Vous avez choisi un ratio de compression
- [ ] Vous avez lancé: `./programme resume document.txt ratio`
- [ ] Vous avez lu le rapport (cherchez les ⚠️ et ℹ️)
- [ ] Vous avez comparé avec l'original si doubtful
- [ ] Vous avez ajusté le ratio si nécessaire et relancé

---

## 🎓 CONCLUSION

Le système est conçu pour vous donner:
1. **Automatisation**: Détection type sans effort
2. **Qualité**: Adaptation à chaque type
3. **Protection**: Limites intelligentes contre mauvaise compression
4. **Contrôle**: Vous fixez toujours le ratio
5. **Transparence**: Rapport détaillé explique ce qui s'est passé

**Conseil:** Utilisez ratios 0.3-0.7 pour meilleur résultat. Laissez le système adapter pour encyclopédique.

