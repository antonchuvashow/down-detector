'use strict';


const API_BASE = '/api/v1';
const ROUTES_ENDPOINT = `${API_BASE}/routes`;
const REPORTS_ENDPOINT = `${API_BASE}/reports`;


const form = document.getElementById('report-form');
const routeSelect = document.getElementById('route-select');
const routeError = document.getElementById('route-error');
const btnSuccess = document.getElementById('btn-success');
const btnError = document.getElementById('btn-error');
const statusValue = document.getElementById('status-value');
const errorTypesGrp = document.getElementById('error-types-group');
const errorTypesErr = document.getElementById('error-types-error');
const submitBtn = document.getElementById('submit-btn');
const successBanner = document.getElementById('success-banner');
const errorBanner = document.getElementById('error-banner');
const errorMsg = document.getElementById('error-msg');

function detectPlatform() {
  const raw = [
    navigator.userAgentData?.platform || '',
    navigator.platform || '',
    navigator.userAgent || '',
  ]
      .join(' ')
      .toLowerCase();

  if (raw.includes('android')) return 'android';
  if (/iphone|ipad|ipod|ios|macos/.test(raw)) return 'ios';
  if (raw.includes('windows')) return 'windows';
  if (raw.includes('linux') && !raw.includes('android')) return 'linux';

  return 'unknown';
}

async function measureLatencyMs(url) {
  if (!url) return null;
  let timeoutMs = 2000

  const startedAt = performance.now();
  const controller = new AbortController();

  const timeoutId = setTimeout(() => controller.abort(), timeoutMs);

  try {
    await fetch(url, {
      method: 'GET',
      mode: 'no-cors',
      cache: 'no-store',
      signal: controller.signal,
    });

    return Math.round(performance.now() - startedAt);
  } catch (err) {
    console.warn('[detector] latency measurement failed:', err);
    return null;
  } finally {
    clearTimeout(timeoutId);
  }
}

function getGeo() {
  return new Promise((resolve) => {
    if (!navigator.geolocation) {
      console.warn('[detector] geolocation is not supported');
      resolve({ latitude: null, longitude: null });
      return;
    }

    navigator.geolocation.getCurrentPosition(
        (pos) => {
          resolve({
            latitude: pos.coords.latitude,
            longitude: pos.coords.longitude,
          });
        },
        (err) => {
          console.warn('[detector] geolocation failed:', {
            code: err.code,
            message: err.message,
            secureContext: window.isSecureContext,
            origin: window.location.origin,
          });

          resolve({ latitude: null, longitude: null });
        },
        {
          enableHighAccuracy: false,
          timeout: 5000,
          maximumAge: 60000,
        }
    );
  });
}

async function collectReportMeta(routeUrl) {
  const [latencyMs, geo] = await Promise.all([
    measureLatencyMs(routeUrl),
    getGeo(),
  ]);

  return {
    latency_ms: latencyMs,
    latitude: geo.latitude,
    longitude: geo.longitude,
    platform: detectPlatform(),
  };
}


async function loadRoutes() {
  try {
    const res = await fetch(ROUTES_ENDPOINT);
    if (!res.ok) throw new Error(`HTTP ${res.status}`);

    const routes = await res.json();
    routeSelect.innerHTML = '';

    if (!routes?.routes || routes.routes.length === 0) {
      const opt = document.createElement('option');
      opt.value = '';
      opt.disabled = true;
      opt.selected = true;
      opt.textContent = 'Нет доступных маршрутов';
      routeSelect.appendChild(opt);
      return;
    }

    routes.routes.forEach((route) => {
      const opt = document.createElement('option');
      opt.value = route.id;
      opt.textContent = route.url || route.id;
      opt.dataset.url = route.url || '';
      routeSelect.appendChild(opt);
    });
  } catch (err) {
    console.error('[detector] loadRoutes:', err);
    routeSelect.innerHTML = '<option value="" disabled selected>Ошибка загрузки маршрутов</option>';
  }
}

function setStatus(isSuccess) {
  statusValue.value = String(isSuccess);

  btnSuccess.setAttribute('aria-pressed', String(isSuccess));
  btnError.setAttribute('aria-pressed', String(!isSuccess));

  if (isSuccess) {
    errorTypesGrp.classList.add('hidden');
    errorTypesErr.classList.add('hidden');
  } else {
    errorTypesGrp.classList.remove('hidden');
  }
}

btnSuccess.addEventListener('click', () => setStatus(true));
btnError.addEventListener('click', () => setStatus(false));

function validate() {
  let valid = true;

  const routeId = routeSelect.value;
  if (!routeId) {
    routeError.textContent = 'Выберите маршрут';
    valid = false;
  } else {
    routeError.textContent = '';
  }

  const isSuccess = statusValue.value === 'true';
  const checkboxes = form.querySelectorAll('input[name="error_types"]:checked');
  const errorTypes = Array.from(checkboxes).map((cb) => cb.value);

  if (!isSuccess && errorTypes.length === 0) {
    errorTypesErr.classList.remove('hidden');
    valid = false;
  } else {
    errorTypesErr.classList.add('hidden');
  }

  return { valid, routeId, isSuccess, errorTypes };
}


form.addEventListener('submit', async (e) => {
  e.preventDefault();

  successBanner.classList.add('hidden');
  errorBanner.classList.add('hidden');

  const { valid, routeId, isSuccess, errorTypes } = validate();
  if (!valid) return;

  submitBtn.classList.add('btn--loading');
  submitBtn.disabled = true;

  try {
    const selectedOption = routeSelect.selectedOptions[0];
    const routeUrl = selectedOption?.dataset?.url || '';

    const meta = await collectReportMeta(routeUrl);

    const payload = {
      route_id: routeId,
      success: isSuccess,
      error_types: errorTypes,
      platform: meta.platform,
      latency_ms: meta.latency_ms,
      latitude: meta.latitude,
      longitude: meta.longitude,
      source: 'user',
    };

    const res = await fetch(REPORTS_ENDPOINT, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    });

    if (res.status === 204 || res.ok) {
      successBanner.classList.remove('hidden');
      form.reset();
      setStatus(false);
      successBanner.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
    } else {
      const body = await res.json().catch(() => ({}));
      throw new Error(body.error || `Сервер вернул ${res.status}`);
    }
  } catch (err) {
    errorBanner.classList.remove('hidden');
    errorMsg.textContent = `${err.message}`;
    console.error('[detector] submit report:', err);
  } finally {
    submitBtn.classList.remove('btn--loading');
    submitBtn.disabled = false;
  }
});


document.addEventListener('DOMContentLoaded', () => {
  loadRoutes();
  setStatus(false);
});