/**
 * PULSO CUBANO — Economic Intelligence Terminal
 * Alset-JS-Runtime only. Honest empty/partial/live states.
 */
import {
  AlsetInspector, Column, Row, Text, mod, alsetState, Theme
} from './vendor/AlsetPulseCore.js';

Theme.set({
  primary: '#E8C547',
  secondary: '#5B8DEF',
  background: '#0B0E14',
  surface: '#141A24',
  radius: 10
});

const GOLD = '#E8C547';
const BLUE = '#5B8DEF';
const GREEN = '#3DDC97';
const RED = '#F07178';
const MUTED = '#8B93A7';
const LINE = 'rgba(232,197,71,0.12)';
const SURFACE = '#141A24';
const BG = '#0B0E14';

const term = alsetState(null);
const section = alsetState('overview');
const status = alsetState('idle');
const errMsg = alsetState('');
const searchQ = alsetState('');
const selected = alsetState(null);

async function api(path, opts) {
  const r = await fetch(path, opts);
  if (!r.ok) throw new Error(await r.text());
  return r.json();
}

async function load() {
  try {
    status.set('loading');
    term.set(await api('/api/v1/terminal'));
    status.set('ok');
    errMsg.set('');
  } catch (e) {
    status.set('error');
    errMsg.set(String(e.message || e));
  }
}

async function runSensors() {
  try {
    status.set('sensing');
    await api('/api/v1/genes/run', { method: 'POST' });
    await load();
  } catch (e) {
    status.set('error');
    errMsg.set(String(e.message || e));
  }
}

function card(fn) {
  return Column(mod().padding(16).gap(8)
    .addStyle('background', SURFACE)
    .addStyle('border', '1px solid ' + LINE)
    .addStyle('borderRadius', '12px'), fn);
}

function muted(t) { return Text(t, mod().sizeText(12).color(MUTED).addStyle('lineHeight', '1.45')); }
function title(t) { return Text(t, mod().sizeText(15).weight('700').color('#F4F6FA')); }
function h2(t) {
  return Text(t, mod().sizeText(12).weight('700').color(GOLD)
    .addStyle('letterSpacing', '0.1em').addStyle('textTransform', 'uppercase').margin('0 0 10px 0'));
}
function emptyBox(msg) {
  return card(() => {
    Text('Sin datos capturados', mod().sizeText(14).weight('600').color(MUTED));
    muted(msg || 'Ejecute sensores o conecte fuentes Nivel A/B. Esta sección no inventa actividad de mercado.');
  });
}
function stateBadge(st) {
  const c = st === 'live' || st === 'connected' || st === 'ready' ? GREEN
    : st === 'partial' || st === 'pending' || st === 'observed' ? GOLD
    : st === 'error' ? RED : MUTED;
  return Text(String(st || '—'), mod().sizeText(11).weight('600').color(c)
    .padding('3px 8px').addStyle('border', '1px solid ' + c).addStyle('borderRadius', '999px'));
}

function MetricTile(m) {
  const c = m.state === 'live' ? GREEN : m.state === 'delayed' ? GOLD : MUTED;
  return Column(mod().padding(14).gap(4)
    .addStyle('background', SURFACE)
    .addStyle('border', '1px solid ' + LINE)
    .addStyle('borderRadius', '12px')
    .addStyle('minWidth', '150px')
    .addStyle('flex', '1'), () => {
    Text(m.label, mod().sizeText(11).color(MUTED).addStyle('letterSpacing', '0.06em'));
    Text(m.value + (m.unit ? ' ' + m.unit : ''), mod().sizeText(26).weight('800').color(c).margin('4px 0'));
    muted(m.hint || '');
  });
}

function NavBtn(id, label) {
  const on = section.get() === id;
  return Column(mod().padding('8px 10px').margin('0 0 2px 0')
    .addStyle('borderRadius', '8px')
    .addStyle('cursor', 'pointer')
    .addStyle('background', on ? 'rgba(232,197,71,0.12)' : 'transparent')
    .addStyle('borderLeft', on ? '2px solid ' + GOLD : '2px solid transparent')
    .clickable(() => { section.set(id); selected.set(null); }), () => {
    Text(label, mod().sizeText(12).weight(on ? '700' : '500').color(on ? GOLD : '#C5CAD6'));
  });
}

function ObsRow(o) {
  const on = selected.get() && selected.get().id === o.id;
  return Column(mod().padding(12).margin('0 0 6px 0').gap(4)
    .addStyle('background', on ? 'rgba(91,141,239,0.12)' : SURFACE)
    .addStyle('border', '1px solid ' + (on ? BLUE : LINE))
    .addStyle('borderRadius', '10px')
    .addStyle('cursor', 'pointer')
    .clickable(() => selected.set(o)), () => {
    Text(o.text || (o.product + ' · ' + o.source), mod().sizeText(13).weight('600').color('#EEF1F6'));
    muted([o.source, o.type, o.category, o.at].filter(Boolean).join(' · '));
    if (o.price != null) Text(o.price + ' ' + (o.currency || ''), mod().sizeText(13).color(GOLD).weight('700'));
  });
}

function SignalRow(s) {
  return Column(mod().padding(12).margin('0 0 6px 0').gap(4)
    .addStyle('background', SURFACE).addStyle('border', '1px solid ' + LINE)
    .addStyle('borderRadius', '10px').addStyle('cursor', 'pointer')
    .clickable(() => selected.set(s)), () => {
    Row(mod().align('center', 'space-between'), () => {
      Text(s.title, mod().sizeText(13).weight('600').color('#EEF1F6').addStyle('flex', '1'));
      stateBadge(s.state);
    });
    muted(s.evidence + ' · ' + (s.at || ''));
  });
}

function EvidencePanel() {
  const s = selected.get();
  if (!s) {
    return card(() => {
      h2('Evidencia');
      muted('Seleccione una observación o señal para ver procedencia, confianza y límites de interpretación.');
    });
  }
  return card(() => {
    h2('Evidencia');
    title(s.title || s.text || s.id || 'Ítem');
    muted('ID: ' + (s.id || '—'));
    if (s.source) muted('Fuente: ' + s.source);
    if (s.kind || s.type) muted('Tipo: ' + (s.kind || s.type));
    if (s.layer) muted('Capa: ' + s.layer);
    if (s.confidence != null) muted('Confianza: ' + Math.round(s.confidence * 100) + '%');
    if (s.evidence) muted(s.evidence);
    if (s.url) muted('URL: ' + s.url);
    if (s.at) muted('Fecha: ' + s.at);
    Text('Una coincidencia temporal no implica causalidad. Un precio anunciado no es precio de mercado universal.',
      mod().sizeText(11).color(GOLD).margin('8px 0 0 0').addStyle('lineHeight', '1.4'));
  });
}

function Scaffold(name, desc, forWhom) {
  return Column(mod().gap(12), () => {
    h2(name);
    card(() => {
      title('Estructura de producto lista · datos pendientes');
      muted(desc);
      muted('Audiencia: ' + forWhom);
      muted('Estado: scaffold — navegable, sin cifras inventadas. Se activará al conectar fuentes Nivel A/B y flujo de normalización.');
    });
    emptyBox('Cuando existan observaciones normalizadas para esta sección, aparecerán aquí con enlace a evidencia.');
  });
}

/* ——— Sections ——— */
function Overview(t) {
  return Column(mod().gap(14), () => {
    h2('01 · Market Overview');
    muted('Estado de la economía observable en el territorio activo. Explore datos reales o vea cobertura vacía — no marketing.');
    Row(mod().gap(10).addStyle('flexWrap', 'wrap'), () => {
      for (const m of (t.header_metrics || [])) MetricTile(m);
    });
    card(() => {
      title('Modo de terminal: ' + (t.mode || '—'));
      muted(t.disclaimer || '');
      muted('Corte: ' + (t.cutoff || '—') + ' · ' + (t.territory || '') + ' · ' + (t.period || ''));
    });
    h2('Radar rápido');
    const sigs = t.signals || [];
    if (!sigs.length) emptyBox('Sin señales nuevas. Ejecute «Actualizar sensores» o conecte fuentes.');
    else for (const s of sigs.slice(0, 6)) SignalRow(s);
    h2('Salud de fuentes');
    for (const sh of (t.source_health || []).slice(0, 8)) {
      Row(mod().gap(10).align('center').margin('0 0 6px 0').padding(10)
        .addStyle('background', SURFACE).addStyle('borderRadius', '8px')
        .addStyle('border', '1px solid ' + LINE), () => {
        stateBadge(sh.status);
        Column(mod().gap(2).addStyle('flex', '1'), () => {
          Text(sh.name, mod().sizeText(13).weight('600').color('#EEF1F6'));
          muted(sh.note);
        });
      });
    }
  });
}

function Radar(t) {
  return Column(mod().gap(12), () => {
    h2('02 · Market Radar');
    muted('Cambios detectados que merecen investigación. Cada ítem abre evidencia a la derecha.');
    const sigs = t.signals || [];
    if (!sigs.length) emptyBox('No hay señales observadas en esta corrida.');
    else for (const s of sigs) SignalRow(s);
  });
}

function Explorer(t) {
  return Column(mod().gap(12), () => {
    h2('03 · Market Explorer');
    muted('Consulta multidimensional: territorio × sector × período × tipo de observación.');
    card(() => {
      title('Contexto activo');
      muted('Territorio: ' + (t.territory || 'Cuba / Guantánamo'));
      muted('Mercado: ' + (t.market || 'General'));
      muted('Período: ' + (t.period || '—'));
      muted('Filtros persistentes y árbol geo se activan con catálogo normalizado.');
    });
    const q = (searchQ.get() || '').toLowerCase();
    const obs = (t.observations || []).filter(o => {
      if (!q) return true;
      return JSON.stringify(o).toLowerCase().includes(q);
    });
    h2('Resultados de observación');
    if (!obs.length) emptyBox('Sin observaciones para el filtro actual.');
    else for (const o of obs.slice(0, 30)) ObsRow(o);
  });
}

function Prices(t) {
  return Column(mod().gap(12), () => {
    h2('04 · Price Intelligence');
    muted('Precios anunciados — no precio de mercado universal. Muestra, moneda y antigüedad visibles.');
    const list = t.prices || [];
    if (!list.length) emptyBox('Aún no hay muestras de precio capturadas (conecte El Toque, clasificados o CSV).');
    else for (const o of list) ObsRow(o);
  });
}

function News(t) {
  return Column(mod().gap(12), () => {
    h2('09 · News & Events');
    muted('Contexto económico con fuente. Coincidencia temporal ≠ causalidad.');
    const list = t.news || [];
    if (!list.length) emptyBox('Sin titulares capturados. El gene RSS/Wiki alimenta esta lista cuando hay red.');
    else for (const o of list) ObsRow(o);
  });
}

function Methodology(t) {
  return Column(mod().gap(12), () => {
    h2('14 · Data & Methodology');
    card(() => {
      title('Capas de datos');
      muted('official · observed_market · reference · first_party · user_provided');
      muted('Nunca se mezclan sin etiqueta. Un anuncio no es una venta.');
    });
    card(() => {
      title('Prioridad de fuentes');
      muted('A — API/datos oficiales y propios con consentimiento');
      muted('B — Web pública indexable con ToS');
      muted('C — Plataformas restringidas solo con autorización');
      muted('D — Investigación primaria y validación de campo');
    });
    h2('Registro de fuentes');
    for (const sh of (t.source_health || [])) {
      card(() => {
        Row(mod().gap(8).align('center'), () => {
          stateBadge(sh.status);
          Text(sh.name, mod().sizeText(13).weight('600'));
        });
        muted(sh.id + ' · ' + sh.type + ' · access=' + sh.access);
        muted(sh.note);
      });
    }
  });
}

function APISection() {
  return Column(mod().gap(12), () => {
    h2('15 · API & Data Products');
    card(() => {
      title('Endpoints actuales');
      muted('GET /api/v1/health');
      muted('GET /api/v1/terminal  — estado de la terminal');
      muted('GET /api/v1/dashboard — brief legacy');
      muted('POST /api/v1/genes/run — corrida de sensores');
      muted('GET /api/v1/sources · /api/v1/pulse · /api/v1/queries/gtmo');
    });
    muted('Productos comerciales (CSV, series, licencias) requieren cobertura estable y derechos claros — no se publican dumps sin licencia.');
  });
}

function SectionBody(t) {
  const id = section.get();
  if (id === 'overview') return Overview(t);
  if (id === 'radar') return Radar(t);
  if (id === 'explorer') return Explorer(t);
  if (id === 'prices') return Prices(t);
  if (id === 'news') return News(t);
  if (id === 'methodology') return Methodology(t);
  if (id === 'api') return APISection();
  const meta = (t.sections || []).find(s => s.id === id);
  return Scaffold(
    (meta && meta.title) || id,
    (meta && meta.description) || 'Sección planificada de la terminal.',
    (meta && meta.audience) || 'Profesionales'
  );
}

AlsetInspector(() => {
  const t = term.get() || { nav: [], header_metrics: [], source_health: [], signals: [], observations: [] };

  Column(mod().addStyle('minHeight', '100dvh').addStyle('background', BG), () => {
    // Shell header
    Row(mod().padding('12px 18px').align('center', 'space-between')
      .addStyle('borderBottom', '1px solid ' + LINE)
      .addStyle('background', '#0F131A')
      .addStyle('flexWrap', 'wrap').gap(12), () => {
      Column(mod().gap(2), () => {
        Text('PULSO CUBANO', mod().sizeText(16).weight('900').color(GOLD).addStyle('letterSpacing', '0.08em'));
        Text('ECONOMIC INTELLIGENCE TERMINAL', mod().sizeText(11).color(MUTED).addStyle('letterSpacing', '0.14em'));
      });
      Column(mod().gap(2).addStyle('flex', '1').addStyle('minWidth', '200px'), () => {
        Text('LIVE DATA · ' + (t.territory || 'CUBA / GUANTÁNAMO') + ' · ' + (t.market || 'MERCADO GENERAL'),
          mod().sizeText(11).color(BLUE).weight('600'));
        muted('Corte ' + (t.cutoff || '—') + ' · modo ' + (t.mode || '—'));
      });
      Row(mod().gap(8).align('center'), () => {
        stateBadge(status.get());
        Column(mod().padding('8px 12px').addStyle('background', GOLD).addStyle('borderRadius', '8px')
          .addStyle('cursor', 'pointer').clickable(() => runSensors()), () => {
          Text('Actualizar sensores', mod().sizeText(12).weight('700').color('#111'));
        });
        Column(mod().padding('8px 12px').addStyle('border', '1px solid ' + BLUE).addStyle('borderRadius', '8px')
          .addStyle('cursor', 'pointer').clickable(() => load()), () => {
          Text('Refrescar', mod().sizeText(12).weight('600').color(BLUE));
        });
      });
    });

    if (errMsg.get()) {
      Column(mod().padding(12), () => {
        Text(errMsg.get(), mod().sizeText(13).color(RED));
      });
    }

    // Body: nav | workspace | evidence
    Row(mod().align('start', 'start').addStyle('minHeight', 'calc(100dvh - 64px)'), () => {
      // Nav
      Column(mod().padding(12).gap(2)
        .addStyle('width', '200px')
        .addStyle('minWidth', '180px')
        .addStyle('borderRight', '1px solid ' + LINE)
        .addStyle('background', '#0F131A'), () => {
        muted('NAVEGACIÓN');
        for (const n of (t.nav || [])) NavBtn(n.id, n.label);
      });

      // Workspace
      Column(mod().padding(18).gap(12).addStyle('flex', '1').addStyle('minWidth', '0')
        .addStyle('maxHeight', 'calc(100dvh - 64px)').addStyle('overflow', 'auto'), () => {
        // Universal search (client filter on loaded obs)
        Column(mod().padding(10).margin('0 0 8px 0')
          .addStyle('background', SURFACE).addStyle('border', '1px solid ' + LINE)
          .addStyle('borderRadius', '10px'), () => {
          Text('Búsqueda en observaciones cargadas', mod().sizeText(11).color(MUTED));
          // simple prompt via repeated clicks not ideal — use text display of query tip
          muted('Use Explorer y el radar; filtro textual ampliado en siguiente iteración de Input Alset.');
          muted('Consulta ejemplo: equipos informáticos · Guantánamo · 90 días · ofertas vs solicitudes');
        });
        SectionBody(t);
      });

      // Detail
      Column(mod().padding(12)
        .addStyle('width', '280px')
        .addStyle('minWidth', '240px')
        .addStyle('borderLeft', '1px solid ' + LINE)
        .addStyle('background', '#0F131A')
        .addStyle('maxHeight', 'calc(100dvh - 64px)')
        .addStyle('overflow', 'auto'), () => {
        EvidencePanel();
      });
    });
  });
});

load();
