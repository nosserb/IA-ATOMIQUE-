async function summarizeText() {
    const inputText = document.getElementById('inputText').value.trim();
    const outputBox = document.getElementById('outputText');
    const statsBox = document.getElementById('stats');
    const btn = document.getElementById('summarizeBtn');

    if (!inputText) {
        outputBox.innerHTML = '<p class="placeholder" style="color: #d32f2f;">Veuillez entrer du texte!</p>';
        return;
    }

    // Afficher le loader
    outputBox.classList.add('loading');
    outputBox.innerHTML = '<p class="placeholder">Résumé en cours...</p>';
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

        // Afficher le résumé
        outputBox.classList.remove('loading');
        outputBox.innerHTML = data.summary || 'Aucun résumé généré';

        // Afficher les statistiques
        if (data.stats) {
            statsBox.innerHTML = `
                <strong>Statistiques:</strong><br>
                • Phrases analysées: ${data.stats.phrases || 0}<br>
                • Confiance moyenne: ${data.stats.confidence || '0'}%<br>
                • Catégories détectées: ${data.stats.categories || 0}
            `;
        }

    } catch (error) {
        outputBox.classList.remove('loading');
        outputBox.innerHTML = '<p class="placeholder" style="color: #d32f2f;">Erreur: ' + error.message + '</p>';
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
