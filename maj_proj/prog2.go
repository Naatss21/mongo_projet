package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

// Structure pour les métriques
type Metrics struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemUsage    float64 `json:"mem_usage"`
	DiskUsage   float64 `json:"disk_usage"`
	CoreNumber  int     `json:"core_number"`
	ThreadNumber int    `json:"thread_number"`
}

// Récupération des métriques système
func getSysMetrics() Metrics {
	cpuUsage, _ := cpu.Percent(time.Second, false)
	memStats, _ := mem.VirtualMemory()
	diskStats, _ := disk.Usage("/")
	logicalCores, _ := cpu.Counts(true)
	physicalCores, _ := cpu.Counts(false)

	return Metrics{
		CPUUsage:    cpuUsage[0],
		MemUsage:    memStats.UsedPercent,
		DiskUsage:   diskStats.UsedPercent,
		CoreNumber:  physicalCores,
		ThreadNumber: logicalCores,
	}
}

// Stockage des métriques dans InfluxDB
func storeMetricsInInflux(metrics Metrics) {
	influxURL := "http://localhost:8086"
	influxOrg := "mongo0_1"
	influxBucket := "metrics0_1"
	token := os.Getenv("INFLUX_TOKEN")

	client := influxdb2.NewClient(influxURL, token)
	defer client.Close()

	writeAPI := client.WriteAPIBlocking(influxOrg, influxBucket)
	p := influxdb2.NewPointWithMeasurement("system_metrics").
		AddField("cpu_usage", metrics.CPUUsage).
		AddField("mem_usage", metrics.MemUsage).
		AddField("disk_usage", metrics.DiskUsage).
		AddField("core_number", metrics.CoreNumber).
		AddField("thread_number", metrics.ThreadNumber).
		SetTime(time.Now())

	err := writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		log.Println("❌ Erreur lors de l'écriture dans InfluxDB:", err)
	} else {
		log.Println("✅ Métriques stockées dans InfluxDB")
	}
}

// Handler pour l'API
func metricsHandler(w http.ResponseWriter, r *http.Request) {
	metrics := getSysMetrics()
	storeMetricsInInflux(metrics)  // Stocke les métriques dans InfluxDB
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// Fonction principale
func main() {
	fmt.Println("🌍 Serveur démarré sur http://localhost:8080")
	http.HandleFunc("/api/metrics", metricsHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Erreur lors du démarrage du serveur :", err)
	}
}
