package ingestion

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/yecharlot/Pulso-Cubano/domain"
)

// Google CSE requires GOOGLE_CSE_API_KEY + GOOGLE_CSE_CX.
func FetchGoogleCSE(ctx context.Context, query string) ([]domain.Observation, error) {
	key := os.Getenv("GOOGLE_CSE_API_KEY")
	cx := os.Getenv("GOOGLE_CSE_CX")
	if key == "" || cx == "" {
		return nil, fmt.Errorf("google cse: set GOOGLE_CSE_API_KEY and GOOGLE_CSE_CX")
	}
	u := fmt.Sprintf(
		"https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s",
		url.QueryEscape(key), url.QueryEscape(cx), url.QueryEscape(query),
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
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
		return nil, fmt.Errorf("google cse: HTTP %d: %s", res.StatusCode, truncate(string(body), 180))
	}
	var raw struct {
		Items []struct {
			Title string `json:"title"`
			Link  string `json:"link"`
			Snippet string `json:"snippet"`
		} `json:"items"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	var out []domain.Observation
	for i, it := range raw.Items {
		if i >= 8 {
			break
		}
		out = append(out, domain.Observation{
			ID: fmt.Sprintf("gsearch-%d-%d", now.Unix(), i),
			SourceID: "google-cse", SourceType: domain.SourceSearch,
			URL: it.Link, ObservedAt: now, RetrievedAt: now, Country: "CU",
			Province: "Guantánamo", Category: "search", Product: query,
			ObservationType: domain.ObsSearch,
			Text: it.Title + " — " + it.Snippet,
			Fingerprint: fmt.Sprintf("gse:%s", it.Link),
			Confidence: 0.45, Coverage: 0.1,
			Provenance: map[string]string{"layer": string(domain.LayerObserved), "query": query},
		})
	}
	return out, nil
}
