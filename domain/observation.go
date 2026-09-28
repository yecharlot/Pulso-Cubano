package domain

import "time"

// Observation is the atomic market fact: what was seen, where, when, from which source.
// It is NOT a sale, NOT market share, and NOT invented demand.
type Observation struct {
	ID              string            `json:"id"`
	SourceID        string            `json:"source_id"`
	SourceType      SourceType        `json:"source_type"`
	URL             string            `json:"url,omitempty"`
	Reference       string            `json:"reference,omitempty"`
	ObservedAt      time.Time         `json:"observed_at"`
	RetrievedAt     time.Time         `json:"retrieved_at"`
	Country         string            `json:"country"`
	Province        string            `json:"province,omitempty"`
	Municipality    string            `json:"municipality,omitempty"`
	Locality        string            `json:"locality,omitempty"`
	Category        string            `json:"category,omitempty"`
	Subcategory     string            `json:"subcategory,omitempty"`
	Entity          string            `json:"entity,omitempty"`
	Product         string            `json:"product,omitempty"`
	Service         string            `json:"service,omitempty"`
	ObservationType ObservationType   `json:"observation_type"`
	Price           *float64          `json:"price,omitempty"`
	Currency        string            `json:"currency,omitempty"`
	Quantity        *float64          `json:"quantity,omitempty"`
	Text            string            `json:"text,omitempty"`
	Fingerprint     string            `json:"fingerprint"`
	Confidence      float64           `json:"confidence"` // 0..1
	Coverage        float64           `json:"coverage"`   // 0..1 estimate of market coverage of this signal
	Provenance      map[string]string `json:"provenance,omitempty"`
	Evidence        []string          `json:"evidence,omitempty"`
	Metadata        map[string]any    `json:"metadata,omitempty"`
}

type ObservationType string

const (
	ObsListing   ObservationType = "listing"
	ObsIntent    ObservationType = "intent" // busco/compro/necesito — NOT a transaction
	ObsPrice     ObservationType = "price_quote"
	ObsSearch    ObservationType = "search_interest"
	ObsNews      ObservationType = "news"
	ObsOfficial  ObservationType = "official"
	ObsUser      ObservationType = "user_provided"
	ObsUnknown   ObservationType = "unknown"
)

// MarketLayer distinguishes REAL vs OBSERVED market claims.
type MarketLayer string

const (
	LayerObserved MarketLayer = "observed_market"
	LayerOfficial MarketLayer = "official"
	LayerEstimate MarketLayer = "estimate"
	LayerForecast MarketLayer = "forecast"
)
