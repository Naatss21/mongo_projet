package main


import (
		"fmt"
		"log"
		"time"
		"os"
		"github.com/influxdata/influxdb-client-go/v2/api/write"
		
		"encoding/json"
		"net/http"
		
		"github.com/influxdata/influxdb-client-go/v2"
		"github.com/shirou/gopsutil/v3/cpu"
		"github.com/shirou/gopsutil/v3/mem"
		"github.com/shirou/gopsutil/v3/disk"
		"context"
		//"github.com/go-sql-driver/mysql"
)


// Structure pour les métriques
type Metrics struct {
	CPUUsage float64 `json:"cpu_usage"`
	MemUsage float64 `json:"mem_usage"`
	DiskUsage float64 `json:"disk_usage"`
	CoreNumber int `json:"core_number"`
	ThreadNumber int `json:"thread_number"` 
}

func getSysMetrics () Metrics {
	var cpuUsage, memUsage, diskUsage float64
	var corenumber, threadnumber int

	// Récupérer l'utilisation du CPU
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Println("Erreur lors de la récupération de l'utilisation du CPU:", err)
	} else {
		cpuUsage = percentages[0]
	}

	// Récupérer l'utilisation de la mémoire
	memStats, err := mem.VirtualMemory()
	if err != nil {
		log.Println("Erreur lors de la récupération de la mémoire:", err)
	} else {
		memUsage = memStats.UsedPercent
	}

	// Récupérer l'utilisation du disque
	diskStats, err := disk.Usage("/")
	if err != nil {
		log.Println("Erreur lors de la récupération de l'utilisation du disque:", err)
	} else {
		diskUsage = diskStats.UsedPercent
	}
	// Récupérer le nombre de coeur(s) et thread(s)
	logicalCores, err := cpu.Counts(true)
	if err != nil {
		log.Fatal("Erreur lors de la récupération des cœurs logiques:", err)
	} else {
	fmt.Println("Nombre de cœurs logiques :", logicalCores)
	threadnumber = logicalCores 
		
	}
	
	physicalCores, err := cpu.Counts(false)
	if err != nil {
		log.Fatal("Erreur lors de la récupération des cœurs physiques:", err)
	} else {
		
	fmt.Println("Nombre de cœurs physiques :", physicalCores)
	corenumber = physicalCores	
	}
	
	
	
	// Retourner les métriques
	return Metrics{
		CPUUsage:  cpuUsage,
		MemUsage:  memUsage,
		DiskUsage: diskUsage,
		CoreNumber: corenumber,
		ThreadNumber: threadnumber,
	}
	
}

/*
func storeMetricsInInflux(metrics Metrics) {
	// 🚀 Configuration InfluxDB
	
	const (
		influxURL   = "http://localhost:8086"
		influxOrg   = "mongo0_1"
		influxBucket = "metrics0_1"
	)
	token := os.Getenv("INFLUX_TOKEN")
	
	// Créer un client InfluxDB
	client := influxdb2.NewClient(influxURL, token)
	defer client.Close()

	// Ecriture des données
	writeAPI := client.WriteAPIBlocking(influxOrg, influxBucket)
	p := influxdb2.NewPointWithMeasurement("system_metrics").
		AddField("cpu_usage", metrics.CPUUsage).
		AddField("mem_usage", metrics.MemUsage).
		AddField("disk_usage", metrics.DiskUsage).
		AddField("core_number", metrics.CoreNumber).
		AddField("thread_number", metrics.ThreadNumber).
		SetTime(time.Now())
	
	// 🚀 Envoi des données à InfluxDB
	err := writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		log.Println("❌ Erreur lors de l'écriture dans InfluxDB:", err)
	} else {
		log.Println("✅ Métriques stockées dans InfluxDB")
	}
}
	
*/

// Handler pour l'API
func metricsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(" Requête reçue sur /api/metrics") // Ajout de logs
	// Récupérer les données système
	metrics := getSysMetrics()
	// Définir l'en-tête Content-Type pour JSON
	// Stocker les métriques dans InfluxDB
	//storeMetricsInInflux(metrics)  // <-- Ajout de cet appel
	w.Header().Set("Content-Type", "application/json")
	// Encoder les données en JSON et les envoyer dans la réponse
	err := json.NewEncoder(w).Encode(metrics)
	if err != nil {
		http.Error(w, "Erreur lors de l'encodage JSON", http.StatusInternalServerError)
	}
}


func main () {
	token := os.Getenv("INFLUX_TOKEN")
	url := "http://localhost:8086"
	client := influxdb2.NewClient(url, token)
	
	fmt.Println("Version TEST!")
	fmt.Println("retrouve tes metriques sur : http://localhost:8080/api/metrics")
	fmt.Println("🌍 Serveur démarré sur http://localhost:8080")

		
	// Route pour l'API qui renvoie les métriques
	http.HandleFunc("/api/metrics", metricsHandler)
	// Route pour servir la page web (index.html)
	http.Handle("/", http.FileServer(http.Dir("./mypage")))

	// Démarrer le serveur HTTP avec gestion d'erreur
	err := http.ListenAndServe(":8080", nil)
	
	if err != nil {
		log.Fatal("Erreur lors du démarrage du serveur :", err)
			
		
	}
	
	//modif
	org := "mongo0_1"
	bucket := "metrics0_1"
	writeAPI := client.WriteAPIBlocking(org, bucket)
	for value := 0; value < 5; value++ {
		tags := map[string]string{
			"tagname1": "tagvalue1",}
		fields := map[string]interface{}{
		"field1": value,
		}
		point := write.NewPoint("measurement1", tags, fields, time.Now())
		time.Sleep(1 * time.Second) // separate points by 1 second
		if err := writeAPI.WritePoint(context.Background(), point); err != nil {
			log.Fatal(err)
		}
}
	
	
	

}
