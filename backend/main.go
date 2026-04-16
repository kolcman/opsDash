package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Metrics struct {
	CPUPercent float64 `json:"cpuPercent"`
	RAMPercent float64 `json:"ramPercent"`
	UpdatedAt  string  `json:"updatedAt"`
}

type ServiceStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type AlertItem struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

var serviceNames = []string{
	"API Gateway",
	"Auth Service",
	"Billing",
	"Notifications",
	"Worker Queue",
}

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	promURL := envOrDefault("PROMETHEUS_URL", "")
	prom := newPromClient(promURL)
	mux := http.NewServeMux()

	mux.HandleFunc("/api/metrics", withCORS(func(w http.ResponseWriter, r *http.Request) {
		metrics, err := buildMetrics(rng, prom)
		if err != nil {
			log.Printf("metrics: falling back to demo values: %v", err)
			metrics = buildDemoMetrics(rng)
		}
		writeJSON(w, http.StatusOK, metrics)
	}))

	mux.HandleFunc("/api/services", withCORS(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, buildServices(rng))
	}))

	mux.HandleFunc("/api/alerts", withCORS(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, buildAlerts(rng))
	}))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	port := envOrDefault("PORT", "8080")
	log.Printf("OpsDash backend listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

func buildDemoMetrics(rng *rand.Rand) Metrics {
	return Metrics{
		CPUPercent: round1(randomRange(rng, 22, 86)),
		RAMPercent: round1(randomRange(rng, 35, 92)),
		UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
	}
}

func buildMetrics(rng *rand.Rand, prom *promClient) (Metrics, error) {
	if prom == nil || prom.baseURL == "" {
		return Metrics{}, errors.New("PROMETHEUS_URL is not configured")
	}

	// CPU usage % over last minute (across all CPUs).
	// Use 5m window to avoid empty result right after start; add vector(0) fallback.
	cpuQ := `100 * (1 - avg(rate(node_cpu_seconds_total{mode="idle"}[5m]))) or on() vector(0)`
	// RAM usage % (available vs total).
	ramQ := `100 * (1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) or on() vector(0)`

	cpu, err := prom.queryFloat(cpuQ)
	if err != nil {
		return Metrics{}, fmt.Errorf("prometheus cpu query: %w", err)
	}
	ram, err := prom.queryFloat(ramQ)
	if err != nil {
		return Metrics{}, fmt.Errorf("prometheus ram query: %w", err)
	}

	// keep within bounds, avoid NaN.
	if cpu < 0 {
		cpu = 0
	}
	if cpu > 100 {
		cpu = 100
	}
	if ram < 0 {
		ram = 0
	}
	if ram > 100 {
		ram = 100
	}

	return Metrics{
		CPUPercent: round1(cpu),
		RAMPercent: round1(ram),
		UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func buildServices(rng *rand.Rand) []ServiceStatus {
	out := make([]ServiceStatus, 0, len(serviceNames))
	for _, name := range serviceNames {
		status := weightedStatus(rng)
		out = append(out, ServiceStatus{
			Name:   name,
			Status: status,
		})
	}
	return out
}

func buildAlerts(rng *rand.Rand) []AlertItem {
	alertPool := []AlertItem{
		{
			ID:       "cpu-spike",
			Title:    "CPU usage spike",
			Severity: "warning",
			Message:  "CPU is above normal baseline for the last 5 minutes.",
		},
		{
			ID:       "disk-queue",
			Title:    "Disk queue pressure",
			Severity: "warning",
			Message:  "I/O wait increased on storage node 2.",
		},
		{
			ID:       "auth-timeout",
			Title:    "Auth latency critical",
			Severity: "critical",
			Message:  "p95 auth response time exceeded SLA threshold.",
		},
		{
			ID:       "billing-retry",
			Title:    "Billing retries elevated",
			Severity: "warning",
			Message:  "Payment callback retries are above expected range.",
		},
	}

	roll := rng.Intn(100)
	switch {
	case roll < 35:
		return []AlertItem{}
	case roll < 70:
		return []AlertItem{alertPool[rng.Intn(len(alertPool))]}
	default:
		first := alertPool[rng.Intn(len(alertPool))]
		second := alertPool[rng.Intn(len(alertPool))]
		if second.ID == first.ID {
			second = alertPool[(rng.Intn(len(alertPool)-1)+1)%len(alertPool)]
		}
		return []AlertItem{first, second}
	}
}

func weightedStatus(rng *rand.Rand) string {
	roll := rng.Intn(100)
	switch {
	case roll < 80:
		return "up"
	case roll < 94:
		return "degraded"
	default:
		return "down"
	}
}

func randomRange(rng *rand.Rand, min, max float64) float64 {
	return min + rng.Float64()*(max-min)
}

func round1(v float64) float64 {
	return float64(int(v*10)) / 10
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

type promClient struct {
	baseURL string
	client  *http.Client
}

func newPromClient(base string) *promClient {
	base = strings.TrimSpace(base)
	if base == "" {
		return &promClient{
			baseURL: "",
			client: &http.Client{
				Timeout: 4 * time.Second,
			},
		}
	}
	return &promClient{
		baseURL: strings.TrimRight(base, "/"),
		client: &http.Client{
			Timeout: 4 * time.Second,
		},
	}
}

type promQueryResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Value []any `json:"value"`
		} `json:"result"`
	} `json:"data"`
	ErrorType string `json:"errorType"`
	Error     string `json:"error"`
}

func (p *promClient) queryFloat(promQL string) (float64, error) {
	u, err := url.Parse(p.baseURL + "/api/v1/query")
	if err != nil {
		return 0, err
	}
	q := u.Query()
	q.Set("query", promQL)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return 0, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, fmt.Errorf("prometheus http %d", resp.StatusCode)
	}

	var decoded promQueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return 0, err
	}

	if decoded.Status != "success" {
		if decoded.Error != "" {
			return 0, fmt.Errorf("%s: %s", decoded.ErrorType, decoded.Error)
		}
		return 0, fmt.Errorf("prometheus status=%s", decoded.Status)
	}

	if len(decoded.Data.Result) == 0 || len(decoded.Data.Result[0].Value) < 2 {
		return 0, errors.New("no prometheus result")
	}

	raw, ok := decoded.Data.Result[0].Value[1].(string)
	if !ok {
		return 0, errors.New("unexpected prometheus value type")
	}

	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}
