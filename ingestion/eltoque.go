// Package ingestion: legal/public market sensors only.
// El Toque TRMI requires an API token (request at tasas.eltoque.com). Never scrape private groups.
package ingestion

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/yecharlot/Pulso-Cubano/domain"
)

// ElToqueClient reads informal FX rates via official API when ELTOQUE_API_TOKEN is set.
type ElToqueClient struct {
	BaseURL string
	Token   string
	HTTP    *http.Client
}

func NewElToqueFromEnv() *ElToqueClient {
	tok := os.Getenv("ELTOQUE_API_TOKEN")
	base := os.Getenv("ELTOQUE_API_BASE")
	if base == "" {
		base = "https://tasas.eltoque.com"
	}
	return &ElToqueClient{
		BaseURL: base,
		Token:   tok,
		HTTP:    &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *ElToqueClient) Enabled() bool { return c != nil && c.Token != "" }

// FetchTRMI returns observations for informal FX (observed_market, not official BNC).
// Endpoint path may vary by API version — configure ELTOQUE_API_PATH if needed.
func (c *ElToqueClient) FetchTRMI(ctx context.Context) ([]domain.Observation, error) {
	if !c.Enabled() {
		return nil, fmt.Errorf("eltoque: ELTOQUE_API_TOKEN not set (request token at tasas.eltoque.com)")
	}
	path := os.Getenv("ELTOQUE_API_PATH")
	if path == "" {
		path = "/v1/trmi" // placeholder; operator must align with their token docs
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Accept", "application/json")
	res, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("eltoque: HTTP %d: %s", res.StatusCode, truncate(string(body), 200))
	}
	return parseElToqueJSON(body)
}

func parseElToqueJSON(body []byte) ([]domain.Observation, error) {
	// Flexible parse: accept { "tasas": { "USD": 720 } } or array of rates
	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	var out []domain.Observation
	if tasas, ok := raw["tasas"].(map[string]any); ok {
		for cur, v := range tasas {
			f, ok := toFloat(v)
			if !ok {
				continue
			}
			p := f
			out = append(out, domain.Observation{
				ID: fmt.Sprintf("eltoque-%s-%d", cur, now.Unix()),
				SourceID: "eltoque-trmi", SourceType: domain.SourceMarketData,
				ObservedAt: now, RetrievedAt: now,
				Country: "CU", Category: "fx", Product: cur,
				ObservationType: domain.ObsPrice, Price: &p, Currency: "CUP",
				Text: fmt.Sprintf("TRMI %s → CUP (elTOQUE, referential)", cur),
				Fingerprint: fmt.Sprintf("eltoque:%s:%d", cur, now.Unix()/3600),
				Confidence: 0.7, Coverage: 0.4,
				Provenance: map[string]string{
					"layer": string(domain.LayerObserved),
					"source": "eltoque-api",
					"disclaimer": "referential_informal_fx_not_official",
				},
			})
		}
	}
	return out, nil
}

func toFloat(v any) (float64, bool) {
	switch t := v.(type) {
	case float64:
		return t, true
	case json.Number:
		f, err := t.Float64()
		return f, err == nil
	default:
		return 0, false
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
