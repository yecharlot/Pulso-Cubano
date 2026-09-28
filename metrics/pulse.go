package metrics

// PulseScore is the explainable market pulse — never a black box.
// All components are OBSERVED-market metrics unless labeled otherwise.
type PulseScore struct {
	Demand          float64 `json:"demand"`
	Supply          float64 `json:"supply"`
	Momentum        float64 `json:"momentum"`
	Price           float64 `json:"price"`
	Competition     float64 `json:"competition"`
	SearchInterest  float64 `json:"search_interest"`
	Evidence        float64 `json:"evidence"`
	Overall         float64 `json:"overall"`
	Layer           string  `json:"layer"` // observed_market
	MethodVersion   string  `json:"method_version"`
	ObservationN    int     `json:"observation_n"`
	Notes           string  `json:"notes,omitempty"`
}

// ComputePulse is a transparent weighted average for MVP (v0).
// Returns zeros if n==0 — does not invent activity.
func ComputePulse(demand, supply, momentum, price, competition, search, evidence float64, n int) PulseScore {
	if n <= 0 {
		return PulseScore{Layer: "observed_market", MethodVersion: "pulse-v0"}
	}
	w := []float64{0.18, 0.15, 0.15, 0.12, 0.12, 0.13, 0.15}
	parts := []float64{demand, supply, momentum, price, competition, search, evidence}
	var sum float64
	for i := range parts {
		sum += w[i] * parts[i]
	}
	return PulseScore{
		Demand: demand, Supply: supply, Momentum: momentum, Price: price,
		Competition: competition, SearchInterest: search, Evidence: evidence,
		Overall: sum, Layer: "observed_market", MethodVersion: "pulse-v0", ObservationN: n,
	}
}
