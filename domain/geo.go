package domain

// Geography hierarchy: country → province → municipality → locality.
// GTMO is the first configured market; architecture is not Cuba-hardcoded beyond seed data.

type Geo struct {
	Country      string `json:"country"`
	Province     string `json:"province,omitempty"`
	Municipality string `json:"municipality,omitempty"`
	Locality     string `json:"locality,omitempty"`
}

const CountryCuba = "CU"
const ProvinceGuantanamo = "Guantánamo"

var GuantanamoMunicipalities = []string{
	"Guantánamo", "Baracoa", "El Salvador", "Manuel Tames",
	"Niceto Pérez", "Caimanera", "Yateras", "Imías", "Maisí",
}
