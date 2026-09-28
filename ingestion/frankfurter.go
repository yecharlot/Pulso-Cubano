package ingestion

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/yecharlot/Pulso-Cubano/domain"
)

// Frankfurter: free ECB-based FX (context only — not Cuban informal market).
// https://www.frankfurter.app/
func FetchFrankfurterUSD(ctx context.Context) ([]domain.Observation, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://api.frankfurter.app/latest?from=USD&to=EUR,CAD,MXN", nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("frankfurter: HTTP %d", res.StatusCode)
	}
	var raw struct {
		Date  string             `json:"date"`
		Rates map[string]float64 `json:"rates"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	var out []domain.Observation
	for cur, rate := range raw.Rates {
		r := rate
		out = append(out, domain.Observation{
			ID: fmt.Sprintf("fx-global-usd-%s-%s", cur, raw.Date),
			SourceID: "frankfurter", SourceType: domain.SourceMarketData,
			ObservedAt: now, RetrievedAt: now, Country: "XX",
			Category: "fx_global", Product: "USD/"+cur,
			ObservationType: domain.ObsPrice, Price: &r, Currency: cur,
			Text: fmt.Sprintf("Global FX USD→%s = %v (%s)", cur, rate, raw.Date),
			Fingerprint: fmt.Sprintf("frankfurter:USD:%s:%s", cur, raw.Date),
			Confidence: 0.85, Coverage: 0.3,
			Provenance: map[string]string{
				"layer": string(domain.LayerObserved),
				"note":  "global_context_not_cuban_market",
			},
		})
	}
	return out, nil
}
