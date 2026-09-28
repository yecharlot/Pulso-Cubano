package sources

import "github.com/yecharlot/Pulso-Cubano/domain"

// SeedRegistry registers planned sensors. Collectors are NOT implemented here yet —
// only metadata for legal/public access posture.
func SeedRegistry() *domain.SourceRegistry {
	return domain.NewSourceRegistry(
		domain.Source{
			ID: "web-search", Name: "Web Search Sensor", Type: domain.SourceSearch,
			Coverage: "Cuba / public web results", Frequency: "scheduled",
			Reliability: 0.5, FreshnessHrs: 24, Access: domain.AccessAllowed,
			LegalNotes: "Public search results only; no private credentials.",
			Adapter: "search.public", Active: true,
		},
		domain.Source{
			ID: "revolico", Name: "Revolico", Type: domain.SourceClassifieds,
			Coverage: "Cuba classifieds (public pages)", Frequency: "scheduled",
			Reliability: 0.55, FreshnessHrs: 12, Access: domain.AccessLimited,
			LegalNotes: "Public listings only; respect robots/ToS; no account takeover.",
			Adapter: "classifieds.revolico", Active: false, // collector not wired in Phase 0-5
		},
		domain.Source{
			ID: "cubisima", Name: "Cubisima", Type: domain.SourceClassifieds,
			Coverage: "Housing / classifieds public", Frequency: "scheduled",
			Reliability: 0.5, FreshnessHrs: 24, Access: domain.AccessLimited,
			Adapter: "classifieds.cubisima", Active: false,
		},
		domain.Source{
			ID: "google-trends", Name: "Google Trends", Type: domain.SourceTrendData,
			Coverage: "Search interest proxies", Frequency: "daily",
			Reliability: 0.45, FreshnessHrs: 48, Access: domain.AccessLimited,
			Adapter: "trends.google", Active: false,
		},
		domain.Source{
			ID: "manual-csv", Name: "Manual datasets", Type: domain.SourceManualDataset,
			Coverage: "User-provided / research imports", Frequency: "on-demand",
			Reliability: 0.7, FreshnessHrs: 0, Access: domain.AccessAllowed,
			Adapter: "manual.import", Active: true,
		},
		domain.Source{
			ID: "first-party", Name: "Business first-party", Type: domain.SourceFirstParty,
			Coverage: "Voluntary business data", Frequency: "on-demand",
			Reliability: 0.8, FreshnessHrs: 0, Access: domain.AccessAllowed,
			Adapter: "firstparty.api", Active: true,
		},
	)
}
