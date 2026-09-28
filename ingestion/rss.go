package ingestion

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/yecharlot/Pulso-Cubano/domain"
)

type rssFeed struct {
	Channel struct {
		Title string `xml:"title"`
		Item  []struct {
			Title   string `xml:"title"`
			Link    string `xml:"link"`
			PubDate string `xml:"pubDate"`
		} `xml:"item"`
	} `xml:"channel"`
}

// FetchRSS pulls public RSS/Atom-ish RSS 2.0 items as news observations.
func FetchRSS(ctx context.Context, feedURL, sourceID string, limit int) ([]domain.Observation, error) {
	if limit <= 0 {
		limit = 5
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "PulsoCubano/0.2")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("rss: HTTP %d", res.StatusCode)
	}
	var feed rssFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("rss parse: %w", err)
	}
	now := time.Now().UTC()
	var out []domain.Observation
	for i, it := range feed.Channel.Item {
		if i >= limit {
			break
		}
		out = append(out, domain.Observation{
			ID: fmt.Sprintf("%s-%d-%d", sourceID, now.Unix(), i),
			SourceID: sourceID, SourceType: domain.SourcePublicWeb,
			URL: it.Link, ObservedAt: now, RetrievedAt: now, Country: "CU",
			Category: "news", Product: "headline",
			ObservationType: domain.ObsNews,
			Text: it.Title, Fingerprint: fmt.Sprintf("rss:%s:%s", sourceID, it.Link),
			Confidence: 0.5, Coverage: 0.15,
			Provenance: map[string]string{"layer": string(domain.LayerObserved), "feed": feedURL},
		})
	}
	return out, nil
}

// DefaultCubaNewsFeeds — public feeds (may change; errors are soft in genes).
func DefaultCubaNewsFeeds() []struct{ ID, URL string } {
	return []struct{ ID, URL string }{
		{"bbc-mundo", "https://feeds.bbci.co.uk/mundo/rss.xml"},
		{"reuters-world", "http://feeds.reuters.com/Reuters/worldNews"},
	}
}
