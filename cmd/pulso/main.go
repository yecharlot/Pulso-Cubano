package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/yecharlot/Pulso-Cubano/api"
	"github.com/yecharlot/Pulso-Cubano/metrics"
	"github.com/yecharlot/Pulso-Cubano/queries"
	"github.com/yecharlot/Pulso-Cubano/sources"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(`Pulso Cubano

  pulso version
  pulso sources
  pulso queries gtmo
  pulso pulse-demo
  pulso serve [addr]     # API + Alset-JS UI (default :8090)`)
		os.Exit(0)
	}
	switch os.Args[1] {
	case "version":
		fmt.Println("pulso-cubano 0.2.0 (alset-js ui + market genes)")
	case "sources":
		for _, s := range sources.SeedRegistry().List() {
			fmt.Printf("%-16s %-18s active=%v\n", s.ID, s.Type, s.Active)
		}
	case "queries":
		for _, q := range queries.ExpandAll(queries.GTMOSeedTemplates(), nil) {
			fmt.Println(q)
		}
	case "pulse-demo":
		p := metrics.ComputePulse(72, 48, 65, 55, 40, 70, 60, 25)
		b, _ := json.MarshalIndent(p, "", "  ")
		fmt.Println(string(b))
	case "serve":
		addr := ":8090"
		if len(os.Args) > 2 {
			addr = os.Args[2]
		}
		uiDir := os.Getenv("PULSO_UI")
		if uiDir == "" {
			uiDir = "ui"
		}
		data := os.Getenv("PULSO_DATA")
		if data == "" {
			data = filepath.Join(os.TempDir(), "pulso-cubano-data")
		}
		srv, err := api.NewServer(data, uiDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "server: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Pulso Cubano control room (Alset-JS) http://127.0.0.1%s/\n", addr)
		fmt.Printf("API /api/v1/pulse  POST /api/v1/genes/run\n")
		fmt.Printf("Tokens optional: ELTOQUE_API_TOKEN, QVAPAY_TOKEN\n")
		if err := http.ListenAndServe(addr, srv.Handler()); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown: %s\n", os.Args[1])
		os.Exit(1)
	}
}
