// Package genes implements Pulso market genes: collect → normalize → Mind propose.
// Genes run in the vertical; Core only provides organism memory, pulse, policy, Mind.
package genes

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/yecharlot/Pulso-Cubano/adapter/prismatec"
	"github.com/yecharlot/Pulso-Cubano/domain"
	"github.com/yecharlot/Pulso-Cubano/ingestion"
	"github.com/yecharlot/Pulso-Cubano/metrics"
)

// Gene is a named capability of a market observer organism.
type Gene interface {
	Name() string
	Run(ctx context.Context, ad *prismatec.Adapter, orgID string) (GeneResult, error)
}

type GeneResult struct {
	Name          string    `json:"name"`
	Observations  int       `json:"observations"`
	MindSelected  string    `json:"mind_selected"`
	PulseOverall  float64   `json:"pulse_overall"`
	Errors        []string  `json:"errors,omitempty"`
	At            time.Time `json:"at"`
}

// FXGene collects informal FX (El Toque) when token present.
type FXGene struct{}

func (FXGene) Name() string { return "gene.fx.eltoque" }

func (g FXGene) Run(ctx context.Context, ad *prismatec.Adapter, orgID string) (GeneResult, error) {
	res := GeneResult{Name: g.Name(), At: time.Now().UTC()}
	client := ingestion.NewElToqueFromEnv()
	if !client.Enabled() {
		res.Errors = append(res.Errors, "eltoque disabled (set ELTOQUE_API_TOKEN)")
		return res, nil
	}
	obs, err := client.FetchTRMI(ctx)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		return res, nil
	}
	return ingestAndMind(ctx, ad, orgID, g.Name(), obs, res)
}

// QvaPayGene collects P2P averages when token present.
type QvaPayGene struct{}

func (QvaPayGene) Name() string { return "gene.fx.qvapay" }

func (g QvaPayGene) Run(ctx context.Context, ad *prismatec.Adapter, orgID string) (GeneResult, error) {
	res := GeneResult{Name: g.Name(), At: time.Now().UTC()}
	client := ingestion.NewQvaPayFromEnv()
	if !client.Enabled() {
		res.Errors = append(res.Errors, "qvapay disabled (set QVAPAY_TOKEN)")
		return res, nil
	}
	obs, err := client.FetchP2PAverages(ctx)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		return res, nil
	}
	return ingestAndMind(ctx, ad, orgID, g.Name(), obs, res)
}

// ManualSeedGene injects synthetic demo observations for offline UI (never as "official").
type ManualSeedGene struct {
	Items []domain.Observation
}

func (ManualSeedGene) Name() string { return "gene.manual.seed" }

func (g ManualSeedGene) Run(ctx context.Context, ad *prismatec.Adapter, orgID string) (GeneResult, error) {
	res := GeneResult{Name: g.Name(), At: time.Now().UTC()}
	return ingestAndMind(ctx, ad, orgID, g.Name(), g.Items, res)
}

func ingestAndMind(ctx context.Context, ad *prismatec.Adapter, orgID, gene string, obs []domain.Observation, res GeneResult) (GeneResult, error) {
	org, err := ad.Node.Organisms().Get(orgID)
	if err != nil {
		return res, err
	}
	for i := range obs {
		if obs[i].Fingerprint == "" {
			obs[i].Fingerprint = hashFP(obs[i].Text + obs[i].Product + obs[i].SourceID)
		}
		key := prismatec.ObservationMemoryKey(obs[i].Fingerprint)
		_ = ad.StoreObservationJSON(orgID, key, obs[i])
		d, err := ad.AnalyzeObservation(ctx, org, obs[i])
		if err == nil && d != nil {
			res.MindSelected = d.Selected
		}
		ad.EmitMarketPulse(orgID, "market.observation", gene, map[string]any{
			"product": obs[i].Product, "category": obs[i].Category, "source": obs[i].SourceID,
		})
	}
	res.Observations = len(obs)
	// Pulse v0 from simple counts (observed only)
	p := metrics.ComputePulse(60, 50, 55, 70, 40, 50, min(100, float64(len(obs)*10)), len(obs))
	res.PulseOverall = p.Overall
	b, _ := json.Marshal(p)
	_, _ = ad.Node.Organisms().PutMemory(orgID, "market.pulse", string(b))
	return res, nil
}

func hashFP(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:8])
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// RunAll executes all genes in order (FX → QvaPay → optional manual).
func RunAll(ctx context.Context, ad *prismatec.Adapter, orgID string, extra ...Gene) ([]GeneResult, error) {
	list := []Gene{FXGene{}, QvaPayGene{}}
	list = append(list, extra...)
	var out []GeneResult
	for _, g := range list {
		r, err := g.Run(ctx, ad, orgID)
		if err != nil {
			return out, fmt.Errorf("%s: %w", g.Name(), err)
		}
		out = append(out, r)
	}
	return out, nil
}
