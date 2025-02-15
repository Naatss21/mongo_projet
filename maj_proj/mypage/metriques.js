async function fetchMetrics() {
    try {
        const response = await fetch('/api/metrics');
        if (!response.ok) {
            throw new Error(`Erreur HTTP: ${response.status}`);
        }

        const data = await response.json();

        document.getElementById('cpu').innerText = `CPU Usage: ${data.cpu_usage.toFixed(2)}%`;
        document.getElementById('memory').innerText = `Memory Usage: ${data.mem_usage.toFixed(2)}%`;
        document.getElementById('disk').innerText = `Disk Usage: ${data.disk_usage.toFixed(2)}%`;
        document.getElementById('cores').innerText = `Cœurs physiques: ${data.core_number}`;
        document.getElementById('threads').innerText = `Threads (cœurs logiques): ${data.thread_number}`;

    } catch (error) {
        console.error("Erreur lors de la récupération des métriques :", error);
        document.getElementById('cpu').innerText = "Erreur de chargement...";
        document.getElementById('memory').innerText = "Erreur de chargement...";
        document.getElementById('disk').innerText = "Erreur de chargement...";
        document.getElementById('cores').innerText = "Erreur de chargement...";
        document.getElementById('threads').innerText = "Erreur de chargement...";
    }
}

// Lancer la récupération des données au chargement de la page
window.onload = () => {
    fetchMetrics();
    setInterval(fetchMetrics, 3000);
};
