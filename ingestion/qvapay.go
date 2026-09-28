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

// QvaPayClient reads public/aggregate market signals when QVAPAY_TOKEN is set.
// Uses documented API (api.qvapay.com). P2P averages are observed market activity, not GDP.
type QvaPayClient struct {
	BaseURL string
	Token   string
	HTTP    *http.Client
}

func NewQvaPayFromEnv() *QvaPayClient {
	base := os.Getenv("QVAPAY_API_BASE")
	if base == "" {
		base = "https://api.qvapay.com"
	}
	return &QvaPayClient{
		BaseURL: base,
		Token:   os.Getenv("QVAPAY_TOKEN"),
		HTTP:    &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *QvaPayClient) Enabled() bool { return c != nil && c.Token != "" }

// FetchP2PAverages maps public P2P average endpoints into Observations (layer observed).
func (c *QvaPayClient) FetchP2PAverages(ctx context.Context) ([]domain.Observation, error) {
	if !c.Enabled() {
		return nil, fmt.Errorf("qvapay: QVAPAY_TOKEN not set")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/p2p/average", nil)
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
		return nil, fmt.Errorf("qvapay: HTTP %d: %s", res.StatusCode, truncate(string(body), 200))
	}
	now := time.Now().UTC()
	// Store raw envelope as one observation + try numeric fields
	obs := domain.Observation{
		ID: fmt.Sprintf("qvapay-p2p-%d", now.Unix()),
		SourceID: "qvapay-p2p", SourceType: domain.SourceMarketData,
		ObservedAt: now, RetrievedAt: now,
		Country: "CU", Category: "fx_p2p", Product: "p2p_average",
		ObservationType: domain.ObsPrice,
		Text: "QvaPay P2P average snapshot (observed)",
		Fingerprint: fmt.Sprintf("qvapay:p2p:%d", now.Unix()/3600),
		Confidence: 0.55, Coverage: 0.25,
		Provenance: map[string]string{
			"layer": string(domain.LayerObserved),
			"source": "qvapay-api",
			"endpoint": "/p2p/average",
		},
		Metadata: map[string]any{"raw": json.RawMessage(body)},
	}
	var m map[string]any
	if json.Unmarshal(body, &m) == nil {
		if avg, ok := toFloat(m["average"]); ok {
			obs.Price = &avg
			obs.Currency = "CUP"
		}
	}
	return []domain.Observation{obs}, nil
}
