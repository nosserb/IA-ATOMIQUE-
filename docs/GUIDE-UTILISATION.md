#  GUIDE D'UTILISATION - SYST�ME DE R�SUM� ADAPTATIF

##  Concept Cle

Le syst�me detecte automatiquement le **type de texte** et adapte le pipeline de resume en consequence. Vous obtenez des resumes optimises sans configuration supplementaire.

##  3 Types de Textes Supportes

### 1. **ENCYCLOP�DIQUE** 
Textes factuels, scientifiques, informatifs

**Caracteristiques detectees:**
- Definitions ("est un", "est une", "est le")
- Processus ("photosynth�se", "reaction")
- Structures ("molecule", "chloroplaste", "organite")
- �quations et formules scientifiques

**Comportement du syst�me:**
```
 Resume par PHRASES (pas par mots)
 Compression limite max 30%
 SKIP Phase X+1 (conserve faits concrets)
 SKIP enrichissement (garde structure pure)
 SKIP humanisation (pas de connecteurs parasites)
 Phase X+4 optionnelle si ultra-compression
```

**Qualite obtenue:**
- Coherence maximale
- Lisibilite excellente
- Faits preserves
- Pas d'hallucinations IA

**Exemples:** Photosynth�se, biologie, chimie, physique, geographie

---

### 2. **NARRATIF** 
Histoires, recits, experiences, temoignages

**Caracteristiques detectees:**
- Personnages ("heros", "protagoniste")
- Dialogues ("dit", "demanda", "s'exclama")
- Actions ("aventure", "voyage", "conflit")
- Ambiance temporelle ("autrefois", "soudain", "alors")

**Comportement du syst�me:**
```
 Resume par mots (plus flexible)
 Respecte ratio utilisateur (0.1  0.9)
 Phase X+1 appliquee si abstraction < 60%
 Enrichissement normal
 Humanisation normale (Phase X+3)
```

**Qualite obtenue:**
- Rythme preserve
- Tension narrative maintenue
- �motions transcrites
- Coherence acceptable

**Exemples:** Romans, contes, memoires, temoignages

---

### 3. **CONCEPTUEL** 
Philosophie, analyse, theorie, argumentation

**Caracteristiques detectees:**
- Concepts ("idee", "principe", "theorie")
- Logique ("cause", "consequence", "donc")
- Abstraction ("philosophie", "metaphysique")
- Argumentation ("par contre", "neanmoins", "cependant")

**Comportement du syst�me:**
```
 Resume par mots (compression agressive OK)
 Accepte ratio eleve (0.7-0.9)
 Phase X+1 for�ee si abstraction < 60%
 Humanisation Phase X+3 appliquee
 Convertit concret  abstrait
```

**Qualite obtenue:**
- Abstraction forcee
- Essence conservee
- Connecteurs naturels
- Argumention simplifiee

**Exemples:** Philosophie, essais, analyses politiques, theories

---

##  UTILISATION RECOMMAND�E

### Utilisation Simple (RECOMMAND�E)
```bash
./programme resume votre_document.txt 0.5
```
Le syst�me detecte le type et adapte automatiquement! 

**Resultat selon type:**
- Encyclopedique + ratio 0.1  30% garde automatiquement (protection)
- Narratif + ratio 0.5  50% garde (respecte)
- Conceptuel + ratio 0.7  30% garde (agressif, acceptable)

### Utilisation Avancee (PR�DICTION)

**Pour textes encyclopediques:**
```bash
# Vous voullez 15%  demandez 0.15  re�oit 30% (meilleur compromis)
./programme resume physique.txt 0.15

# Vous voullez 40%  demandez 0.4  re�oit 30-40%
./programme resume chimie.txt 0.4
```

**Pour textes narratifs:**
```bash
# Respecte votre ratio exactement
./programme resume roman.txt 0.3     # 30% garde
./programme resume histoire.txt 0.7  # 70% garde
```

**Pour textes conceptuels:**
```bash
# Permet compression agressive
./programme resume philo.txt 0.5      # Compression normale
./programme resume essai.txt 0.9      # Ultra-compression acceptable
```

---

##  INTERPR�TATION DU RAPPORT

### Exemple 1: Texte Encyclopedique

```
[PHASE 15] �tape 0: Detection du type de texte...
  �  Type: Encyclopedique (compression: 0.10)
    Phase X+1 desactivee pour ce type

[PHASE 15] �tape 2: Resume atomique (Phase 13+++)...
    Compression limitee � 30% pour texte encyclopedique (demandee: 10%)
   Resume genere: 1847 caract�res
```

**Interpretation:**
-  Type correct identifie
- � Phase X+1 skippee (fait correct, pas d'abstraction)
-  Vous avez demande 10% mais re�u 30%  Protection appliquee
-  Resume bien genere

**En bas du rapport:**
```
Text Type:         Encyclopedique
Abstraction:       SKIPPED (facts preservation)
Compression:       82.4%
```
 Tout est correct

---

### Exemple 2: Texte Narratif

```
[PHASE 15] �tape 0: Detection du type de texte...
  �  Type: Narratif (compression: 0.50)

[PHASE 15] �tape 2: Resume atomique...
   Resume genere: 12500 caract�res
```

**Interpretation:**
-  Type narratif detecte
- Pas de limite appliquee (ratio respecte)
- Abstraction appliquee si necessaire

---

##  MESSAGES SP�CIAUX

### " Phase X+1 desactivee pour ce type"
**Signifie:** Texte encyclopedique detecte  faits conserves  c'est normal et bon

### " Compression limitee � 30%"
**Signifie:** Vous avez demande < 30% sur encyclopedique  protection automatique appliquee

### "� Texte encyclopedique: conservation des faits concrets"
**Signifie:** Phase X+1 skippee  pas de for�age d'abstraction

### "� Texte encyclopedique: conservation structure originale"
**Signifie:** Phase X+3 skippee  pas de connecteurs parasites ajoutes

---

##  CAS D'USAGE PRATIQUES

### Cas 1: Resumer un article scientifique
```bash
./programme resume article_biologie.txt 0.2
```
- Detection: ENCYCLOP�DIQUE
- Limite: 30% garde (pas 20%)
- Resultat: Article coherent et factuel

### Cas 2: Resumer un roman
```bash
./programme resume roman_250pages.txt 0.1
```
- Detection: NARRATIF
- Limite: Aucune (10% respecte)
- Resultat: Resume court du rythme narratif

### Cas 3: Resumer un essai philosophique
```bash
./programme resume essai_politique.txt 0.5
```
- Detection: CONCEPTUEL
- Limite: Aucune (50% respecte)
- Resultat: Essence conceptuelle extraite

---

##  CONSEILS D'OPTIMISATION

### Pour Encyclopedique
- Demandez 30-40% pour meilleur resultat
- Ne demandez pas < 20% (perte importante)
- Le syst�me corrigera automatiquement

### Pour Narratif
- Demandez 30-50% pour bonne coherence
- Demandez 10-20% pour resume tr�s court
- Demandez 70-90% pour simple reduction

### Pour Conceptuel
- Demandez 50-70% pour abstraction moderee
- Demandez 80-90% pour ultra-dense
- Demandez 30-40% pour preserver plus de details

---

##  METRIQUES � SURVEILLER

### Compression
- < -10%: Expansion (texte enrichi, normal pour encyclopedique)
- 0-50%: Modere (bon compromis)
- 50-80%: Agressif (synopsis)
- > 80%: Ultra-dense (une phrase)

### Coherence Score
- > 50%: Excellent
- 20-50%: Bon
- < 20%: � ameliorer (trop agressif)

### Grammar Score
- > 75%: Excellent
- 50-75%: Bon
- < 50%: Problematique

---

##  D�PANNAGE

### "Le resume est trop court"
 Augmentez le ratio (ex: 0.3  0.5)
 Ou laissez le syst�me adapter pour encyclopedique

### "Le resume garde trop de contenu"
 Diminuez le ratio (ex: 0.7  0.5)
 Conceptuel accepte 0.9, encyclopedique a limite 0.3

### "Beaucoup de connecteurs inutiles"
 C'est Narratif/Conceptuel (normal)
 Encyclopedique a protection contre �a

### "Perte de faits importants"
 Verifiez que c'est detecte ENCYCLOP�DIQUE
 Augmentez ratio si < 0.3

---

##  CHECKLIST UTILISATION

- [ ] Vous avez un document � resumer
- [ ] Vous avez devine le type (encyclopedique/narratif/conceptuel)
- [ ] Vous avez choisi un ratio de compression
- [ ] Vous avez lance: `./programme resume document.txt ratio`
- [ ] Vous avez lu le rapport (cherchez les  et �)
- [ ] Vous avez compare avec l'original si doubtful
- [ ] Vous avez ajuste le ratio si necessaire et relance

---

##  CONCLUSION

Le syst�me est con�u pour vous donner:
1. **Automatisation**: Detection type sans effort
2. **Qualite**: Adaptation � chaque type
3. **Protection**: Limites intelligentes contre mauvaise compression
4. **Controle**: Vous fixez toujours le ratio
5. **Transparence**: Rapport detaille explique ce qui s'est passe

**Conseil:** Utilisez ratios 0.3-0.7 pour meilleur resultat. Laissez le syst�me adapter pour encyclopedique.

