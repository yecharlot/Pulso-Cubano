# Terminal de inteligencia económica

Pulso Cubano se presenta como **Economic Intelligence Terminal**, no como panel de debug.

## Experiencias

1. **Terminal de mercado** — Overview, Radar, Explorer, Prices, News  
2. **Inteligencia económica** — Directory, Network, Sectors, Opportunities  
3. **Inteligencia operativa** — Business Terminal, Campaigns, Reports, API  

## Estados de datos

| Estado | Significado |
|--------|-------------|
| `empty` | Sin observaciones capturadas |
| `partial` | Hay sensores con datos reales (genes públicos) |
| `live` | Métrica respaldada por observaciones en memoria |
| `scaffold` | Sección navegable sin cifras inventadas |

## API

`GET /api/v1/terminal` — estado completo de la terminal.

## Regla

No se muestran indicadores de mercado fabricados. Las secciones 05–08 y 10–13 arrancan en scaffold hasta fuentes Nivel A/B estables.
