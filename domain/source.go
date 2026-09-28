package domain

import "time"

type SourceType string

const (
	SourcePublicWeb     SourceType = "PUBLIC_WEB"
	SourceSearch        SourceType = "SEARCH"
	SourceClassifieds   SourceType = "CLASSIFIEDS"
	SourceSocialPublic  SourceType = "SOCIAL_PUBLIC"
	SourceTelegramPub   SourceType = "TELEGRAM_PUBLIC"
	SourceOfficial      SourceType = "OFFICIAL"
	SourceFirstParty    SourceType = "FIRST_PARTY"
	SourceUserProvided  SourceType = "USER_PROVIDED"
	SourceMarketData    SourceType = "MARKET_DATA"
	SourceEconomicData  SourceType = "ECONOMIC_DATA"
	SourceTrendData     SourceType = "TREND_DATA"
	SourceManualDataset SourceType = "MANUAL_DATASET"
)

type AccessStatus string

const (
	AccessAllowed AccessStatus = "allowed"
	AccessLimited AccessStatus = "limited"
	AccessBlocked AccessStatus = "blocked"
	AccessUnknown AccessStatus = "unknown"
)

// Source describes a market sensor. Vertical-owned; Core is unaware of sources.
type Source struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Type         SourceType   `json:"type"`
	Coverage     string       `json:"coverage"` // e.g. "Cuba / Guantánamo / public listings"
	Frequency    string       `json:"frequency"`
	Reliability  float64      `json:"reliability"` // 0..1
	FreshnessHrs int          `json:"freshness_hrs"`
	Access       AccessStatus `json:"access"`
	LegalNotes   string       `json:"legal_notes,omitempty"`
	TermsURL     string       `json:"terms_url,omitempty"`
	Adapter      string       `json:"adapter"` // registry key for collector
	LastCapture  *time.Time   `json:"last_capture,omitempty"`
	LastError    string       `json:"last_error,omitempty"`
	Active       bool         `json:"active"`
}

// SourceRegistry is in-memory for MVP; later persisted via Core memory/organism.
type SourceRegistry struct {
	byID map[string]Source
}

func NewSourceRegistry(sources ...Source) *SourceRegistry {
	r := &SourceRegistry{byID: map[string]Source{}}
	for _, s := range sources {
		r.byID[s.ID] = s
	}
	return r
}

func (r *SourceRegistry) Get(id string) (Source, bool) {
	s, ok := r.byID[id]
	return s, ok
}

func (r *SourceRegistry) List() []Source {
	out := make([]Source, 0, len(r.byID))
	for _, s := range r.byID {
		out = append(out, s)
	}
	return out
}

func (r *SourceRegistry) Upsert(s Source) { r.byID[s.ID] = s }
