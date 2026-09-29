/**
 * PULSO CUBANO — Investor & Business Control Room (Alset-JS only)
 */
import {
  AlsetInspector, Column, Row, Text, mod, alsetState, Theme
} from './vendor/AlsetPulseCore.js';

Theme.set({
  primary: '#F4B400',
  secondary: '#3dd68c',
  background: '#07090e',
  surface: '#12161f',
  radius: 14
});

const GOLD = '#F4B400';
const CYAN = '#4cc9f0';
const GREEN = '#3dd68c';
const RED = '#ff6b6b';
const MUTED = '#8b93a7';
const CARD = 'rgba(255,255,255,0.04)';
const LINE = 'rgba(244,180,0,0.14)';

const dash = alsetState(null);
const status = alsetState('idle');
const errMsg = alsetState('');
const tab = alsetState('overview');

async function api(path, opts) {
  const r = await fetch(path, opts);
  if (!r.ok) throw new Error(await r.text());
  return r.json();
}

async function loadDashboard() {
  try {
    status.set('loading');
    dash.set(await api('/api/v1/dashboard'));
    status.set('ok');
    errMsg.set('');
  } catch (e) {
    status.set('error');
    errMsg.set(String(e.message || e));
  }
}

async function senseMarket() {
  try {
    status.set('sensing');
    await api('/api/v1/genes/run', { method: 'POST' });
    await loadDashboard();
  } catch (e) {
    errMsg.set(String(e.message || e));
    status.set('error');
  }
}

function SectionTitle(title, sub) {
  return Column(mod().gap(4).margin('0 0 12px 0'), () => {
    Text(title, mod().sizeText(13).weight('700').color(GOLD)
      .addStyle('letterSpacing', '0.12em').addStyle('textTransform', 'uppercase'));
    if (sub) Text(sub, mod().sizeText(13).color(MUTED));
  });
}

function Card(fn) {
  return Column(mod().padding(18).gap(10)
    .addStyle('background', CARD)
    .addStyle('border', '1px solid ' + LINE)
    .addStyle('borderRadius', '16px')
    .addStyle('flex', '1')
    .addStyle('minWidth', '200px'), fn);
}

function Bar(label, value, color) {
  const v = Math.max(0, Math.min(100, Number(value) || 0));
  return Column(mod().gap(6).margin('0 0 10px 0'), () => {
    Row(mod().align('center', 'space-between'), () => {
      Text(label, mod().sizeText(13).color('#d0d4de'));
      Text(String(Math.round(v)), mod().sizeText(13).weight('700').color(color || GOLD));
    });
    Column(mod().height(8).addStyle('background', 'rgba(255,255,255,0.06)')
      .addStyle('borderRadius', '99px').addStyle('overflow', 'hidden'), () => {
      Column(mod().height(8).addStyle('width', v + '%')
        .addStyle('background', color || GOLD).addStyle('borderRadius', '99px'), () => {});
    });
  });
}

function Pill(label, key) {
  const on = tab.get() === key;
  return Column(mod().padding('8px 14px')
    .addStyle('borderRadius', '999px')
    .addStyle('cursor', 'pointer')
    .addStyle('background', on ? 'rgba(244,180,0,0.18)' : 'transparent')
    .addStyle('border', on ? '1px solid ' + GOLD : '1px solid rgba(255,255,255,0.1)')
    .clickable(() => tab.set(key)), () => {
    Text(label, mod().sizeText(12).weight('600').color(on ? GOLD : MUTED));
  });
}

function BtnPrimary(label, fn) {
  return Column(mod().padding('12px 18px').addStyle('background', GOLD)
    .addStyle('borderRadius', '12px').addStyle('cursor', 'pointer').clickable(fn), () => {
    Text(label, mod().sizeText(13).weight('700').color('#111'));
  });
}

function BtnGhost(label, fn) {
  return Column(mod().padding('12px 18px').addStyle('border', '1px solid ' + CYAN)
    .addStyle('borderRadius', '12px').addStyle('cursor', 'pointer').clickable(fn), () => {
    Text(label, mod().sizeText(13).weight('600').color(CYAN));
  });
}

function Overview(d) {
  const p = d.pulse || {};
  return Column(mod().gap(18), () => {
    SectionTitle('Resumen para decisión', 'Qué está pasando y por qué importa a negocio e inversor');
    Card(() => {
      Text(d.headline || '—', mod().sizeText(17).weight('600').color('#f2f4f8').addStyle('lineHeight', '1.45'));
      Text(d.disclaimer || '', mod().sizeText(12).color(MUTED).margin('10px 0 0 0').addStyle('lineHeight', '1.4'));
    });
    Row(mod().gap(12).addStyle('flexWrap', 'wrap'), () => {
      Card(() => {
        Text('PULSO GENERAL', mod().sizeText(11).color(MUTED).addStyle('letterSpacing', '0.08em'));
        Text(p.overall != null ? String(Math.round(p.overall)) : '—',
          mod().sizeText(42).weight('800').color(GOLD).margin('6px 0'));
        Text('Método ' + (p.method_version || 'pulse-v0') + ' · ' + (p.layer || 'observed_market'),
          mod().sizeText(12).color(MUTED));
      });
      Card(() => {
        Text('COMPONENTES', mod().sizeText(11).color(MUTED));
        Bar('Demanda', p.demand, GREEN);
        Bar('Oferta', p.supply, CYAN);
        Bar('Precio / tensión', p.price, GOLD);
        Bar('Evidencia', p.evidence, '#a78bfa');
      });
      Card(() => {
        const inf = d.inflation || {};
        Text('INFLACIÓN (PROXY)', mod().sizeText(11).color(MUTED));
        Text(inf.score != null ? String(Math.round(inf.score)) : '—',
          mod().sizeText(42).weight('800').color(RED));
        Text(inf.proxy_label || '', mod().sizeText(12).color(MUTED));
        Text(inf.explanation || '', mod().sizeText(12).color('#b8bfcc').margin('8px 0 0 0').addStyle('lineHeight', '1.4'));
      });
    });
    SectionTitle('Categorías · Guantánamo', 'Demanda vs oferta observable — la brecha orienta oportunidad');
    Row(mod().gap(12).addStyle('flexWrap', 'wrap'), () => {
      for (const c of (d.categories || [])) {
        Card(() => {
          Text(c.name, mod().sizeText(15).weight('700').color('#fff'));
          Text(c.momentum || '', mod().sizeText(12).color(GREEN));
          Bar('Demanda', c.demand, GREEN);
          Bar('Oferta', c.supply, CYAN);
          Text('Brecha ' + Math.round(c.gap) + ' pts', mod().sizeText(13).weight('700').color(GOLD));
          Text(c.evidence || '', mod().sizeText(11).color(MUTED));
        });
      }
    });
    SectionTitle('Lecturas para inversores', 'Hipótesis de trabajo — no consejo de inversión automático');
    for (const t of (d.takeaways || [])) {
      Card(() => { Text('→  ' + t, mod().sizeText(14).color('#e2e6ef').addStyle('lineHeight', '1.45')); });
    }
  });
}

function FXView(d) {
  return Column(mod().gap(16), () => {
    SectionTitle('Mercado de divisas', 'Oficial ≠ informal ≠ global');
    for (const x of (d.fx || [])) {
      Card(() => {
        Row(mod().align('center', 'space-between').addStyle('flexWrap', 'wrap').gap(8), () => {
          Text(x.pair, mod().sizeText(16).weight('700'));
          Text((x.value ? x.value : '—') + (x.unit ? ' ' + x.unit : ''),
            mod().sizeText(22).weight('800').color(GOLD));
        });
        Text('Fuente: ' + (x.source || '—') + ' · confianza ' + Math.round((x.confidence || 0) * 100) + '%',
          mod().sizeText(12).color(MUTED));
        Text(x.note || '', mod().sizeText(13).color('#c5cad6').margin('6px 0 0 0'));
      });
    }
    Card(() => {
      Text('Cómo usarlo en negocio', mod().sizeText(14).weight('700').color(CYAN));
      Text('Costeo de importación con escenarios · pricing local sin anclarse a un solo print · el spread es fricción, no arbitraje garantizado.',
        mod().sizeText(13).color('#c5cad6').addStyle('lineHeight', '1.5'));
    });
  });
}

function DemandView(d) {
  return Column(mod().gap(16), () => {
    SectionTitle('Qué se necesita vs qué se ofrece', 'Presión observable — no ranking oficial de ventas');
    Row(mod().gap(12).addStyle('flexWrap', 'wrap'), () => {
      Card(() => {
        Text('MÁS NECESITADOS (demanda)', mod().sizeText(12).color(GREEN).weight('700'));
        for (const p of (d.demand_top || [])) {
          Column(mod().margin('10px 0 0 0'), () => {
            Text(p.product, mod().sizeText(15).weight('600'));
            Bar(p.metric, p.score, GREEN);
            Text(p.evidence, mod().sizeText(11).color(MUTED));
          });
        }
      });
      Card(() => {
        Text('OFERTA MÁS VISIBLE', mod().sizeText(12).color(CYAN).weight('700'));
        for (const p of (d.supply_top || [])) {
          Column(mod().margin('10px 0 0 0'), () => {
            Text(p.product, mod().sizeText(15).weight('600'));
            Bar(p.metric, p.score, CYAN);
            Text(p.evidence, mod().sizeText(11).color(MUTED));
          });
        }
      });
    });
    SectionTitle('Oportunidades', 'Con evidencia y confianza');
    for (const o of (d.opportunities || [])) {
      Card(() => {
        Text(o.title, mod().sizeText(16).weight('700').color(GOLD));
        Text(o.why, mod().sizeText(14).color('#e2e6ef').margin('6px 0'));
        Text('Audiencia: ' + o.audience, mod().sizeText(12).color(CYAN));
        Text('Evidencia: ' + o.evidence + ' · confianza ' + Math.round((o.confidence || 0) * 100) + '%',
          mod().sizeText(12).color(MUTED));
      });
    }
  });
}

function CampaignsView(d) {
  return Column(mod().gap(16), () => {
    SectionTitle('Constructor de campañas', 'De la señal al borrador comercial');
    Text('No publica ads solo. Arma el argumento a partir de brechas observadas.', mod().sizeText(13).color(MUTED));
    for (const c of (d.campaigns || [])) {
      Card(() => {
        Text(c.name, mod().sizeText(17).weight('800').color('#fff'));
        Text('Para: ' + c.audience, mod().sizeText(13).color(CYAN).margin('6px 0'));
        Text('"' + c.message + '"', mod().sizeText(15).color(GOLD).addStyle('fontStyle', 'italic').addStyle('lineHeight', '1.4'));
        Text('Canal: ' + c.channel, mod().sizeText(13).color('#c5cad6').margin('8px 0 0 0'));
        Text('Basado en: ' + c.based_on, mod().sizeText(12).color(MUTED));
      });
    }
  });
}

function MethodView(d) {
  const e = d.evidence || {};
  return Column(mod().gap(16), () => {
    SectionTitle('Metodología y evidencia', 'Límites claros = confianza real');
    Card(() => {
      Text('Observaciones (última corrida): ' + (e.observation_n || 0), mod().sizeText(15).weight('600'));
      Text('Fuentes activas: ' + (e.sources_active || 0) + ' / ' + (e.sources_total || 0), mod().sizeText(14).color(MUTED));
      Text('Método: ' + (e.method || ''), mod().sizeText(13).color('#c5cad6').margin('8px 0'));
      Text(e.coverage_note || '', mod().sizeText(13).color(MUTED).addStyle('lineHeight', '1.45'));
    });
    for (const g of (d.genes_last || [])) {
      Card(() => {
        Text((g.name || '') + ' · obs=' + (g.observations || 0) + ' · mind=' + (g.mind_selected || '—'),
          mod().sizeText(13).weight('600'));
        const errs = g.errors || [];
        if (errs.length) Text(errs.join(' · '), mod().sizeText(12).color(RED));
      });
    }
  });
}

AlsetInspector(() => {
  const d = dash.get() || {};
  const key = tab.get();

  Column(mod().padding(22).gap(18).addStyle('maxWidth', '1180px').addStyle('margin', '0 auto')
    .addStyle('minHeight', '100dvh'), () => {

    Row(mod().gap(12).align('center', 'space-between').addStyle('flexWrap', 'wrap'), () => {
      Column(mod().gap(4), () => {
        Text('PULSO CUBANO', mod().sizeText(24).weight('900').color(GOLD).addStyle('letterSpacing', '0.04em'));
        Text((d.market || 'Cuba') + ' · inteligencia de mercado para negocios e inversores',
          mod().sizeText(13).color(MUTED));
      });
      Row(mod().gap(10).align('center').addStyle('flexWrap', 'wrap'), () => {
        Text(status.get(), mod().sizeText(12).color(status.get() === 'error' ? RED : GREEN)
          .padding('6px 12px').addStyle('border', '1px solid rgba(255,255,255,0.12)')
          .addStyle('borderRadius', '999px'));
        BtnPrimary('Actualizar sensores', () => senseMarket());
        BtnGhost('Refrescar panel', () => loadDashboard());
      });
    });

    if (errMsg.get()) Text(errMsg.get(), mod().sizeText(13).color(RED));

    Row(mod().gap(8).addStyle('flexWrap', 'wrap'), () => {
      Pill('Resumen', 'overview');
      Pill('Divisas', 'fx');
      Pill('Demanda & oferta', 'demand');
      Pill('Campañas', 'campaigns');
      Pill('Método', 'method');
    });

    if (key === 'overview') Overview(d);
    else if (key === 'fx') FXView(d);
    else if (key === 'demand') DemandView(d);
    else if (key === 'campaigns') CampaignsView(d);
    else MethodView(d);

    Text('UI Alset-JS · PrismaTec-Core · observed_market · no inventamos ventas oficiales',
      mod().sizeText(11).color(MUTED).margin('28px 0 8px 0'));
  });
});

loadDashboard();
