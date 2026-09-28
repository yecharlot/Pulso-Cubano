# Universo de fuentes — sensor de mercado (SaaS)

Pulso Cubano debe **sensar el máximo de información legítima**.  
El Toque es **una** referencia de FX informal, no la única fuente.

## Capas de datos

| Capa | Significado |
|------|-------------|
| `official` | Banco Central, ONEI, Gaceta, ministerios |
| `observed_market` | Clasificados, búsqueda, P2P, tasas referenciales |
| `reference` | Wikipedia, documentación, taxonomías |
| `first_party` | Datos que un negocio entrega voluntariamente |
| `user_provided` | CSV/export que sube el cliente |

## Fuentes objetivo (ampliación)

### Oficiales / institucionales
- Banco Central de Cuba (tipos oficiales, comunicados)
- ONEI / estadísticas públicas
- Gaceta Oficial (regulatorio)
- Ministerios con portales abiertos (turismo, comercio, etc.)

### FX e informal referencial
- El Toque TRMI API (token, ToS)
- Otras series públicas de tasas cuando existan APIs formales

### Digital / fintech observada
- QvaPay (API autenticada, mínimos privilegios)
- Otras plataformas con API pública documentada

### Clasificados y web pública
- Revolico, Cubisima y similares (**solo HTML/API público**, robots.txt, sin cuentas robadas)
- Resultados de búsqueda web (Google Programmable Search / CSE con API key propia)

### Tendencias e interés
- Google Trends (vía fuentes legales / datasets / API oficial cuando aplique)
- Google News / RSS de medios económicos cubanos e internacionales sobre Cuba

### Conocimiento estructurado
- Wikipedia / Wikidata (categorías, entidades, geografía)

### Mercado internacional (contexto)
- APIs gratuitas de FX global (p. ej. open.er-api, Frankfurter) para **comparar** CUP vs USD/EUR — no sustituyen mercado cubano
- Commodity/open data cuando aporte contexto (alimentos, energía)

### First-party y SaaS
- Onboarding de MIPYMES: subida de precios, inventario, ventas agregadas
- Webhooks de clientes enterprise

## Lo que NO hacemos

- Acceso clandestino a **grupos privados** de Facebook, WhatsApp o Telegram  
- Credenciales compartidas, cookies robadas, bypass de login  
- Presentar un anuncio o un mensaje de grupo como “venta confirmada”  
- Mezclar tasa BNC y tasa informal sin etiqueta  

### WhatsApp / Facebook / Telegram

| Escenario | Política Pulso |
|-----------|----------------|
| Página o canal **público** con URL abierta | Evaluable (respetando ToS y rate limits) |
| Grupo **privado** | Solo si el **titular** exporta datos o conecta first-party |
| Scraping masivo de grupos cerrados | **Prohibido** en el diseño del producto |

Para SaaS: el valor está en **agregación + histórico + metodología**, no en invadir privacidad.

## Genes previstos (roadmap)

`gene.fx.official` · `gene.fx.informal` · `gene.search.web` · `gene.classifieds.*` · `gene.trends` · `gene.wiki.ref` · `gene.news.rss` · `gene.firstparty` · `gene.manual`

Cada gene escribe `Observation` con `source_id`, `confidence`, `coverage`, `provenance.layer`.
