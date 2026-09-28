# PrismaTec Integration Map — Pulso Cubano

**Rule:** PrismaTec-Core is the motor. Pulso Cubano is the product vertical.  
**Core is never modified for market ontology.**

## Audit summary (PrismaTec-Core ~0.3.0)

| Pulso need | Core capability | Interface / type | Adapter needed? | Notes |
|------------|-----------------|------------------|-----------------|-------|
| Node identity | `core.Node`, `identity.Identity` | `Node.ID()`, `Identity()` | thin | Host for observers |
| Digital organisms | `organism.Manager` | `Create`, `PutMemory`, `Get`, `List` | thin | `GTMOObserver` etc. live in vertical |
| Working / semantic memory | organism memory maps | `PutMemory`, `PutSemantic` | thin | Store observations as JSON blobs (MVP) |
| Events | `events.Bus` | Subscribe / publish | optional | Future market event fan-out |
| Pulse | `pulse.Hub` | `Emit`, `Recent` | thin | Market observation pulses to Studio |
| Provenance | `provenance.Log` | append/query | later | Chain observation → signal |
| Policy | `policy.Engine` | `Evaluate` | thin | Gate campaigns / exports |
| Execution | `runtime/execution` | Engine interface | later | Only authorized actions |
| Inference | `inference.Provider` | `Infer` | optional | Annotations only |
| Mind | `runtime/mind.Engine` | `Observe`, `Evaluate` | thin | Propose signals; Core authorizes |
| Zyrion | `runtime/zyrion` | ternary + RuleSet | thin | Epistemic SUPPORT/UNKNOWN/NEGATED |
| AIP | `api/aip` | HTTP/SSE | later | External market API surface may mirror patterns |
| Network / replication | `network.Transport`, `replication.Service` | TCP/libp2p | later | Multi-node observatory |
| Persistence / CID | organism snapshot + storage/cid | Snapshot/Adopt | later | Market graph history |

## What stays ONLY in Pulso Cubano

- Observation, SourceRegistry, QueryTemplate  
- Market graph, metrics (Demand Pressure…), Pulse Score  
- Opportunity / Campaign engines  
- GTMO geography seed, taxonomies  
- Collectors (Revolico, Trends, …) when legal  

## Adapter package

`adapter/prismatec` is the **only** import path from vertical → Core.

```
Pulso domain → Adapter → Core organisms / pulse / mind
```

## Modification policy for Core

Modify Core **only** if a primitive is proven generic for many verticals (e.g. shared time-series store).  
Market-specific fields never enter Core packages.
