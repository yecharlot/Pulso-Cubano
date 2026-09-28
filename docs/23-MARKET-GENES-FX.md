# Market genes & real sensors

## Architecture

```
Scheduler / HTTP
    → gene.fx.eltoque   (ELTOQUE_API_TOKEN)
    → gene.fx.qvapay    (QVAPAY_TOKEN)
    → gene.manual.seed  (demo offline)
         ↓
    Observation → organism memory
         ↓
    Mind.Evaluate → propose record_signal
         ↓
    Pulse emit → Alset-JS control room
```

## El Toque

- Official API: https://tasas.eltoque.com/  
- Token by application form; cite elTOQUE; rates **referential**.  
- Env: `ELTOQUE_API_TOKEN`, optional `ELTOQUE_API_BASE`, `ELTOQUE_API_PATH`.  
- Layer: `observed_market` (informal), never mixed with BNC official without label.

## QvaPay

- API: https://api.qvapay.com — Bearer token, least privilege.  
- P2P averages = **observed** activity, not national accounts.  
- Env: `QVAPAY_TOKEN`.

## Principles

- No private Telegram/WhatsApp scraping.  
- No inventing sales.  
- Genes report errors in UI when tokens missing; seed gene still runs for Mind path demo.
