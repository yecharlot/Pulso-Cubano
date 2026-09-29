package api

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/yecharlot/Pulso-Cubano/domain"
	"github.com/yecharlot/Pulso-Cubano/sources"
)

// TerminalState feeds the economic intelligence terminal UI.
// Explicit: demo vs observed vs empty. No fabricated market activity.
type TerminalState struct {
	GeneratedAt   string            `json:"generated_at"`
	Market        string            `json:"market"`
	Territory     string            `json:"territory"`
	Period        string            `json:"period"`
	Mode          string            `json:"mode"` // live | partial | empty
	CutOff        string            `json:"cutoff"`
	Disclaimer    string            `json:"disclaimer"`
	HeaderMetrics []HeaderMetric    `json:"header_metrics"`
	SourceHealth  []SourceHealth    `json:"source_health"`
	Observations  []ObsView         `json:"observations"`
	Signals       []SignalView      `json:"signals"`
	News          []ObsView         `json:"news"`
	Prices        []ObsView         `json:"prices"`
	Sections      []SectionMeta     `json:"sections"`
	Nav           []NavItem         `json:"nav"`
}

type HeaderMetric struct {
	ID     string `json:"id"`
	Label  string `json:"label"`
	Value  string `json:"value"`
	Unit   string `json:"unit,omitempty"`
	State  string `json:"state"` // live | empty | delayed
	Hint   string `json:"hint"`
}

type SourceHealth struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Active    bool   `json:"active"`
	Access    string `json:"access"`
	Status    string `json:"status"` // connected | pending | disabled | error
	Note      string `json:"note"`
}

type ObsView struct {
	ID         string  `json:"id"`
	Source     string  `json:"source"`
	Type       string  `json:"type"`
	Category   string  `json:"category,omitempty"`
	Product    string  `json:"product,omitempty"`
	Text       string  `json:"text"`
	Price      *float64 `json:"price,omitempty"`
	Currency   string  `json:"currency,omitempty"`
	Confidence float64 `json:"confidence"`
	Layer      string  `json:"layer"`
	At         string  `json:"at"`
	URL        string  `json:"url,omitempty"`
}

type SignalView struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Kind     string `json:"kind"`
	Evidence string `json:"evidence"`
	State    string `json:"state"` // observed | pending_analysis
	At       string `json:"at"`
}

type SectionMeta struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Audience    string `json:"audience"`
	Status      string `json:"status"` // ready | partial | scaffold
	Description string `json:"description"`
}

type NavItem struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Group string `json:"group"`
}

func terminalNav() []NavItem {
	return []NavItem{
		{ID: "overview", Label: "01 Overview", Group: "terminal"},
		{ID: "radar", Label: "02 Radar", Group: "terminal"},
		{ID: "explorer", Label: "03 Explorer", Group: "terminal"},
		{ID: "prices", Label: "04 Prices", Group: "intelligence"},
		{ID: "supply", Label: "05 Supply & Demand", Group: "intelligence"},
		{ID: "directory", Label: "06 Directory", Group: "intelligence"},
		{ID: "network", Label: "07 Network", Group: "intelligence"},
		{ID: "sectors", Label: "08 Sectors", Group: "intelligence"},
		{ID: "news", Label: "09 News", Group: "intelligence"},
		{ID: "opportunities", Label: "10 Opportunities", Group: "ops"},
		{ID: "business", Label: "11 Business Terminal", Group: "ops"},
		{ID: "campaigns", Label: "12 Campaigns", Group: "ops"},
		{ID: "reports", Label: "13 Reports", Group: "ops"},
		{ID: "methodology", Label: "14 Methodology", Group: "infra"},
		{ID: "api", Label: "15 API & Data", Group: "infra"},
	}
}

func sectionCatalog() []SectionMeta {
	return []SectionMeta{
		{ID: "overview", Title: "Market Overview", Audience: "Todos", Status: "partial", Description: "Estado general, cobertura y acceso rápido"},
		{ID: "radar", Title: "Market Radar", Audience: "Empresarios, consultores", Status: "partial", Description: "Señales y cambios que merecen investigación"},
		{ID: "explorer", Title: "Market Explorer", Audience: "Todos", Status: "scaffold", Description: "Consulta multidimensional geo × sector × tiempo"},
		{ID: "prices", Title: "Price Intelligence", Audience: "Comerciantes, proveedores", Status: "partial", Description: "Precios anunciados con muestra y antigüedad"},
		{ID: "supply", Title: "Supply & Demand", Audience: "Negocios", Status: "scaffold", Description: "Ofertas vs solicitudes explícitas por separado"},
		{ID: "directory", Title: "Business Directory", Audience: "Todos", Status: "scaffold", Description: "Actores confirmados / probables / perfil"},
		{ID: "network", Title: "Economic Network", Audience: "Consultores", Status: "scaffold", Description: "Relaciones declaradas, observadas o asociativas"},
		{ID: "sectors", Title: "Sector Intelligence", Audience: "Inversores, consultores", Status: "scaffold", Description: "Estudios por industria con capas de evidencia"},
		{ID: "news", Title: "News & Events", Audience: "Todos", Status: "partial", Description: "Contexto económico con enlace a fuente"},
		{ID: "opportunities", Title: "Opportunity Research", Audience: "Emprendedores, inversores", Status: "scaffold", Description: "Hipótesis con evidencia y estado de validación"},
		{ID: "business", Title: "Business Terminal", Audience: "Comerciantes", Status: "scaffold", Description: "Espacio privado del negocio"},
		{ID: "campaigns", Title: "Campaign Intelligence", Audience: "Negocios, agencias", Status: "scaffold", Description: "De señal a prueba medible"},
		{ID: "reports", Title: "Reports & Research", Audience: "Profesionales", Status: "scaffold", Description: "Biblioteca de informes versionados"},
		{ID: "methodology", Title: "Data & Methodology", Audience: "Analistas", Status: "ready", Description: "Fuentes, sesgos, definiciones"},
		{ID: "api", Title: "API & Data Products", Audience: "Instituciones, devs", Status: "partial", Description: "Endpoints y productos de datos"},
	}
}

func (s *Server) handleTerminal(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, s.buildTerminal())
}

func (s *Server) buildTerminal() TerminalState {
	now := time.Now().UTC()
	obs := s.loadObservations()
	prices, news := splitObs(obs)
	signals := signalsFromObs(obs)

	reg := sources.SeedRegistry().List()
	health := make([]SourceHealth, 0, len(reg))
	activeN := 0
	for _, src := range reg {
		st := "pending"
		note := "Registrada; collector no conectado o sin token"
		if !src.Active {
			st = "disabled"
			note = "Desactivada hasta acceso legal estable"
		} else if src.ID == "frankfurter" || src.ID == "wikipedia" || src.ID == "news-rss" || src.ID == "manual-csv" {
			st = "connected"
			note = "Gene público disponible"
			activeN++
		} else if src.ID == "eltoque-trmi" || src.ID == "qvapay-p2p" || src.ID == "google-cse" {
			st = "pending"
			note = "Requiere credencial de operador"
			activeN++
		}
		health = append(health, SourceHealth{
			ID: src.ID, Name: src.Name, Type: string(src.Type), Active: src.Active,
			Access: string(src.Access), Status: st, Note: note,
		})
	}
	sort.Slice(health, func(i, j int) bool { return health[i].ID < health[j].ID })

	mode := "empty"
	if len(obs) > 0 {
		mode = "partial"
	}

	actVal, actState := "—", "empty"
	if len(obs) > 0 {
		actVal = itoa(len(obs))
		actState = "live"
	}
	priceVal, priceState := "—", "empty"
	if len(prices) > 0 {
		priceVal = itoa(len(prices))
		priceState = "live"
	}
	sigVal, sigState := "—", "empty"
	if len(signals) > 0 {
		sigVal = itoa(len(signals))
		sigState = "live"
	}

	return TerminalState{
		GeneratedAt: now.Format(time.RFC3339),
		Market:      "Mercado general",
		Territory:   "Cuba / Guantánamo",
		Period:      "Última corrida de sensores",
		Mode:        mode,
		CutOff:      now.Format("2006-01-02 15:04 UTC"),
		Disclaimer:  "Terminal de inteligencia sobre mercado OBSERVADO. No es bolsa ni contabilidad oficial. Cifras solo cuando hay observaciones capturadas; el resto muestra estado vacío o scaffold de producto.",
		HeaderMetrics: []HeaderMetric{
			{ID: "activity", Label: "Actividad observada", Value: actVal, Unit: "obs", State: actState, Hint: "Publicaciones y señales capturadas"},
			{ID: "prices", Label: "Precios", Value: priceVal, Unit: "muestras", State: priceState, Hint: "Últimas observaciones con precio"},
			{ID: "sources", Label: "Fuentes activas", Value: itoa(activeN), Unit: "reg.", State: "live", Hint: "En registro (no todas conectadas)"},
			{ID: "signals", Label: "Señales nuevas", Value: sigVal, State: sigState, Hint: "Pendientes de análisis humano"},
		},
		SourceHealth: health,
		Observations: obs,
		Signals:      signals,
		News:         news,
		Prices:       prices,
		Sections:     sectionCatalog(),
		Nav:          terminalNav(),
	}
}

func (s *Server) loadObservations() []ObsView {
	org, err := s.Adapter.Node.Organisms().Get(s.OrgID)
	if err != nil || org.Memory.Working == nil {
		return nil
	}
	var out []ObsView
	for k, raw := range org.Memory.Working {
		if !strings.HasPrefix(k, "obs:") || raw == "" {
			continue
		}
		var o domain.Observation
		if json.Unmarshal([]byte(raw), &o) != nil {
			continue
		}
		layer := "observed_market"
		if o.Provenance != nil {
			if l, ok := o.Provenance["layer"]; ok {
				layer = l
			}
		}
		at := o.ObservedAt
		if at.IsZero() {
			at = o.RetrievedAt
		}
		out = append(out, ObsView{
			ID: o.ID, Source: o.SourceID, Type: string(o.ObservationType),
			Category: o.Category, Product: o.Product, Text: o.Text,
			Price: o.Price, Currency: o.Currency, Confidence: o.Confidence,
			Layer: layer, At: at.UTC().Format(time.RFC3339), URL: o.URL,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].At > out[j].At })
	if len(out) > 80 {
		out = out[:80]
	}
	return out
}

func splitObs(obs []ObsView) (prices, news []ObsView) {
	for _, o := range obs {
		if o.Price != nil {
			prices = append(prices, o)
		}
		if o.Type == string(domain.ObsNews) || o.Category == "news" || o.Category == "reference" {
			news = append(news, o)
		}
	}
	return
}

func signalsFromObs(obs []ObsView) []SignalView {
	var out []SignalView
	for _, o := range obs {
		title := o.Text
		if len(title) > 90 {
			title = title[:90] + "…"
		}
		if title == "" {
			title = o.Product + " · " + o.Source
		}
		out = append(out, SignalView{
			ID: o.ID, Title: title, Kind: o.Type,
			Evidence: "fuente " + o.Source + " · confianza " + pct(o.Confidence),
			State: "observed", At: o.At,
		})
		if len(out) >= 25 {
			break
		}
	}
	return out
}

func itoa(n int) string {
	return strconv.Itoa(n)
}

func pct(f float64) string {
	return strconv.FormatFloat(f*100, 'f', 0, 64) + "%"
}
