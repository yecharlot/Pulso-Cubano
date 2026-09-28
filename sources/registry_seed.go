package sources

import "github.com/yecharlot/Pulso-Cubano/domain"

// SeedRegistry — universo de sensores (metadatos). Collectors se activan por gene + credenciales legales.
func SeedRegistry() *domain.SourceRegistry {
	return domain.NewSourceRegistry(
		// --- Official ---
		domain.Source{
			ID: "bcc-official", Name: "Banco Central de Cuba (público)", Type: domain.SourceOfficial,
			Coverage: "Tipos/comunicados oficiales cuando estén publicados", Frequency: "as-published",
			Reliability: 0.9, FreshnessHrs: 24, Access: domain.AccessAllowed,
			LegalNotes: "Solo información publicada en portales oficiales.",
			Adapter: "gene.fx.official", Active: true,
		},
		domain.Source{
			ID: "onei-stats", Name: "ONEI / estadísticas públicas", Type: domain.SourceOfficial,
			Coverage: "Indicadores macro publicados", Frequency: "periodic",
			Reliability: 0.85, FreshnessHrs: 168, Access: domain.AccessAllowed,
			Adapter: "gene.econ.official", Active: true,
		},
		domain.Source{
			ID: "gaceta", Name: "Gaceta Oficial", Type: domain.SourceOfficial,
			Coverage: "Cambios regulatorios públicos", Frequency: "as-published",
			Reliability: 0.9, FreshnessHrs: 48, Access: domain.AccessAllowed,
			Adapter: "gene.reg.gaceta", Active: true,
		},
		// --- Informal FX reference ---
		domain.Source{
			ID: "eltoque-trmi", Name: "elTOQUE TRMI API", Type: domain.SourceMarketData,
			Coverage: "FX informal referencial Cuba", Frequency: "hourly",
			Reliability: 0.75, FreshnessHrs: 1, Access: domain.AccessLimited,
			LegalNotes: "Token + ToS elTOQUE; citar fuente; no es tasa BNC.",
			Adapter: "gene.fx.eltoque", Active: true,
		},
		// --- Fintech observed ---
		domain.Source{
			ID: "qvapay-p2p", Name: "QvaPay P2P", Type: domain.SourceMarketData,
			Coverage: "Actividad P2P observada (API)", Frequency: "on-demand",
			Reliability: 0.6, FreshnessHrs: 1, Access: domain.AccessLimited,
			LegalNotes: "Token propio; mínimos privilegios.",
			Adapter: "gene.fx.qvapay", Active: true,
		},
		// --- Search & web ---
		domain.Source{
			ID: "google-cse", Name: "Google Programmable Search", Type: domain.SourceSearch,
			Coverage: "Resultados web parametrizados (API key propia)", Frequency: "scheduled",
			Reliability: 0.5, FreshnessHrs: 12, Access: domain.AccessLimited,
			LegalNotes: "Google CSE / Custom Search API; cuotas de Google.",
			Adapter: "gene.search.web", Active: true,
		},
		domain.Source{
			ID: "google-trends", Name: "Google Trends", Type: domain.SourceTrendData,
			Coverage: "Interés de búsqueda (proxies)", Frequency: "daily",
			Reliability: 0.45, FreshnessHrs: 24, Access: domain.AccessLimited,
			Adapter: "gene.trends.google", Active: true,
		},
		// --- Classifieds public ---
		domain.Source{
			ID: "revolico", Name: "Revolico (público)", Type: domain.SourceClassifieds,
			Coverage: "Listings públicos Cuba", Frequency: "scheduled",
			Reliability: 0.55, FreshnessHrs: 12, Access: domain.AccessLimited,
			LegalNotes: "Solo páginas públicas; robots/ToS; sin cuentas ajenas.",
			Adapter: "gene.classifieds.revolico", Active: false,
		},
		domain.Source{
			ID: "cubisima", Name: "Cubisima (público)", Type: domain.SourceClassifieds,
			Coverage: "Vivienda/clasificados públicos", Frequency: "scheduled",
			Reliability: 0.5, FreshnessHrs: 24, Access: domain.AccessLimited,
			Adapter: "gene.classifieds.cubisima", Active: false,
		},
		// --- Knowledge ---
		domain.Source{
			ID: "wikipedia", Name: "Wikipedia / Wikidata", Type: domain.SourcePublicWeb,
			Coverage: "Geografía, entidades, contexto", Frequency: "on-demand",
			Reliability: 0.7, FreshnessHrs: 720, Access: domain.AccessAllowed,
			Adapter: "gene.wiki.ref", Active: true,
		},
		// --- News RSS ---
		domain.Source{
			ID: "news-rss", Name: "RSS medios económicos", Type: domain.SourcePublicWeb,
			Coverage: "Titulares/públicos sobre economía Cuba", Frequency: "hourly",
			Reliability: 0.5, FreshnessHrs: 6, Access: domain.AccessAllowed,
			Adapter: "gene.news.rss", Active: true,
		},
		// --- Global FX context ---
		domain.Source{
			ID: "frankfurter", Name: "Frankfurter / open FX", Type: domain.SourceMarketData,
			Coverage: "USD/EUR globales (contexto, no mercado cubano)", Frequency: "daily",
			Reliability: 0.8, FreshnessHrs: 24, Access: domain.AccessAllowed,
			Adapter: "gene.fx.global", Active: true,
		},
		// --- First party & manual ---
		domain.Source{
			ID: "first-party", Name: "Datos de negocio (SaaS)", Type: domain.SourceFirstParty,
			Coverage: "Cliente sube precios/actividad agregada", Frequency: "on-demand",
			Reliability: 0.85, FreshnessHrs: 0, Access: domain.AccessAllowed,
			Adapter: "gene.firstparty", Active: true,
		},
		domain.Source{
			ID: "manual-csv", Name: "Import CSV/JSON", Type: domain.SourceManualDataset,
			Coverage: "Datasets de investigación / operador", Frequency: "on-demand",
			Reliability: 0.7, FreshnessHrs: 0, Access: domain.AccessAllowed,
			Adapter: "gene.manual", Active: true,
		},
		// --- Social public only (explicit policy) ---
		domain.Source{
			ID: "social-public", Name: "Redes — solo contenido público", Type: domain.SourceSocialPublic,
			Coverage: "Páginas/canales públicos indexables", Frequency: "scheduled",
			Reliability: 0.35, FreshnessHrs: 24, Access: domain.AccessLimited,
			LegalNotes: "NO grupos privados FB/WA. Solo URLs públicas o export del titular.",
			Adapter: "gene.social.public", Active: false,
		},
	)
}
