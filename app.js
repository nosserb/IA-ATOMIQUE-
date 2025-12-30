// Configuration API
const API_BASE = 'http://localhost:8080/api';

// État global
const state = {
    isTraining: false,
    stats: {
        activeNeurons: 0,
        totalEnergy: 0,
        categoryCount: 0,
        sentenceCount: 0
    }
};

// Analyse de texte
async function analyzeText() {
    const text = document.getElementById('inputText').value.trim();
    
    if (!text) {
        showError('Enter text to analyze');
        return;
    }

    const loading = document.getElementById('analysisLoading');
    const output = document.getElementById('analysisOutput');
    const outputText = document.getElementById('analysisText');

    loading.style.display = 'block';
    output.style.display = 'none';

    try {
        const response = await fetch(`${API_BASE}/analyze`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ text: text })
        });

        if (!response.ok) throw new Error('Server error');

        const data = await response.json();
        
        outputText.innerHTML = formatAnalysisResult(data);
        output.style.display = 'block';
        showSuccess('Analysis complete');
        
    } catch (error) {
        console.error('Error:', error);
        showError('Analysis failed: ' + error.message);
    } finally {
        loading.style.display = 'none';
    }
}

// Format résultat d'analyse
function formatAnalysisResult(data) {
    if (!data) return '<span class="text-gray-500">No results</span>';
    
    let html = '';
    
    if (data.categories && data.categories.length > 0) {
        html += '<div class="mb-6"><span class="text-cyan-400">&gt;</span> CATEGORIES DETECTED<br>';
        data.categories.forEach(cat => {
            const percent = (cat.confidence * 100).toFixed(1);
            html += `├─ ${cat.name} <span class="text-gray-500">[${percent}%]</span><br>`;
        });
        html += '</div>';
    }
    
    if (data.entities && data.entities.length > 0) {
        html += '<div class="mb-6"><span class="text-cyan-400">&gt;</span> KEY ELEMENTS<br>';
        data.entities.forEach(ent => {
            html += `├─ ${ent.text} <span class="text-gray-500">[${ent.type}]</span><br>`;
        });
        html += '</div>';
    }
    
    // SUMMARY - Always show it with better styling
    if (data.summary) {
        html += `<div class="mt-6 p-4 border-l-2 border-cyan-400 bg-cyan-400/5">
                    <span class="text-cyan-400 font-bold">&gt;&gt; ANALYSIS REPORT</span><br>
                    <span class="text-gray-300">${data.summary}</span>
                </div>`;
    } else {
        html += `<div class="mt-6"><span class="text-cyan-400">&gt;</span> ANALYSIS COMPLETE<br><span class="text-gray-400">└─ Text processed successfully</span></div>`;
    }
    
    return html || '<span class="text-gray-500">No analysis data</span>';
}

// Rafraîchir les stats
async function refreshStats() {
    try {
        const response = await fetch(`${API_BASE}/stats`);
        
        if (!response.ok) {
            throw new Error('Stats unavailable');
        }

        const stats = await response.json();
        showStats(stats);
        updateProgressBars(stats);
        showSuccess('Metrics refreshed');
        
    } catch (error) {
        console.error('Error:', error);
        showError('Failed to load metrics');
    }
}

// Afficher les stats
function showStats(stats) {
    document.getElementById('activeNeurons').textContent = 
        (stats.activeNeurons || 0).toString();
    document.getElementById('totalEnergy').textContent = 
        ((stats.totalEnergy || 0).toFixed(1));
    document.getElementById('categoryCount').textContent = 
        (stats.categoryCount || 0).toString();
    document.getElementById('sentenceCount').textContent = 
        (stats.sentenceCount || 0).toString();
    
    state.stats = stats;
}

// Mettre à jour les barres de progression
function updateProgressBars(stats) {
    const maxNeurons = 1000;
    const activity = Math.min((stats.activeNeurons / maxNeurons) * 100, 100);
    const load = Math.min((activity / 2) + Math.random() * 30, 100);
    const memory = Math.min(activity * 0.8, 100);
    
    document.getElementById('activityBar').style.width = activity + '%';
    document.getElementById('activityPercent').textContent = Math.round(activity) + '%';
    
    document.getElementById('loadBar').style.width = load + '%';
    document.getElementById('loadPercent').textContent = Math.round(load) + '%';
    
    document.getElementById('memBar').style.width = memory + '%';
    document.getElementById('memPercent').textContent = Math.round(memory) + '%';
}

// Entraînement
async function startTraining() {
    const file = document.getElementById('trainingFile').files[0];
    const epochs = parseInt(document.getElementById('epochs').value);

    if (!file) {
        showError('Veuillez sélectionner un fichier');
        return;
    }

    state.isTraining = true;
    const statusDiv = document.getElementById('trainingStatus');
    statusDiv.innerHTML = '<div class="loading"><div class="spinner"></div><p>Entraînement en cours...</p></div>';

    try {
        const formData = new FormData();
        formData.append('file', file);
        formData.append('epochs', epochs);

        const response = await fetch(`${API_BASE}/train`, {
            method: 'POST',
            body: formData
        });

        if (!response.ok) throw new Error('Erreur serveur');

        const result = await response.json();
        statusDiv.innerHTML = `
            <div class="success">
                <strong>Entraînement Complété!</strong><br>
                Précision: ${(result.accuracy * 100).toFixed(2)}%<br>
                Temps: ${result.duration}s
            </div>
        `;
        showSuccess('Entraînement terminé');
        
    } catch (error) {
        console.error('Erreur:', error);
        statusDiv.innerHTML = `<div class="error">Erreur: ${error.message}</div>`;
    } finally {
        state.isTraining = false;
    }
}

// Paramètres
function saveSettings() {
    const settings = {
        learningRate: parseFloat(document.getElementById('learningRate').value),
        threshold: parseFloat(document.getElementById('threshold').value)
    };

    localStorage.setItem('iaSettings', JSON.stringify(settings));
    showSuccess('Paramètres enregistrés');
}

// Onglets
function switchTab(tabName) {
    // Masquer tous les onglets
    document.querySelectorAll('.tab-content').forEach(tab => {
        tab.classList.remove('active');
    });
    document.querySelectorAll('.tab-button').forEach(btn => {
        btn.classList.remove('active');
    });

    // Afficher l'onglet sélectionné
    document.getElementById(tabName).classList.add('active');
    event.target.classList.add('active');
}

// Utilitaires
function clearInput() {
    document.getElementById('inputText').value = '';
    document.getElementById('analysisOutput').style.display = 'none';
}

function showError(message) {
    showMessage(message, 'error');
}

function showSuccess(message) {
    showMessage(message, 'success');
}

function showMessage(message, type) {
    // Créer une notification temporaire
    const notification = document.createElement('div');
    notification.className = 'notification';
    if (type === 'error') {
        notification.style.borderColor = 'rgba(239, 68, 68, 0.3)';
        notification.style.color = '#ef4444';
    } else {
        notification.style.borderColor = 'rgba(0, 243, 255, 0.3)';
        notification.style.color = '#00f3ff';
    }
    notification.textContent = message;
    document.body.appendChild(notification);

    setTimeout(() => {
        notification.style.animation = 'slideOut 0.3s ease';
        setTimeout(() => notification.remove(), 300);
    }, 3000);
}

// Charger les paramètres au démarrage
window.addEventListener('DOMContentLoaded', () => {
    const saved = localStorage.getItem('iaSettings');
    if (saved) {
        const settings = JSON.parse(saved);
        document.getElementById('learningRate').value = settings.learningRate;
        document.getElementById('threshold').value = settings.threshold;
    }
    
    refreshStats();
    
    // Animer les barres au démarrage
    setTimeout(() => {
        updateProgressBars({
            activeNeurons: Math.floor(Math.random() * 500) + 100,
            totalEnergy: 2500,
            categoryCount: 10,
            sentenceCount: 25
        });
    }, 500);
});
