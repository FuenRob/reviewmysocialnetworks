import React, { useState } from 'react';
import { Sparkles, ArrowRight, Shield, BarChart3, TrendingUp, Calendar, CheckCircle2 } from 'lucide-react';
import { InstagramIcon } from './InstagramIcon';
import { TikTokIcon } from './TikTokIcon';
import { getAuthURL, getTikTokAuthURL, analyzeDemo, analyzeTikTokDemo } from '../api/client';
import type { AccountReport } from '../types/instagram';

interface Props { onReportLoaded: (report: AccountReport) => void; }
type Platform = 'instagram' | 'tiktok';
type Tier = 'A' | 'B' | 'D' | 'F';

const demos: Record<Platform, Record<Tier, { username: string; detail: string }>> = {
  instagram: {
    A: { username: 'codeandcoffee.dev', detail: '48.5K seg · 6.8% engagement · Reels + Carruseles' },
    B: { username: 'atelier.cafe.madrid', detail: '19.2K seg · base sólida · calidad consistente' },
    D: { username: 'fitness_routine_daily', detail: '14.5K seg · interacción baja · ritmo irregular' },
    F: { username: 'crypto_signals_daily', detail: '32.4K seg · cuenta inactiva · audiencia desconectada' },
  },
  tiktok: {
    A: { username: 'cienciaen60s', detail: '120K seg · 260K vistas/vídeo · alta viralidad' },
    B: { username: 'cocinaconmarta', detail: '52K seg · 43K vistas/vídeo · base sólida' },
    D: { username: 'entrenaencasa', detail: '28K seg · 7K vistas/vídeo · alcance limitado' },
    F: { username: 'ofertas_flash_ya', detail: '46K seg · 3.5K vistas/vídeo · inactividad' },
  },
};

const tierStyles: Record<Tier, string> = {
  A: 'border-emerald-500/30 hover:border-emerald-500/60 text-emerald-400',
  B: 'border-blue-500/30 hover:border-blue-500/60 text-blue-400',
  D: 'border-amber-500/30 hover:border-amber-500/60 text-amber-400',
  F: 'border-rose-500/30 hover:border-rose-500/60 text-rose-400',
};

export const HeroConnect: React.FC<Props> = ({ onReportLoaded }) => {
  const [platform, setPlatform] = useState<Platform>('instagram');
  const [loading, setLoading] = useState(false);
  const [loadingDemoTier, setLoadingDemoTier] = useState<Tier | null>(null);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const isTikTok = platform === 'tiktok';
  const PlatformIcon = isTikTok ? TikTokIcon : InstagramIcon;
  const platformName = isTikTok ? 'TikTok' : 'Instagram';

  const handleOAuthLogin = async () => {
    setLoading(true); setErrorMsg(null);
    try {
      const data = isTikTok ? await getTikTokAuthURL() : await getAuthURL();
      window.location.href = data.auth_url;
    } catch (err: any) {
      setErrorMsg(err.message || `Error al iniciar sesión con ${platformName}`);
      setLoading(false);
    }
  };

  const handleLoadDemo = async (tier: Tier) => {
    setLoadingDemoTier(tier); setErrorMsg(null);
    try {
      const report = isTikTok ? await analyzeTikTokDemo(tier) : await analyzeDemo(tier);
      onReportLoaded(report);
    } catch (err: any) {
      setErrorMsg(err.message || `Error al cargar la demo de ${platformName}`);
    } finally { setLoadingDemoTier(null); }
  };

  const features = isTikTok ? [
    ['Engagement por vista', 'Likes, comentarios y compartidos sobre visualizaciones.', BarChart3],
    ['Alcance y viralidad', 'Vistas por seguidor, compartidos y vídeos virales.', TrendingUp],
    ['Cadencia de vídeos', 'Frecuencia, regularidad y mejores momentos.', Calendar],
    ['Plan de crecimiento', 'Acciones priorizadas según tus datos reales.', CheckCircle2],
  ] : [
    ['Engagement real', 'Likes, comentarios, guardados e interacción por post.', BarChart3],
    ['Mezcla de contenido', 'Rendimiento de Reels, carruseles y fotos.', TrendingUp],
    ['Cadencia de posts', 'Frecuencia, regularidad y mejores momentos.', Calendar],
    ['Plan de crecimiento', 'Acciones priorizadas según tus datos reales.', CheckCircle2],
  ];

  return <div className="space-y-10 max-w-5xl mx-auto py-6">
    <div className="flex justify-center">
      <div className="p-1.5 rounded-2xl bg-slate-900 border border-slate-800 flex gap-1.5" role="tablist" aria-label="Red social a analizar">
        {(['instagram', 'tiktok'] as const).map((value) => {
          const Icon = value === 'instagram' ? InstagramIcon : TikTokIcon;
          return <button key={value} onClick={() => { setPlatform(value); setErrorMsg(null); }} role="tab" aria-selected={platform === value}
            className={`px-5 py-2.5 rounded-xl text-sm font-bold flex items-center gap-2 transition-all ${platform === value ? 'bg-indigo-600 text-white shadow-lg' : 'text-slate-400 hover:text-white hover:bg-slate-800'}`}>
            <Icon className="w-4 h-4" />{value === 'instagram' ? 'Instagram' : 'TikTok'}
          </button>;
        })}
      </div>
    </div>

    <div className="text-center space-y-4">
      <div className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full bg-indigo-500/10 border border-indigo-500/20 text-indigo-300 text-xs font-semibold">
        <Sparkles className="w-3.5 h-3.5" /> Auditoría oficial · Calificación A, B, D o F
      </div>
      <h1 className="text-4xl sm:text-6xl font-black text-white tracking-tight leading-tight">Audita y mejora tu<br className="hidden sm:inline" /> <span className={isTikTok ? 'bg-gradient-to-r from-cyan-400 via-white to-rose-400 bg-clip-text text-transparent' : 'bg-gradient-to-r from-yellow-400 via-pink-500 to-purple-500 bg-clip-text text-transparent'}>cuenta de {platformName}</span></h1>
      <p className="text-base sm:text-lg text-slate-400 max-w-3xl mx-auto leading-relaxed">Conecta tu cuenta para obtener un informe completo, una nota comparable y recomendaciones concretas para mejorar el alcance, la interacción y el crecimiento.</p>
    </div>

    {errorMsg && <div className="p-4 rounded-2xl bg-rose-500/10 border border-rose-500/30 text-rose-400 text-sm font-medium max-w-2xl mx-auto">⚠️ {errorMsg}</div>}

    <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-stretch">
      <section className="lg:col-span-7 bg-slate-900/90 border border-slate-800 rounded-3xl p-6 sm:p-8 flex flex-col justify-between space-y-8 shadow-2xl">
        <div className="space-y-6">
          <div className="flex items-center gap-3"><div className={`w-12 h-12 rounded-2xl p-0.5 ${isTikTok ? 'bg-gradient-to-tr from-cyan-400 to-rose-500' : 'bg-gradient-to-tr from-pink-500 to-purple-600'}`}><div className="w-full h-full bg-slate-950 rounded-[14px] flex items-center justify-center"><PlatformIcon className="w-6 h-6 text-white" /></div></div><div><h2 className="text-xl font-bold text-white">Conexión con {platformName}</h2><p className="text-xs text-slate-400">Acceso seguro y de solo lectura mediante OAuth</p></div></div>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">{features.map(([title, detail, Icon]: any) => <div key={title} className="p-3 rounded-2xl bg-slate-950/60 border border-slate-800/80 flex items-start gap-2.5"><Icon className="w-4 h-4 text-indigo-400 shrink-0 mt-0.5" /><div><h4 className="text-xs font-bold text-white">{title}</h4><p className="text-[11px] text-slate-400 mt-0.5">{detail}</p></div></div>)}</div>
          <button onClick={handleOAuthLogin} disabled={loading} className={`w-full py-4 px-6 rounded-2xl text-white font-bold text-base shadow-xl flex items-center justify-center gap-3 transition-all active:scale-95 disabled:opacity-50 group ${isTikTok ? 'bg-gradient-to-r from-cyan-500 via-slate-800 to-rose-500' : 'bg-gradient-to-r from-pink-500 via-purple-600 to-indigo-600'}`}><PlatformIcon className="w-5 h-5" />{loading ? `Conectando con ${platformName}...` : `Iniciar sesión con ${platformName}`}<ArrowRight className="w-4 h-4 ml-auto" /></button>
        </div>
        <div className="text-[11px] text-slate-500 flex items-start gap-2 border-t border-slate-800/80 pt-4"><Shield className="w-4 h-4 text-emerald-400 shrink-0" /><span>Autenticación oficial. Solo leemos los datos necesarios para generar el informe; no publicamos ni modificamos tu cuenta.</span></div>
      </section>

      <section className="lg:col-span-5 bg-slate-900/90 border border-slate-800 rounded-3xl p-6 sm:p-8 shadow-2xl">
        <div className="flex items-center gap-2 mb-1"><Sparkles className="w-4 h-4 text-amber-400" /><h3 className="text-lg font-bold text-white">Probar cuentas demo</h3></div>
        <p className="text-xs text-slate-400 mb-4">Explora los cuatro niveles con datos de ejemplo de {platformName}:</p>
        <div className="space-y-2.5">{(['A','B','D','F'] as Tier[]).map((tier) => <button key={tier} onClick={() => handleLoadDemo(tier)} disabled={Boolean(loadingDemoTier)} className={`w-full p-3.5 rounded-2xl bg-slate-950/70 border text-left transition-all flex items-center justify-between group ${tierStyles[tier]}`}><div><div className="flex items-center gap-2"><span className="w-6 h-6 rounded-lg bg-current text-slate-950 font-black text-xs flex items-center justify-center"><span className="text-slate-950">{tier}</span></span><span className="text-xs font-bold text-white">@{demos[platform][tier].username}</span></div><p className="text-[11px] text-slate-400 mt-1 pl-8">{demos[platform][tier].detail}</p></div><span className="text-xs font-bold">{loadingDemoTier === tier ? '...' : 'Ver →'}</span></button>)}</div>
        <div className="text-[11px] text-slate-500 bg-slate-950/40 p-3 rounded-2xl border border-slate-800 mt-5">💡 La demo muestra el informe completo sin conectar una cuenta real.</div>
      </section>
    </div>
  </div>;
};
