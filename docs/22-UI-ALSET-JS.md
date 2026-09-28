# UI = Alset-JS-Runtime only

Pulso Cubano control room lives in `ui/` and imports:

```js
from '../../Alset-JS-Runtime/src/core/AlsetPulseCore.js'
```

Primitives used: `AlsetInspector`, `alsetState`, `Column`, `Row`, `Text`, `Theme`, `mod`.

**No React / Vue / Svelte for product UI.** Static file server from `pulso serve`.

Data: `GET /api/v1/pulse`, `POST /api/v1/genes/run` → PrismaTec organisms + Mind.
