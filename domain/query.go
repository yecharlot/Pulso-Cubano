package domain

import (
	"strings"
)

// QueryTemplate is a parameterized search — not a hardcoded list of 100 queries.
type QueryTemplate struct {
	ID       string            `json:"id"`
	Template string            `json:"template"` // e.g. "{{municipality}} {{intent}} {{product}}"
	Intent   string            `json:"intent,omitempty"`
	Category string            `json:"category,omitempty"`
	Vars     map[string]string `json:"vars,omitempty"`
}

// Expand substitutes {{key}} from vars (municipality, product, province, intent, ...).
func (q QueryTemplate) Expand(vars map[string]string) string {
	s := q.Template
	merged := map[string]string{}
	for k, v := range q.Vars {
		merged[k] = v
	}
	for k, v := range vars {
		merged[k] = v
	}
	for k, v := range merged {
		s = strings.ReplaceAll(s, "{{"+k+"}}", v)
	}
	return strings.TrimSpace(s)
}
