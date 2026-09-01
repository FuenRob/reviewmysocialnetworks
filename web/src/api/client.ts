import type { AccountReport } from '../types/instagram';

const API_BASE = '/api';

export async function getAuthURL(mode?: 'business' | 'basic'): Promise<{ auth_url: string }> {
  const url = mode === 'basic' ? `${API_BASE}/instagram/auth/url?mode=basic` : `${API_BASE}/instagram/auth/url`;
  const res = await fetch(url);
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || 'Error al generar la URL de autorización de Instagram');
  }
  return res.json();
}

export async function getAuthResult(): Promise<AccountReport> {
  const res = await fetch(`${API_BASE}/instagram/auth/result`, { credentials: 'same-origin' });
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || 'No se pudo completar el análisis autenticado');
  }
  return res.json();
}

export async function analyzeDemo(tier: 'A' | 'B' | 'D' | 'F'): Promise<AccountReport> {
  const res = await fetch(`${API_BASE}/instagram/analyze/demo`, {
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

export async function getTikTokAuthURL(): Promise<{ auth_url: string }> {
  const res = await fetch(`${API_BASE}/tiktok/auth/url`);
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || 'Error al generar la URL de autorización de TikTok');
  }
  return res.json();
}

export async function getTikTokAuthResult(): Promise<AccountReport> {
  const res = await fetch(`${API_BASE}/tiktok/auth/result`, { credentials: 'same-origin' });
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || 'No se pudo completar el análisis de TikTok');
  }
  return res.json();
}

export async function analyzeTikTokDemo(tier: 'A' | 'B' | 'D' | 'F'): Promise<AccountReport> {
  const res = await fetch(`${API_BASE}/tiktok/analyze/demo`, {
    method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ tier }),
  });
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || 'Error al cargar la cuenta demo de TikTok');
  }
  return res.json();
}
