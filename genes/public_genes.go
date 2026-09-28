package genes

import (
	"context"
	"time"

	"github.com/yecharlot/Pulso-Cubano/adapter/prismatec"
	"github.com/yecharlot/Pulso-Cubano/ingestion"
)

type GlobalFXGene struct{}

func (GlobalFXGene) Name() string { return "gene.fx.global" }

func (g GlobalFXGene) Run(ctx context.Context, ad *prismatec.Adapter, orgID string) (GeneResult, error) {
	res := GeneResult{Name: g.Name(), At: time.Now().UTC()}
	obs, err := ingestion.FetchFrankfurterUSD(ctx)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		return res, nil
	}
	return ingestAndMind(ctx, ad, orgID, g.Name(), obs, res)
}

type WikiGene struct{}

func (WikiGene) Name() string { return "gene.wiki.ref" }

func (g WikiGene) Run(ctx context.Context, ad *prismatec.Adapter, orgID string) (GeneResult, error) {
	res := GeneResult{Name: g.Name(), At: time.Now().UTC()}
	obs, err := ingestion.FetchWikipediaSummary(ctx, "Economía_de_Cuba")
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		return res, nil
	}
	return ingestAndMind(ctx, ad, orgID, g.Name(), obs, res)
}

type NewsRSSGene struct{}

func (NewsRSSGene) Name() string { return "gene.news.rss" }

func (g NewsRSSGene) Run(ctx context.Context, ad *prismatec.Adapter, orgID string) (GeneResult, error) {
	res := GeneResult{Name: g.Name(), At: time.Now().UTC()}
	for _, f := range ingestion.DefaultCubaNewsFeeds() {
		obs, err := ingestion.FetchRSS(ctx, f.URL, f.ID, 4)
		if err != nil {
			res.Errors = append(res.Errors, f.ID+": "+err.Error())
			continue
		}
		partial, err := ingestAndMind(ctx, ad, orgID, g.Name(), obs, GeneResult{Name: g.Name(), At: time.Now().UTC()})
		if err != nil {
			res.Errors = append(res.Errors, err.Error())
			continue
		}
		res.Observations += partial.Observations
		if partial.MindSelected != "" {
			res.MindSelected = partial.MindSelected
		}
		res.PulseOverall = partial.PulseOverall
	}
	return res, nil
}

type SearchWebGene struct {
	Query string
}

func (SearchWebGene) Name() string { return "gene.search.web" }

func (g SearchWebGene) Run(ctx context.Context, ad *prismatec.Adapter, orgID string) (GeneResult, error) {
	res := GeneResult{Name: g.Name(), At: time.Now().UTC()}
	q := g.Query
	if q == "" {
		q = "Guantánamo venta precio negocio"
	}
	obs, err := ingestion.FetchGoogleCSE(ctx, q)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		return res, nil
	}
	return ingestAndMind(ctx, ad, orgID, g.Name(), obs, res)
}
