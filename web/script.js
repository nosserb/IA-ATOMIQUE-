async function summarizeText() {
    const inputText = document.getElementById('inputText').value.trim();
    const keyPhrasesBox = document.getElementById('keyPhrases');
    const keyIdeasBox = document.getElementById('keyIdeas');
    const summaryBox = document.getElementById('summary');
    const statsBox = document.getElementById('stats');
    const btn = document.getElementById('summarizeBtn');

    if (!inputText) {
        keyPhrasesBox.innerHTML = '<p class="placeholder" style="color: #d32f2f;">Veuillez entrer du texte!</p>';
        return;
    }

    // Afficher le loader
    keyPhrasesBox.innerHTML = '<p class="placeholder">Analyse en cours...</p>';
    keyIdeasBox.innerHTML = '<p class="placeholder">Analyse en cours...</p>';
    summaryBox.innerHTML = '<p class="placeholder">Analyse en cours...</p>';
    statsBox.innerHTML = '';
    btn.disabled = true;

    try {
        const response = await fetch('/api/summarize', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                text: inputText
            })
        });

        if (!response.ok) {
            throw new Error('Erreur serveur: ' + response.statusText);
        }

        const data = await response.json();

        // Afficher les phrases clés
        if (data.keyPhrases && data.keyPhrases.length > 0) {
            keyPhrasesBox.innerHTML = data.keyPhrases
                .map(phrase => `<div style="margin-bottom: 10px;">• ${phrase}</div>`)
                .join('');
        } else {
            keyPhrasesBox.innerHTML = '<p class="placeholder">Aucune phrase clé détectée</p>';
        }

        // Afficher les idées clés
        if (data.keyIdeas && data.keyIdeas.length > 0) {
            keyIdeasBox.innerHTML = data.keyIdeas
                .map(idea => `<div style="margin-bottom: 10px;">• ${idea}</div>`)
                .join('');
        } else {
            keyIdeasBox.innerHTML = '<p class="placeholder">Aucune idée clée détectée</p>';
        }

        // Afficher le résumé
        if (data.summary) {
            summaryBox.innerHTML = `<p>${data.summary}</p>`;
        } else {
            summaryBox.innerHTML = '<p class="placeholder">Aucun résumé généré</p>';
        }

        // Afficher les statistiques
        if (data.stats) {
            statsBox.innerHTML = `
                <strong>📊 Statistiques:</strong><br>
                • Phrases analysées: ${data.stats.phrases || 0}<br>
                • Confiance moyenne: ${data.stats.confidence || '0'}%<br>
                • Catégories détectées: ${data.stats.categories || 0}
            `;
        }

    } catch (error) {
        keyPhrasesBox.innerHTML = '<p class="placeholder" style="color: #d32f2f;">Erreur: ' + error.message + '</p>';
        keyIdeasBox.innerHTML = '<p class="placeholder" style="color: #d32f2f;">Erreur</p>';
        summaryBox.innerHTML = '<p class="placeholder" style="color: #d32f2f;">Erreur</p>';
    } finally {
        btn.disabled = false;
    }
}

// Permettre Ctrl+Entrée pour soumettre
document.addEventListener('keydown', function(event) {
    if (event.ctrlKey && event.key === 'Enter') {
        summarizeText();
    }
});
