package api

import (
	"context"
	"encoding/json"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/yecharlot/Pulso-Cubano/adapter/prismatec"
	"github.com/yecharlot/Pulso-Cubano/domain"
	"github.com/yecharlot/Pulso-Cubano/genes"
	"github.com/yecharlot/Pulso-Cubano/metrics"
	"github.com/yecharlot/Pulso-Cubano/queries"
	"github.com/yecharlot/Pulso-Cubano/sources"
)

// Server exposes Pulso JSON for the Alset-JS control room.
type Server struct {
	Adapter *prismatec.Adapter
	OrgID   string
	UIDir   string
	mu      sync.RWMutex
	Last    []genes.GeneResult
}

func NewServer(dataDir, uiDir string) (*Server, error) {
	ad, err := prismatec.NewAdapter(filepath.Join(dataDir, "node"), "pulso-cubano")
	if err != nil {
		return nil, err
	}
	org, err := ad.EnsureObserverOrganism("GTMOObserver")
	if err != nil {
		return nil, err
	}
	return &Server{Adapter: ad, OrgID: org.ID, UIDir: uiDir}, nil
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"ok": true, "product": "pulso-cubano", "ui": "alset-js-runtime"})
	})
	mux.HandleFunc("/api/v1/sources", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, sources.SeedRegistry().List())
	})
	mux.HandleFunc("/api/v1/queries/gtmo", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, queries.ExpandAll(queries.GTMOSeedTemplates(), nil))
	})
	mux.HandleFunc("/api/v1/pulse", s.handlePulse)
	mux.HandleFunc("/api/v1/genes/run", s.handleGenesRun)
	mux.HandleFunc("/api/v1/terminal", s.handleTerminal)
	mux.HandleFunc("/api/v1/dashboard", s.handleDashboard)
	mux.HandleFunc("/api/v1/genes/last", func(w http.ResponseWriter, r *http.Request) {
		s.mu.RLock()
		defer s.mu.RUnlock()
		writeJSON(w, 200, s.Last)
	})
	if s.UIDir != "" {
		mux.Handle("/", http.FileServer(http.Dir(s.UIDir)))
	}
	return mux
}

func (s *Server) handlePulse(w http.ResponseWriter, r *http.Request) {
	org, err := s.Adapter.Node.Organisms().Get(s.OrgID)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	raw, ok := org.Memory.Working["market.pulse"]
	if !ok || raw == "" {
		writeJSON(w, 200, metrics.ComputePulse(0, 0, 0, 0, 0, 0, 0, 0))
		return
	}
	var p metrics.PulseScore
	if json.Unmarshal([]byte(raw), &p) != nil {
		writeJSON(w, 200, map[string]string{"raw": raw})
		return
	}
	writeJSON(w, 200, p)
}

func (s *Server) handleGenesRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		w.WriteHeader(405)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
	defer cancel()
	// Always include a small observed-market demo seed so UI works offline
	now := time.Now().UTC()
	price := 720.0
	seed := genes.ManualSeedGene{Items: []domain.Observation{{
		ID: "seed-usd", SourceID: "manual-csv", SourceType: domain.SourceManualDataset,
		ObservedAt: now, RetrievedAt: now, Country: "CU", Province: "Guantánamo",
		Municipality: "Guantánamo", Category: "fx", Product: "USD",
		ObservationType: domain.ObsPrice, Price: &price, Currency: "CUP",
		Text: "seed referential USD/CUP for UI demo", Fingerprint: "seed-usd-hour",
		Confidence: 0.3, Coverage: 0.05,
		Provenance: map[string]string{"layer": "observed_market", "demo": "true"},
	}}}
	res, err := genes.RunAll(ctx, s.Adapter, s.OrgID, seed)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	s.mu.Lock()
	s.Last = res
	s.mu.Unlock()
	writeJSON(w, 200, res)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
