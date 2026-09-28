package ingestion

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/yecharlot/Pulso-Cubano/domain"
)

// Wikipedia summary for Cuba economy context (reference layer).
func FetchWikipediaSummary(ctx context.Context, title string) ([]domain.Observation, error) {
	if title == "" {
		title = "Economía_de_Cuba"
	}
	u := "https://es.wikipedia.org/api/rest_v1/page/summary/" + url.PathEscape(title)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "PulsoCubano/0.2 (market-intelligence; contact: local-dev)")
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("wikipedia: HTTP %d", res.StatusCode)
	}
	var raw struct {
		Title       string `json:"title"`
		Extract     string `json:"extract"`
		ContentURL  struct {
			Desktop struct {
				Page string `json:"page"`
			} `json:"desktop"`
		} `json:"content_urls"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	text := raw.Extract
	if len(text) > 500 {
		text = text[:500] + "…"
	}
	return []domain.Observation{{
		ID: fmt.Sprintf("wiki-%s-%d", title, now.Unix()/86400),
		SourceID: "wikipedia", SourceType: domain.SourcePublicWeb,
		URL: raw.ContentURL.Desktop.Page,
		ObservedAt: now, RetrievedAt: now, Country: "CU",
		Category: "reference", Product: raw.Title,
		ObservationType: domain.ObsNews,
		Text: text, Fingerprint: fmt.Sprintf("wiki:%s:%d", title, now.Unix()/86400),
		Confidence: 0.65, Coverage: 0.2,
		Provenance: map[string]string{"layer": "reference", "source": "wikipedia_rest"},
	}}, nil
}
