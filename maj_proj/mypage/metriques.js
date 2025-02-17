let historique = []; // Tableau pour stocker l'historique des métriques

async function fetchMetrics() {
    try {
        const response = await fetch('/api/metrics');
        if (!response.ok) {
            throw new Error(`Erreur HTTP: ${response.status}`);
        }

        const data = await response.json();

        // Affichage des valeurs actuelles
        document.getElementById('cpu').innerText = `CPU Usage: ${data.cpu_usage.toFixed(2)}%`;
        document.getElementById('memory').innerText = `Memory Usage: ${data.mem_usage.toFixed(2)}%`;
        document.getElementById('disk').innerText = `Disk Usage: ${data.disk_usage.toFixed(2)}%`;
        document.getElementById('cores').innerText = `Cœurs physiques: ${data.core_number}`;
        document.getElementById('threads').innerText = `Threads (cœurs logiques): ${data.thread_number}`;

        // Ajouter la métrique actuelle à l'historique
        const currentMetrics = {
            cpu: data.cpu_usage.toFixed(2),
            memory: data.mem_usage.toFixed(2),
            disk: data.disk_usage.toFixed(2),
            cores: data.core_number,
            threads: data.thread_number
        };

        historique.push(currentMetrics);
        updateHistorique();

		} 
			catch (error) {
				console.error("Erreur lors de la récupération des métriques :", error);
				document.getElementById('cpu').innerText = "Erreur de chargement...";
				document.getElementById('memory').innerText = "Erreur de chargement...";
				document.getElementById('disk').innerText = "Erreur de chargement...";
				document.getElementById('cores').innerText = "Erreur de chargement...";
				document.getElementById('threads').innerText = "Erreur de chargement...";
    }
}

// Fonction pour afficher l'historique
function updateHistorique() {
    const historiqueList = document.getElementById('historique-list');
    historiqueList.innerHTML = ''; // Vider l'historique avant de le mettre à jour

    // Limiter l'historique à 10 éléments pour éviter trop de données
    const historiqueLimite = historique.slice(-10);

    historiqueLimite.forEach(metrics => {
        const li = document.createElement('li');
        li.innerHTML = `
            CPU: ${metrics.cpu}% | 
            Memory: ${metrics.memory}% | 
            Disk: ${metrics.disk}% |
        `;
        historiqueList.appendChild(li);
    });
}

// Lancer la récupération des données au chargement de la page
window.onload = () => {
    fetchMetrics();
    setInterval(fetchMetrics, 3000); // Mettre à jour toutes les 3 secondes
};
