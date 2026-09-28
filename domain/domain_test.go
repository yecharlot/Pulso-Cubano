package domain

import "testing"

func TestQueryExpand(t *testing.T) {
	q := QueryTemplate{Template: "{{municipality}} {{intent}} {{product}}", Vars: map[string]string{"intent": "busco"}}
	got := q.Expand(map[string]string{"municipality": "Baracoa", "product": "moto"})
	if got != "Baracoa busco moto" {
		t.Fatalf("%q", got)
	}
}

func TestSourceRegistry(t *testing.T) {
	r := NewSourceRegistry(Source{ID: "revolico", Name: "Revolico", Type: SourceClassifieds, Active: true})
	s, ok := r.Get("revolico")
	if !ok || s.Type != SourceClassifieds {
		t.Fatal(s)
	}
}
