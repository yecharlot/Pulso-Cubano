package prismatec

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/yecharlot/Pulso-Cubano/domain"
)

func TestAdapterObserverAndMind(t *testing.T) {
	dir := t.TempDir()
	ad, err := NewAdapter(filepath.Join(dir, "node"), "pulso-test")
	if err != nil {
		t.Fatal(err)
	}
	org, err := ad.EnsureObserverOrganism("GTMOObserver")
	if err != nil {
		t.Fatal(err)
	}
	price := 350.0
	obs := domain.Observation{
		ID: "obs-1", SourceID: "manual", SourceType: domain.SourceManualDataset,
		ObservedAt: time.Now().UTC(), RetrievedAt: time.Now().UTC(),
		Country: "CU", Province: "Guantánamo", Municipality: "Guantánamo",
		Category: "electronics", Product: "telefono",
		ObservationType: domain.ObsIntent, Price: &price, Currency: "USD",
		Text: "busco telefono Guantánamo", Fingerprint: "fp-test-1",
		Confidence: 0.6, Coverage: 0.1,
		Provenance: map[string]string{"layer": string(domain.LayerObserved)},
	}
	if err := ad.StoreObservationJSON(org.ID, ObservationMemoryKey(obs.Fingerprint), obs); err != nil {
		t.Fatal(err)
	}
	d, err := ad.AnalyzeObservation(context.Background(), org, obs)
	if err != nil {
		t.Fatal(err)
	}
	if d.Selected != "record_signal" {
		t.Fatalf("expected record_signal, got %+v", d)
	}
	ad.EmitMarketPulse(org.ID, "market.observation", "pulso-cubano", map[string]any{
		"fingerprint": obs.Fingerprint,
		"category":    obs.Category,
	})
	if ad.Node.Pulses().Count() < 1 {
		t.Fatal("expected pulse emitted")
	}
}
