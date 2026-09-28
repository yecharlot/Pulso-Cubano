# Pulso Cubano

**Market intelligence vertical** sobre [PrismaTec-Core](https://github.com/yecharlot/PrismaTec-Core).

```
PrismaTec-Core (motor)
        ↓
   Pulso Cubano (producto)
```

Core **no** depende de este repo. Este repo **sí** depende de Core vía `adapter/prismatec`.

## Quick start

```bash
# desde el monorepo local (replace en go.mod → ../PrismaTec-Core)
go test ./...
go run ./cmd/pulso sources
go run ./cmd/pulso queries gtmo
go run ./cmd/pulso pulse-demo
```

## Docs

Ver `docs/` — empezar por `PRISMATEC_INTEGRATION_MAP.md` y `19-BUSINESS_MODEL.md`.

## Principio

No inventar ventas, share ni demanda real. Observatory primero.
