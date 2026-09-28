// Package prismatec is the ONLY bridge from Pulso Cubano to PrismaTec-Core.
// Core does not import this package. Vertical owns all market ontology.
package prismatec

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/yecharlot/PrismaTec-Core/core"
	"github.com/yecharlot/PrismaTec-Core/core/organism"
	"github.com/yecharlot/PrismaTec-Core/core/policy"
	"github.com/yecharlot/PrismaTec-Core/core/pulse"
	"github.com/yecharlot/PrismaTec-Core/runtime/inference"
	"github.com/yecharlot/PrismaTec-Core/runtime/mind"
	"github.com/yecharlot/PrismaTec-Core/runtime/zyrion"

	"github.com/yecharlot/Pulso-Cubano/domain"
)

// Adapter wraps a Core node for market-intelligence workloads.
type Adapter struct {
	Node *core.Node
}

func NewAdapter(dataDir, name string) (*Adapter, error) {
	n, err := core.NewNode(core.Config{DataDir: dataDir, Name: name})
	if err != nil {
		return nil, err
	}
	return &Adapter{Node: n}, nil
}

func (a *Adapter) NodeID() string {
	if a.Node == nil || a.Node.Identity() == nil {
		return ""
	}
	return string(a.Node.ID())
}

// EnsureObserverOrganism creates a vertical observer organism if missing.
func (a *Adapter) EnsureObserverOrganism(name string) (*organism.Organism, error) {
	mgr := a.Node.Organisms()
	for _, o := range mgr.List() {
		if o.Name == name {
			return o, nil
		}
	}
	return mgr.Create(organism.CreateOptions{
		Name: name,
		Capabilities: []organism.Capability{
			"memory.read", "memory.write", "inference",
		},
	})
}

// StoreObservationJSON persists observation JSON in organism working memory.
func (a *Adapter) StoreObservationJSON(orgID, key string, obs domain.Observation) error {
	b, err := json.Marshal(obs)
	if err != nil {
		return err
	}
	_, err = a.Node.Organisms().PutMemory(orgID, key, string(b))
	return err
}

// EmitMarketPulse publishes a Core pulse for Studio/AIP observers.
func (a *Adapter) EmitMarketPulse(target, typ, source string, data map[string]any) {
	if a.Node.Pulses() == nil {
		return
	}
	a.Node.Pulses().Emit(pulse.Pulse{
		Target:    target,
		Type:      typ,
		Data:      data,
		Source:    source,
		Timestamp: time.Now().UTC(),
	})
}

func (a *Adapter) PolicyEngine(rules ...policy.Rule) *policy.Engine {
	return policy.NewEngine(rules...)
}

// AnalyzeObservation uses Mind to propose record_signal; never authorizes by itself.
// Beliefs use Core's evidence.for ingest path with proposition market.intent|market.listing.
func (a *Adapter) AnalyzeObservation(ctx context.Context, org *organism.Organism, obs domain.Observation) (*mind.Decision, error) {
	rules := zyrion.RuleSet{Version: "pulso-v0", Rules: []zyrion.Rule{
		{ID: "signal-intent", Version: "pulso-v0", When: "market.intent", Then: "record_signal"},
		{ID: "signal-listing", Version: "pulso-v0", When: "market.listing", Then: "record_signal"},
	}}
	m := mind.New(mind.Config{
		OrganismID: org.ID,
		Rules:      rules,
		Provider:   inference.EchoProvider{},
	})
	prop := "market.listing"
	if obs.ObservationType == domain.ObsIntent {
		prop = "market.intent"
	}
	_ = m.Observe(mind.Observation{
		ID:   obs.ID,
		Type: "evidence.for",
		Payload: map[string]any{
			"proposition": prop,
			"category":    obs.Category,
			"product":     obs.Product,
			"geo":         obs.Municipality,
			"layer":       string(domain.LayerObserved),
		},
		At: time.Now().UTC(),
	})
	ev, err := m.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	return &ev.Decision, nil
}

func ObservationMemoryKey(fp string) string {
	return fmt.Sprintf("obs:%s", fp)
}
