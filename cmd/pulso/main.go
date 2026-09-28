package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/yecharlot/Pulso-Cubano/metrics"
	"github.com/yecharlot/Pulso-Cubano/queries"
	"github.com/yecharlot/Pulso-Cubano/sources"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(`Pulso Cubano — market intelligence vertical on PrismaTec-Core

Usage:
  pulso version
  pulso sources
  pulso queries gtmo
  pulso pulse-demo`)
		os.Exit(0)
	}
	switch os.Args[1] {
	case "version":
		fmt.Println("pulso-cubano 0.1.0-mvp (observatory scaffold)")
	case "sources":
		for _, s := range sources.SeedRegistry().List() {
			fmt.Printf("%-16s %-18s active=%v access=%s\n", s.ID, s.Type, s.Active, s.Access)
		}
	case "queries":
		qs := queries.ExpandAll(queries.GTMOSeedTemplates(), nil)
		for _, q := range qs {
			fmt.Println(q)
		}
	case "pulse-demo":
		// Demo only — not real market data
		p := metrics.ComputePulse(72, 48, 65, 55, 40, 70, 60, 25)
		b, _ := json.MarshalIndent(p, "", "  ")
		fmt.Println(string(b))
		fmt.Println("layer=observed_market — not official GDP or sales")
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
