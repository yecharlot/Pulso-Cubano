package genes

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/yecharlot/Pulso-Cubano/adapter/prismatec"
	"github.com/yecharlot/Pulso-Cubano/domain"
)

func TestManualSeedGeneMind(t *testing.T) {
	dir := t.TempDir()
	ad, err := prismatec.NewAdapter(filepath.Join(dir, "n"), "g")
	if err != nil {
		t.Fatal(err)
	}
	org, err := ad.EnsureObserverOrganism("GTMOObserver")
	if err != nil {
		t.Fatal(err)
	}
	price := 700.0
	g := ManualSeedGene{Items: []domain.Observation{{
		ID: "t1", SourceID: "manual", SourceType: domain.SourceManualDataset,
		ObservedAt: time.Now().UTC(), RetrievedAt: time.Now().UTC(),
		Country: "CU", Category: "fx", Product: "USD", ObservationType: domain.ObsPrice,
		Price: &price, Currency: "CUP", Text: "test", Fingerprint: "t1",
		Confidence: 0.5, Coverage: 0.1,
	}}}
	res, err := g.Run(context.Background(), ad, org.ID)
	if err != nil {
		t.Fatal(err)
	}
	if res.Observations != 1 {
		t.Fatalf("%+v", res)
	}
	if res.MindSelected != "record_signal" {
		t.Fatalf("mind=%q", res.MindSelected)
	}
}
