package queries

import "github.com/yecharlot/Pulso-Cubano/domain"

// GTMOSeedTemplates — initial Guantánamo observatory queries (expandable).
func GTMOSeedTemplates() []domain.QueryTemplate {
	muni := "Guantánamo"
	base := []struct {
		id, tpl, intent, cat string
	}{
		{"gtmo-venta", "{{municipality}} venta", "offer", "general"},
		{"gtmo-vendo", "{{municipality}} vendo", "offer", "general"},
		{"gtmo-compro", "{{municipality}} compro", "intent", "general"},
		{"gtmo-busco", "{{municipality}} busco", "intent", "general"},
		{"gtmo-necesito", "{{municipality}} necesito", "intent", "general"},
		{"gtmo-precio", "{{municipality}} precio", "price", "general"},
		{"gtmo-negocio", "{{municipality}} negocio", "business", "general"},
		{"gtmo-oferta", "{{municipality}} oferta", "offer", "general"},
		{"gtmo-telefono", "{{municipality}} teléfono", "product", "electronics"},
		{"gtmo-celular", "{{municipality}} celular", "product", "electronics"},
		{"gtmo-samsung", "{{municipality}} Samsung", "product", "electronics"},
		{"gtmo-xiaomi", "{{municipality}} Xiaomi", "product", "electronics"},
		{"gtmo-iphone", "{{municipality}} iPhone", "product", "electronics"},
		{"gtmo-laptop", "{{municipality}} laptop", "product", "electronics"},
		{"gtmo-moto", "{{municipality}} moto", "product", "automotive"},
		{"gtmo-motorina", "{{municipality}} motorina", "product", "automotive"},
		{"gtmo-casa-venta", "{{municipality}} casa venta", "product", "housing"},
		{"gtmo-casa-alquiler", "{{municipality}} casa alquiler", "product", "housing"},
		{"gtmo-pollo", "{{municipality}} pollo", "product", "food"},
		{"gtmo-reparacion", "{{municipality}} reparación", "service", "services"},
	}
	out := make([]domain.QueryTemplate, 0, len(base))
	for _, b := range base {
		out = append(out, domain.QueryTemplate{
			ID: b.id, Template: b.tpl, Intent: b.intent, Category: b.cat,
			Vars: map[string]string{"municipality": muni, "province": "Guantánamo", "country": "CU"},
		})
	}
	return out
}

func ExpandAll(templates []domain.QueryTemplate, extra map[string]string) []string {
	qs := make([]string, 0, len(templates))
	for _, t := range templates {
		qs = append(qs, t.Expand(extra))
	}
	return qs
}
