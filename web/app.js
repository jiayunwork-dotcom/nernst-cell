"use strict";

const $ = (id) => document.getElementById(id);
const errorBox = $("errorBox");

function showError(msg) {
  errorBox.hidden = false;
  errorBox.textContent = msg;
}

function clearError() {
  errorBox.hidden = true;
  errorBox.textContent = "";
}

async function postJSON(url, body) {
  clearError();
  const res = await fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  let data = null;
  try {
    data = await res.json();
  } catch (e) {
    data = {};
  }
  if (!res.ok) {
    throw new Error((data && data.error) ? data.error : ("HTTP " + res.status));
  }
  return data;
}

function fmt(v) {
  if (!isFinite(v)) return String(v);
  const a = Math.abs(v);
  if (a >= 1e4 || (a < 1e-3 && a !== 0)) return v.toExponential(3);
  return v.toFixed(6);
}

function nernstBody() {
  return {
    standard_potential_v: Number($("e0").value),
    electrons: Number($("nElectrons").value),
    temperature_c: Number($("tempC").value),
    ox_activity: Number($("oxActivity").value),
    red_activity: Number($("redActivity").value),
  };
}

function ivBody() {
  return Object.assign(nernstBody(), {
    exchange_current_density: Number($("i0").value),
    alpha: Number($("alpha").value),
    eta_min_v: Number($("etaMin").value),
    eta_max_v: Number($("etaMax").value),
    eta_points: Number($("etaPoints").value),
  });
}

$("loadExample").addEventListener("click", async () => {
  clearError();
  const res = await fetch("/api/examples");
  if (!res.ok) {
    showError("failed to load examples");
    return;
  }
  const all = await res.json();
  const data = all[$("exampleSelect").value];
  if (!data) {
    showError("example not found");
    return;
  }
  $("e0").value = data.standard_potential_v;
  $("nElectrons").value = data.electrons;
  $("tempC").value = data.temperature_c;
  $("oxActivity").value = data.ox_activity;
  $("redActivity").value = data.red_activity;
  if (data.note) {
    $("nernstResult").innerHTML = `<p class="note">${data.note}</p>`;
  }
});

$("nernstBtn").addEventListener("click", async () => {
  try {
    const d = await postJSON("/api/nernst", nernstBody());
    $("nernstResult").innerHTML =
      `<p>E = ${fmt(d.equilibrium_potential_v)} V = ${fmt(d.equilibrium_potential_v * 1000)} mV</p>` +
      `<p>Nernst 斜率 = ${fmt(d.slope_mv_per_decade)} mV/decade · RT/F = ${fmt(d.thermal_voltage_v)} V · 单电子十倍移位 = ${fmt(d.decade_shift_v)} V</p>` +
      `<p>a_ox / a_red = ${fmt(d.activity_ratio)}</p>`;
  } catch (e) {
    showError(e.message);
  }
});

function polyline(points, xOf, yOf, color) {
  return points.map((p, i) => {
    const x = xOf(p.eta_v);
    const y = yOf(p.i);
    return (i === 0 ? "M" : "L") + x.toFixed(2) + " " + y.toFixed(2);
  }).join(" ");
}

function drawIV(d) {
  const svg = $("ivSvg");
  svg.innerHTML = "";
  const W = 680, H = 420;
  const mL = 70, mR = 30, mT = 20, mB = 50;
  const plotW = W - mL - mR;
  const plotH = H - mT - mB;

  const etas = d.points.map((p) => p.eta_v);
  const etaMin = Math.min(...etas), etaMax = Math.max(...etas);
  let maxAbs = 1;
  d.points.forEach((p) => {
    maxAbs = Math.max(maxAbs, Math.abs(p.i_bv), Math.abs(p.i_tafel));
  });

  const xOf = (eta) => mL + ((eta - etaMin) / (etaMax - etaMin)) * plotW;
  const yOf = (i) => mT + plotH / 2 - (i / maxAbs) * (plotH / 2 - 10);

  const defs = document.createElementNS("http://www.w3.org/2000/svg", "defs");
  const marker = document.createElementNS("http://www.w3.org/2000/svg", "marker");
  marker.setAttribute("id", "arrow");
  marker.setAttribute("viewBox", "0 0 10 10");
  marker.setAttribute("refX", "8");
  marker.setAttribute("refY", "5");
  marker.setAttribute("markerWidth", "6");
  marker.setAttribute("markerHeight", "6");
  marker.setAttribute("orient", "auto-start-reverse");
  marker.innerHTML = '<path d="M0 0 L10 5 L0 10 z" fill="#222"></path>';
  defs.appendChild(marker);
  svg.appendChild(defs);

  const hLine = document.createElementNS("http://www.w3.org/2000/svg", "line");
  hLine.setAttribute("x1", mL); hLine.setAttribute("y1", mT + plotH / 2);
  hLine.setAttribute("x2", mL + plotW); hLine.setAttribute("y2", mT + plotH / 2);
  hLine.setAttribute("stroke", "#ccc");
  svg.appendChild(hLine);
  const vLine = document.createElementNS("http://www.w3.org/2000/svg", "line");
  vLine.setAttribute("x1", mL); vLine.setAttribute("y1", mT);
  vLine.setAttribute("x2", mL); vLine.setAttribute("y2", mT + plotH);
  vLine.setAttribute("stroke", "#ccc");
  svg.appendChild(vLine);

  const xLabel = document.createElementNS("http://www.w3.org/2000/svg", "text");
  xLabel.setAttribute("x", mL + plotW / 2); xLabel.setAttribute("y", H - 12);
  xLabel.setAttribute("text-anchor", "middle");
  xLabel.setAttribute("fill", "#222");
  xLabel.textContent = "η (V)";
  svg.appendChild(xLabel);
  const yLabel = document.createElementNS("http://www.w3.org/2000/svg", "text");
  yLabel.setAttribute("x", 14); yLabel.setAttribute("y", mT + plotH / 2);
  yLabel.setAttribute("text-anchor", "middle");
  yLabel.setAttribute("fill", "#222");
  yLabel.textContent = "i";
  svg.appendChild(yLabel);

  const bv = document.createElementNS("http://www.w3.org/2000/svg", "path");
  bv.setAttribute("d", polyline(d.points, xOf, yOf, "#2563eb"));
  bv.setAttribute("fill", "none");
  bv.setAttribute("stroke", "#2563eb");
  bv.setAttribute("stroke-width", "2");
  svg.appendChild(bv);

  const taf = document.createElementNS("http://www.w3.org/2000/svg", "path");
  taf.setAttribute("d", polyline(d.points.map((p) => ({ eta_v: p.eta_v, i: p.i_tafel })), xOf, yOf, "#ea580c"));
  taf.setAttribute("fill", "none");
  taf.setAttribute("stroke", "#ea580c");
  taf.setAttribute("stroke-width", "2");
  taf.setAttribute("stroke-dasharray", "6 4");
  svg.appendChild(taf);

  const xTicks = [etaMin, 0, etaMax].map((v) => {
    const t = document.createElementNS("http://www.w3.org/2000/svg", "text");
    t.setAttribute("x", xOf(v)); t.setAttribute("y", mT + plotH / 2 + 16);
    t.setAttribute("text-anchor", "middle");
    t.setAttribute("fill", "#666");
    t.setAttribute("font-size", "11");
    t.textContent = fmt(v);
    svg.appendChild(t);
  });

  const legend = document.createElementNS("http://www.w3.org/2000/svg", "text");
  legend.setAttribute("x", mL + 8); legend.setAttribute("y", mT + 14);
  legend.setAttribute("fill", "#333");
  legend.setAttribute("font-size", "12");
  legend.textContent = "—— Butler–Volmer i(η)     - - - Tafel 参考 i(η)（来自 /api/iv）";
  svg.appendChild(legend);
}

$("ivBtn").addEventListener("click", async () => {
  try {
    const d = await postJSON("/api/iv", ivBody());
    const rows = d.points.map((p) =>
      `<tr><td>${fmt(p.eta_v)}</td><td>${fmt(p.i_bv)}</td><td>${fmt(p.i_tafel)}</td></tr>`
    ).join("");
    $("ivResult").innerHTML =
      `<p>E_eq = ${fmt(d.equilibrium_potential_v)} V · Tafel 斜率 b = ${fmt(d.tafel_slope_mv_per_decade)} mV/decade · α_a = ${fmt(d.alpha_anodic)} · α_c = ${fmt(d.alpha_cathodic)}</p>` +
      `<table class="tbl"><thead><tr><th>η (V)</th><th>i_BV</th><th>i_Tafel</th></tr></thead><tbody>${rows}</tbody></table>`;
    drawIV(d);
    $("ivNote").innerHTML = `<p class="note">大 |η| 时 BV 与 Tafel 参考重合；η = 0 时净电流为 0。图中点列全部来自 /api/iv 响应。</p>`;
  } catch (e) {
    showError(e.message);
  }
});
