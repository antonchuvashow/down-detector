'use strict';

const API_BASE = '/api/v1';

const DASHBOARDS_ENDPOINT = `${API_BASE}/superset/dashboards`;

const GUEST_TOKEN_ENDPOINT = `${API_BASE}/superset/guest-token`;

const SUPERSET_DOMAIN = window.__SUPERSET_DOMAIN__ || 'http://localhost:8088';

let currentDashboardId = null;
let embeddedDashboard  = null;

const tabsEl       = document.getElementById('tabs');
const embedEl      = document.getElementById('superset-embed');
const placeholder  = document.getElementById('embed-placeholder');
const errorBanner  = document.getElementById('embed-error');
const errorMsg     = document.getElementById('embed-error-msg');

function showError(msg) {
  if (errorBanner) {
    errorBanner.classList.remove('hidden');
    errorMsg.textContent = msg ? ` ${msg}` : '';
  }
}

function hideError() {
  errorBanner?.classList.add('hidden');
}

function showPlaceholder() {
  placeholder?.classList.remove('hidden');
}

function hidePlaceholder() {
  placeholder?.classList.add('hidden');
}

async function fetchGuestToken(dashboardId) {
  const url = `${GUEST_TOKEN_ENDPOINT}?dashboard=${encodeURIComponent(dashboardId)}`;
  const res = await fetch(url);
  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(body.error || `HTTP ${res.status}`);
  }
  const { token } = await res.json();
  return token;
}

async function embedDashboard(dashboardId) {
  hideError();
  showPlaceholder();

  if (embeddedDashboard) {
    try { await embeddedDashboard.unmount(); } catch (_) { /* ignore */ }
    embeddedDashboard = null;
  }

  while (embedEl.firstChild) {
    embedEl.removeChild(embedEl.firstChild);
  }
  embedEl.appendChild(placeholder);

  const sdk = window.supersetEmbeddedSdk;
  if (!sdk) {
    showError('Superset Embedded SDK не загружен.');
    hidePlaceholder();
    return;
  }

  try {
    embeddedDashboard = await sdk.embedDashboard({
      id:           dashboardId,
      supersetDomain: SUPERSET_DOMAIN,
      mountPoint:   embedEl,
      fetchGuestToken: () => fetchGuestToken(dashboardId),
      dashboardUiConfig: {
        hideTitle:         true,
        hideChartControls: false,
        filters: { expanded: true },
      },
    });

    hidePlaceholder();
    currentDashboardId = dashboardId;
  } catch (err) {
    console.error('[detector] embedDashboard error:', err);
    showError(err.message);
    hidePlaceholder();
  }
}

function buildTabs(dashboards) {
  tabsEl.innerHTML = '';

  dashboards.forEach((dash, idx) => {
    const btn = document.createElement('button');
    btn.type = 'button';
    btn.className = 'tab' + (idx === 0 ? ' tab--active' : '');
    btn.textContent = dash.name;
    btn.dataset.id = dash.id;
    btn.setAttribute('role', 'tab');
    btn.setAttribute('aria-selected', idx === 0 ? 'true' : 'false');

    btn.addEventListener('click', () => {
      if (btn.dataset.id === currentDashboardId) return;

      tabsEl.querySelectorAll('.tab').forEach(t => {
        t.classList.remove('tab--active');
        t.setAttribute('aria-selected', 'false');
      });
      btn.classList.add('tab--active');
      btn.setAttribute('aria-selected', 'true');

      embedDashboard(btn.dataset.id);
    });

    tabsEl.appendChild(btn);
  });
}

async function init() {
  let dashboards;

  try {
    const res = await fetch(DASHBOARDS_ENDPOINT);
    if (res.ok) {
      dashboards = await res.json();
    } else {
      throw new Error(`Endpoint returned ${res.status}`);
    }
  } catch (err) {
    console.warn('[detector] Could not fetch dashboards from backend, using fallback.', err);
    dashboards = window.__SUPERSET_DASHBOARDS__ || [];
  }

  if (!dashboards || dashboards.length === 0) {
    showError('Нет доступных дашбордов. Настройте SUPERSET_DASHBOARD_IDS в конфиге сервера.');
    hidePlaceholder();
    return;
  }

  buildTabs(dashboards.dashboards);
  embedDashboard(dashboards.dashboards[0].id);
}

document.addEventListener('DOMContentLoaded', init);
