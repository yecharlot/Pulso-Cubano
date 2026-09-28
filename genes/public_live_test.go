package genes

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/yecharlot/Pulso-Cubano/adapter/prismatec"
)

func TestPublicGenesLiveNetwork(t *testing.T) {
	if testing.Short() {
		t.Skip("network")
	}
	dir := t.TempDir()
	ad, err := prismatec.NewAdapter(filepath.Join(dir, "n"), "live")
	if err != nil {
		t.Fatal(err)
	}
	org, err := ad.EnsureObserverOrganism("GTMOObserver")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	res, err := RunAll(ctx, ad, org.ID)
	if err != nil {
		t.Fatal(err)
	}
	var total int
	for _, r := range res {
		t.Logf("%s obs=%d errs=%v mind=%s", r.Name, r.Observations, r.Errors, r.MindSelected)
		total += r.Observations
	}
	if total == 0 {
		t.Fatal("expected some observations from public genes (frankfurter/wiki/rss)")
	}
}
