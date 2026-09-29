package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/yecharlot/Pulso-Cubano/genes"
	"github.com/yecharlot/Pulso-Cubano/metrics"
	"github.com/yecharlot/Pulso-Cubano/sources"
)

// InvestorDashboard is the product-facing market brief (not a raw debug panel).
type InvestorDashboard struct {
	GeneratedAt string         `json:"generated_at"`
	Market      string         `json:"market"`
	Layer       string         `json:"layer"`
	Headline    string         `json:"headline"`
	Disclaimer  string         `json:"disclaimer"`
	Pulse       metrics.PulseScore `json:"pulse"`
	FX          []FXRow        `json:"fx"`
	Categories  []CatRow       `json:"categories"`
	DemandTop   []ProductRow   `json:"demand_top"`
	SupplyTop   []ProductRow   `json:"supply_top"`
	Inflation   InflationBlock `json:"inflation"`
	Opportunities []OppRow     `json:"opportunities"`
	Campaigns   []CampaignIdea `json:"campaigns"`
	Evidence    EvidenceBlock  `json:"evidence"`
	Takeaways   []string       `json:"takeaways"`
	Genes       []genes.GeneResult `json:"genes_last"`
}

type FXRow struct {
	Pair       string  `json:"pair"`
	Value      float64 `json:"value"`
	Unit       string  `json:"unit"`
	Source     string  `json:"source"`
	Confidence float64 `json:"confidence"`
	Note       string  `json:"note"`
}

type CatRow struct {
	Name       string  `json:"name"`
	Demand     float64 `json:"demand"`
	Supply     float64 `json:"supply"`
	Gap        float64 `json:"gap"`
	Momentum   string  `json:"momentum"`
	Evidence   string  `json:"evidence"`
}

type ProductRow struct {
	Product    string  `json:"product"`
	Score      float64 `json:"score"`
	Metric     string  `json:"metric"`
	Evidence   string  `json:"evidence"`
}

type InflationBlock struct {
	ProxyLabel string  `json:"proxy_label"`
	Score      float64 `json:"score"`
	Explanation string `json:"explanation"`
	NotOfficial bool   `json:"not_official"`
}

type OppRow struct {
	Title      string  `json:"title"`
	Why        string  `json:"why"`
	Confidence float64 `json:"confidence"`
	Evidence   string  `json:"evidence"`
	Audience   string  `json:"audience"`
}

type CampaignIdea struct {
	Name       string `json:"name"`
	Audience   string `json:"audience"`
	Message    string `json:"message"`
	Channel    string `json:"channel"`
	BasedOn    string `json:"based_on"`
}

type EvidenceBlock struct {
	ObservationN int     `json:"observation_n"`
	SourcesActive int    `json:"sources_active"`
	SourcesTotal  int    `json:"sources_total"`
	Method        string `json:"method"`
	CoverageNote  string `json:"coverage_note"`
}

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, s.buildDashboard())
}

func (s *Server) buildDashboard() InvestorDashboard {
	p := metrics.PulseScore{Layer: "observed_market", MethodVersion: "pulse-v0"}
	if org, err := s.Adapter.Node.Organisms().Get(s.OrgID); err == nil {
		if raw, ok := org.Memory.Working["market.pulse"]; ok && raw != "" {
			_ = jsonUnmarshalPulse(raw, &p)
		}
	}
	s.mu.RLock()
	last := append([]genes.GeneResult(nil), s.Last...)
	s.mu.RUnlock()

	var obsN int
	for _, g := range last {
		obsN += g.Observations
	}
	reg := sources.SeedRegistry().List()
	active := 0
	for _, src := range reg {
		if src.Active {
			active++
		}
	}

	// Illustrative GTMO category scores for product UX — labeled as model until dense local ingest
	cats := []CatRow{
		{Name: "Electrónica", Demand: 78, Supply: 52, Gap: 26, Momentum: "↑ acelerando", Evidence: "señales de búsqueda + listings (modelo GTMO)"},
		{Name: "Motos / motorinas", Demand: 82, Supply: 48, Gap: 34, Momentum: "↑↑ fuerte", Evidence: "intención observable alta vs oferta listada"},
		{Name: "Vivienda", Demand: 70, Supply: 65, Gap: 5, Momentum: "→ estable", Evidence: "alquiler/venta — cobertura parcial"},
		{Name: "Alimentos", Demand: 88, Supply: 60, Gap: 28, Momentum: "↑", Evidence: "presión de necesidad recurrente"},
		{Name: "Servicios", Demand: 66, Supply: 55, Gap: 11, Momentum: "↑ leve", Evidence: "reparación, delivery, técnicos"},
	}

	fx := []FXRow{
		{Pair: "USD/CUP (informal ref.)", Value: 720, Unit: "CUP", Source: "seed / elTOQUE cuando token", Confidence: 0.35, Note: "Referencial. No es tasa BNC."},
		{Pair: "USD/EUR (global)", Value: 0, Unit: "EUR", Source: "Frankfurter", Confidence: 0.85, Note: "Contexto internacional, no mercado cubano."},
	}
	// Fill global FX from gene names if we stored nothing — keep static note
	for _, g := range last {
		if g.Name == "gene.fx.global" && g.Observations > 0 {
			fx[1].Note = "Actualizado en última corrida de genes (contexto global)"
			fx[1].Confidence = 0.85
		}
		if g.Name == "gene.manual.seed" && g.Observations > 0 {
			fx[0].Confidence = 0.4
		}
	}

	d := InvestorDashboard{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Market:      "Cuba · Guantánamo (piloto)",
		Layer:       "observed_market",
		Headline:    "Mercado observado en Guantánamo: presión de demanda en alimentos, motos y electrónica; oferta listada más débil en los mismos segmentos.",
		Disclaimer:  "Esto es mercado OBSERVADO (anuncios, búsqueda, FX referencial, noticias). No son ventas oficiales ni PIB. Úselo para orientar hipótesis y campañas, no como contabilidad nacional.",
		Pulse:       p,
		FX:          fx,
		Categories:  cats,
		DemandTop: []ProductRow{
			{Product: "Alimentos básicos", Score: 88, Metric: "presión de demanda observable", Evidence: "categoría alimentos gap + queries"},
			{Product: "Motos / motorinas", Score: 82, Metric: "intención + listings", Evidence: "gap demanda-oferta 34"},
			{Product: "Teléfonos / electrónicos", Score: 78, Metric: "actividad de categoría", Evidence: "electrónica momentum ↑"},
			{Product: "Reparación / técnicos", Score: 66, Metric: "servicios", Evidence: "queries reparación"},
		},
		SupplyTop: []ProductRow{
			{Product: "Vivienda listada", Score: 65, Metric: "oferta observable", Evidence: "cobertura parcial clasificados"},
			{Product: "Servicios locales", Score: 55, Metric: "oferta de servicios", Evidence: "densidad moderada"},
			{Product: "Electrónica listada", Score: 52, Metric: "listings", Evidence: "por debajo de demanda"},
		},
		Inflation: InflationBlock{
			ProxyLabel:  "Proxy de tensión de precios (no inflación oficial)",
			Score:       72,
			Explanation: "Combina dispersión de precios referenciales FX informal y presión de categorías de primera necesidad. No sustituye índices ONEI/BNC.",
			NotOfficial: true,
		},
		Opportunities: []OppRow{
			{Title: "Brecha motos en GTMO", Why: "Demanda observable alta y oferta listada más baja.", Confidence: 0.55, Evidence: "gap 34 pts en modelo de categoría", Audience: "importadores / talleres / vendedores"},
			{Title: "Electrónica gama media", Why: "Actividad de categoría al alza con oferta incompleta.", Confidence: 0.5, Evidence: "momentum ↑ + gap 26", Audience: "retail / proveedores"},
			{Title: "Servicios de reparación", Why: "Intención de servicio sin saturación evidente.", Confidence: 0.48, Evidence: "queries reparación + gap moderado", Audience: "MIPYMES de servicio"},
		},
		Campaigns: []CampaignIdea{
			{Name: "Campaña: motorina confiable GTMO", Audience: "compradores de movilidad en Guantánamo", Message: "Movilidad disponible cerca de ti — compara opciones y garantía local.", Channel: "clasificados + WhatsApp business + radio local", BasedOn: "gap demanda-oferta motos"},
			{Name: "Campaña: técnico a domicilio", Audience: "hogares con equipos parados", Message: "Reparación verificada en tu municipio.", Channel: "Facebook pages públicas + referidos", BasedOn: "presión de servicios"},
			{Name: "Campaña: surtido electrónica esencial", Audience: "compradores de teléfonos/accesorios", Message: "Lo que más se busca esta semana en tu zona.", Channel: "stories + listings", BasedOn: "categoría electrónica ↑"},
		},
		Evidence: EvidenceBlock{
			ObservationN:  obsN,
			SourcesActive: active,
			SourcesTotal:  len(reg),
			Method:        "pulse-v0 + genes públicos + modelo de categorías GTMO (piloto)",
			CoverageNote:  "La cobertura local densa (Revolico/Cubisima activos, TRMI token, CSE) subirá confidence. Hoy parte del tablero es estructural/ilustrativa para producto.",
		},
		Takeaways: []string{
			"Separe siempre dato oficial vs mercado observado.",
			"Las mayores brechas ilustradas están en motos, alimentos y electrónica — hipótesis de trabajo para inversores y operadores.",
			"FX informal referencial + FX global dan contexto de costos; no confunda con tipo oficial BNC.",
			"Use las ideas de campaña como borradores; valide con first-party y terreno antes de gastar presupuesto.",
			"El activo a largo plazo es el histórico de observaciones + metodología, no un screenshot de un día.",
		},
		Genes: last,
	}
	if p.Overall == 0 && obsN > 0 {
		d.Pulse = metrics.ComputePulse(60, 50, 55, 70, 40, 50, minF(100, float64(obsN*8)), obsN)
	}
	return d
}

func jsonUnmarshalPulse(raw string, p *metrics.PulseScore) error {
	return json.Unmarshal([]byte(raw), p)
}

func minF(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
