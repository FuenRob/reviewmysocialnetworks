import type { AccountReport } from '../types/instagram';

const API_BASE = '/api';

export async function getAuthURL(mode?: 'business' | 'basic'): Promise<{ auth_url: string }> {
  const url = mode === 'basic' ? `${API_BASE}/auth/url?mode=basic` : `${API_BASE}/auth/url`;
  const res = await fetch(url);
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || 'Error al generar la URL de autorización de Instagram');
  }
  return res.json();
}

export async function getAuthResult(): Promise<AccountReport> {
  const res = await fetch(`${API_BASE}/auth/result`, { credentials: 'same-origin' });
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || 'No se pudo completar el análisis autenticado');
  }
  return res.json();
}

export async function analyzeWithToken(accessToken: string, accountId?: string): Promise<AccountReport> {
  const res = await fetch(`${API_BASE}/analyze/token`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      access_token: accessToken,
      account_id: accountId || undefined,
    }),
  });
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || 'Error al analizar la cuenta con el token de Instagram');
  }
  return res.json();
}

export async function analyzeDemo(tier: 'A' | 'B' | 'D' | 'F'): Promise<AccountReport> {
  const res = await fetch(`${API_BASE}/analyze/demo`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ tier }),
  });
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || 'Error al cargar cuenta demo');
  }
  return res.json();
}
