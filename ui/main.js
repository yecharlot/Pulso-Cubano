/**
 * PULSO CUBANO — Control Room
 * UI exclusively on Alset-JS-Runtime (AlsetInspector + alsetState + Column/Row/Text).
 */
import {
  AlsetInspector, Column, Row, Text, mod, alsetState, Theme
} from '../../Alset-JS-Runtime/src/core/AlsetPulseCore.js';

Theme.set({
  primary: '#F4B400',
  secondary: '#00ACC1',
  background: '#0a0c10',
  surface: '#141820',
  radius: 16
});

const GOLD = '#F4B400';
const CYAN = '#00ACC1';
const MUTED = '#8b93a7';
const OK = '#3dd68c';
const WARN = '#f5a524';

const pulse = alsetState(null);
const genes = alsetState([]);
const sources = alsetState([]);
const status = alsetState('idle');
const errMsg = alsetState('');

async function api(path, opts) {
  const r = await fetch(path, opts);
  if (!r.ok) throw new Error(await r.text());
  return r.json();
}

async function refresh() {
  try {
    status.set('loading');
    const [p, s, g] = await Promise.all([
      api('/api/v1/pulse'),
      api('/api/v1/sources'),
      api('/api/v1/genes/last').catch(() => [])
    ]);
    pulse.set(p);
    sources.set(s);
    genes.set(Array.isArray(g) ? g : []);
    status.set('ok');
    errMsg.set('');
  } catch (e) {
    status.set('error');
    errMsg.set(String(e.message || e));
  }
}

async function runGenes() {
  try {
    status.set('sensing');
    const res = await api('/api/v1/genes/run', { method: 'POST' });
    genes.set(res);
    await refresh();
  } catch (e) {
    errMsg.set(String(e.message || e));
    status.set('error');
  }
}

function MetricCard(label, value, hint) {
  return Column(mod()
    .padding(16)
    .addStyle('background', 'rgba(255,255,255,0.04)')
    .addStyle('border', '1px solid rgba(244,180,0,0.12)')
    .addStyle('borderRadius', '16px')
    .addStyle('minWidth', '120px')
    .addStyle('flex', '1'), () => {
    Text(label, mod().sizeText(11).color(MUTED).addStyle('letterSpacing', '0.08em').addStyle('textTransform', 'uppercase'));
    Text(value == null || value === '' ? '—' : String(value), mod().sizeText(28).weight('700').color(GOLD).margin('8px 0 4px 0'));
    if (hint) Text(hint, mod().sizeText(12).color(MUTED));
  });
}

AlsetInspector(() => {
  Column(mod().padding(24).gap(20).addStyle('maxWidth', '1100px').addStyle('margin', '0 auto')
    .addStyle('minHeight', '100dvh').addStyle('background', '#0a0c10'), () => {
    Row(mod().gap(12).align('center', 'space-between').addStyle('flexWrap', 'wrap'), () => {
      Column(mod().gap(4), () => {
        Text('PULSO CUBANO', mod().sizeText(22).weight('800').color(GOLD).addStyle('letterSpacing', '0.06em'));
        Text('Market Intelligence · observed market · Alset-JS Runtime', mod().sizeText(13).color(MUTED));
      });
      Text(status.get(), mod().sizeText(12).color(status.get() === 'error' ? WARN : OK).padding('6px 12px')
        .addStyle('border', '1px solid rgba(255,255,255,0.1)').addStyle('borderRadius', '999px'));
    });

    Text('Cuba · Guantánamo first · REAL ≠ OBSERVED. El Toque / QvaPay genes when tokens set.', mod().sizeText(13).color(MUTED));

    const p = pulse.get() || {};
    Row(mod().gap(12).addStyle('flexWrap', 'wrap'), () => {
      MetricCard('Overall Pulse', p.overall != null ? Math.round(p.overall) : '—', p.method_version || 'pulse-v0');
      MetricCard('Demand', p.demand, 'observed');
      MetricCard('Supply', p.supply, 'observed');
      MetricCard('Price', p.price, 'observed');
      MetricCard('Evidence', p.evidence, 'n=' + (p.observation_n || 0));
    });

    Row(mod().gap(12).margin('8px 0'), () => {
      Column(mod().padding('12px 18px').addStyle('background', GOLD).addStyle('borderRadius', '12px')
        .addStyle('cursor', 'pointer').clickable(() => { runGenes(); }), () => {
        Text('RUN GENES (sense market)', mod().sizeText(13).weight('700').color('#111'));
      });
      Column(mod().padding('12px 18px').addStyle('border', '1px solid ' + CYAN)
        .addStyle('borderRadius', '12px').addStyle('cursor', 'pointer').clickable(() => { refresh(); }), () => {
        Text('REFRESH', mod().sizeText(13).weight('600').color(CYAN));
      });
    });

    if (errMsg.get()) Text(errMsg.get(), mod().sizeText(13).color(WARN));

    Text('GENES (last run)', mod().sizeText(12).color(GOLD).addStyle('letterSpacing', '0.1em'));
    const gl = genes.get() || [];
    if (!gl.length) {
      Text('No gene run yet — press RUN GENES.', mod().sizeText(13).color(MUTED));
    }
    for (const g of gl) {
      Column(mod().padding(12).margin('0 0 8px 0')
        .addStyle('background', 'rgba(0,172,193,0.06)')
        .addStyle('border', '1px solid rgba(0,172,193,0.2)')
        .addStyle('borderRadius', '12px'), () => {
        Text((g.name || '') + ' · obs=' + (g.observations ?? 0) + ' · mind=' + (g.mind_selected || '—'),
          mod().sizeText(14).weight('600'));
        const errs = g.errors || [];
        if (errs.length) Text(errs.join(' · '), mod().sizeText(12).color(WARN));
      });
    }

    Text('SOURCES', mod().sizeText(12).color(GOLD).addStyle('letterSpacing', '0.1em').margin('12px 0 0 0'));
    for (const s of (sources.get() || [])) {
      Text((s.id || '') + ' · ' + (s.type || '') + ' · active=' + s.active + ' · ' + (s.access || ''),
        mod().sizeText(13).color(MUTED).margin('4px 0'));
    }

    Text('UI: Alset-JS-Runtime · Motor: PrismaTec-Core · Product: Pulso Cubano',
      mod().sizeText(11).color(MUTED).margin('24px 0 0 0'));
  });
});

refresh();
