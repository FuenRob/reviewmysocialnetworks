import React, { useState } from 'react';
import { Sparkles, ArrowRight, Shield, BarChart3, TrendingUp, Calendar, CheckCircle2 } from 'lucide-react';
import { InstagramIcon } from './InstagramIcon';
import { getAuthURL, analyzeDemo } from '../api/client';
import type { AccountReport } from '../types/instagram';

interface Props {
  onReportLoaded: (report: AccountReport) => void;
}

export const HeroConnect: React.FC<Props> = ({ onReportLoaded }) => {
  const [loading, setLoading] = useState(false);
  const [loadingDemoTier, setLoadingDemoTier] = useState<string | null>(null);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  const handleOAuthLogin = async () => {
    setLoading(true);
    setErrorMsg(null);
    try {
      const data = await getAuthURL();
      window.location.href = data.auth_url;
    } catch (err: any) {
      setErrorMsg(err.message || 'Error al iniciar sesión con Instagram');
      setLoading(false);
    }
  };

  const handleLoadDemo = async (tier: 'A' | 'B' | 'D' | 'F') => {
    setLoadingDemoTier(tier);
    setErrorMsg(null);
    try {
      const report = await analyzeDemo(tier);
      onReportLoaded(report);
    } catch (err: any) {
      setErrorMsg(err.message || 'Error al cargar la cuenta demo');
    } finally {
      setLoadingDemoTier(null);
    }
  };

  return (
    <div className="space-y-12 max-w-5xl mx-auto py-6">
      <div className="text-center space-y-4">
        <div className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full bg-gradient-to-r from-pink-500/10 via-purple-500/10 to-yellow-500/10 border border-pink-500/20 text-pink-300 text-xs font-semibold backdrop-blur-md">
          <Sparkles className="w-3.5 h-3.5 text-pink-400" />
          Auditoría Oficial Instagram Graph API • Calificación Algorítmica A, B, D o F
        </div>

        <h1 className="text-4xl sm:text-6xl font-black text-white tracking-tight leading-tight">
          Audita y analiza tu <br className="hidden sm:inline" />
          <span className="bg-gradient-to-r from-yellow-400 via-pink-500 to-purple-500 bg-clip-text text-transparent">
            Cuenta de Instagram
          </span>
        </h1>

        <p className="text-base sm:text-lg text-slate-400 max-w-2xl mx-auto leading-relaxed">
          Conecta tu perfil de Instagram para obtener un diagnóstico completo sobre tu audiencia, tasa de interacción real, ritmo de publicación y una nota final de rendimiento.
        </p>
      </div>

      {errorMsg && (
        <div className="p-4 rounded-2xl bg-rose-500/10 border border-rose-500/30 text-rose-400 text-sm font-medium flex items-center gap-3 max-w-2xl mx-auto">
          <span className="shrink-0">⚠️</span>
          <span>{errorMsg}</span>
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-stretch">
        <div className="lg:col-span-7 bg-slate-900/90 border border-slate-800 rounded-3xl p-6 sm:p-8 backdrop-blur-xl flex flex-col justify-between space-y-8 shadow-2xl">
          <div className="space-y-6">
            <div className="flex items-center gap-3">
              <div className="w-12 h-12 rounded-2xl bg-gradient-to-tr from-pink-500 to-purple-600 p-0.5 shadow-lg shadow-pink-500/20">
                <div className="w-full h-full bg-slate-950 rounded-[14px] flex items-center justify-center">
                  <InstagramIcon className="w-6 h-6 text-pink-400" />
                </div>
              </div>
              <div>
                <h2 className="text-xl font-bold text-white">
                  Conexión con Instagram
                </h2>
                <p className="text-xs text-slate-400">
                  Acceso seguro y directo vía Instagram OAuth
                </p>
              </div>
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 pt-2">
              <div className="p-3 rounded-2xl bg-slate-950/60 border border-slate-800/80 flex items-start gap-2.5">
                <BarChart3 className="w-4 h-4 text-pink-400 shrink-0 mt-0.5" />
                <div>
                  <h4 className="text-xs font-bold text-white">Engagement Rate Real</h4>
                  <p className="text-[11px] text-slate-400 mt-0.5">Likes, comentarios e interacción por publicación.</p>
                </div>
              </div>

              <div className="p-3 rounded-2xl bg-slate-950/60 border border-slate-800/80 flex items-start gap-2.5">
                <Calendar className="w-4 h-4 text-purple-400 shrink-0 mt-0.5" />
                <div>
                  <h4 className="text-xs font-bold text-white">Cadencia de Posteo</h4>
                  <p className="text-[11px] text-slate-400 mt-0.5">Frecuencia semanal e intervalo entre publicaciones.</p>
                </div>
              </div>

              <div className="p-3 rounded-2xl bg-slate-950/60 border border-slate-800/80 flex items-start gap-2.5">
                <TrendingUp className="w-4 h-4 text-emerald-400 shrink-0 mt-0.5" />
                <div>
                  <h4 className="text-xs font-bold text-white">Calificación A / B / D / F</h4>
                  <p className="text-[11px] text-slate-400 mt-0.5">Diagnóstico ejecutivo y nota global del perfil.</p>
                </div>
              </div>

              <div className="p-3 rounded-2xl bg-slate-950/60 border border-slate-800/80 flex items-start gap-2.5">
                <CheckCircle2 className="w-4 h-4 text-indigo-400 shrink-0 mt-0.5" />
                <div>
                  <h4 className="text-xs font-bold text-white">Plan de Acción</h4>
                  <p className="text-[11px] text-slate-400 mt-0.5">Recomendaciones prioritarias personalizadas.</p>
                </div>
              </div>
            </div>

            <div className="pt-2">
              <button
                onClick={handleOAuthLogin}
                disabled={loading}
                className="w-full py-4 px-6 rounded-2xl bg-gradient-to-r from-pink-500 via-purple-600 to-indigo-600 hover:from-pink-600 hover:via-purple-700 hover:to-indigo-700 text-white font-bold text-base shadow-xl shadow-pink-500/25 flex items-center justify-center gap-3 transition-all transform active:scale-95 disabled:opacity-50 group"
              >
                <InstagramIcon className="w-5 h-5 group-hover:scale-110 transition-transform" />
                {loading ? 'Conectando con Instagram...' : 'Iniciar Sesión con Instagram'}
                <ArrowRight className="w-4 h-4 ml-auto group-hover:translate-x-1 transition-transform" />
              </button>
            </div>
          </div>

          <div className="text-[11px] text-slate-500 flex items-center gap-2 border-t border-slate-800/80 pt-4">
            <Shield className="w-4 h-4 text-emerald-400 shrink-0" />
            <span>Autenticación oficial y segura contra instagram.com. Privacidad garantizada.</span>
          </div>
        </div>

        <div className="lg:col-span-5 bg-slate-900/90 border border-slate-800 rounded-3xl p-6 sm:p-8 backdrop-blur-xl flex flex-col justify-between space-y-5 shadow-2xl">
          <div>
            <div className="flex items-center gap-2 mb-1">
              <Sparkles className="w-4 h-4 text-amber-400" />
              <h3 className="text-lg font-bold text-white">
                Probar Cuentas Demo
              </h3>
            </div>
            <p className="text-xs text-slate-400 mb-4">
              Explora al instante cómo evalúa el algoritmo los 4 niveles de cuenta:
            </p>

            <div className="space-y-2.5">
              <button
                onClick={() => handleLoadDemo('A')}
                disabled={Boolean(loadingDemoTier)}
                className="w-full p-3.5 rounded-2xl bg-slate-950/70 border border-emerald-500/30 hover:border-emerald-500/60 hover:bg-emerald-950/20 text-left transition-all flex items-center justify-between group"
              >
                <div>
                  <div className="flex items-center gap-2">
                    <span className="w-6 h-6 rounded-lg bg-emerald-500 text-slate-950 font-black text-xs flex items-center justify-center shadow-sm shadow-emerald-500/30">
                      A
                    </span>
                    <span className="text-xs font-bold text-white group-hover:text-emerald-300">
                      @codeandcoffee.dev
                    </span>
                  </div>
                  <p className="text-[11px] text-slate-400 mt-1 pl-8">
                    48.5K seg • 6.8% engagement • Reels + Carruseles
                  </p>
                </div>
                <span className="text-xs font-bold text-emerald-400 group-hover:translate-x-1 transition-transform">
                  {loadingDemoTier === 'A' ? '...' : 'Ver →'}
                </span>
              </button>

              <button
                onClick={() => handleLoadDemo('B')}
                disabled={Boolean(loadingDemoTier)}
                className="w-full p-3.5 rounded-2xl bg-slate-950/70 border border-blue-500/30 hover:border-blue-500/60 hover:bg-blue-950/20 text-left transition-all flex items-center justify-between group"
              >
                <div>
                  <div className="flex items-center gap-2">
                    <span className="w-6 h-6 rounded-lg bg-blue-500 text-white font-black text-xs flex items-center justify-center shadow-sm shadow-blue-500/30">
                      B
                    </span>
                    <span className="text-xs font-bold text-white group-hover:text-blue-300">
                      @atelier.cafe.madrid
                    </span>
                  </div>
                  <p className="text-[11px] text-slate-400 mt-1 pl-8">
                    19.2K seg • 2.4% engagement • Calidad consistente
                  </p>
                </div>
                <span className="text-xs font-bold text-blue-400 group-hover:translate-x-1 transition-transform">
                  {loadingDemoTier === 'B' ? '...' : 'Ver →'}
                </span>
              </button>

              <button
                onClick={() => handleLoadDemo('D')}
                disabled={Boolean(loadingDemoTier)}
                className="w-full p-3.5 rounded-2xl bg-slate-950/70 border border-amber-500/30 hover:border-amber-500/60 hover:bg-amber-950/20 text-left transition-all flex items-center justify-between group"
              >
                <div>
                  <div className="flex items-center gap-2">
                    <span className="w-6 h-6 rounded-lg bg-amber-500 text-slate-950 font-black text-xs flex items-center justify-center shadow-sm shadow-amber-500/30">
                      D
                    </span>
                    <span className="text-xs font-bold text-white group-hover:text-amber-300">
                      @fitness_routine_daily
                    </span>
                  </div>
                  <p className="text-[11px] text-slate-400 mt-1 pl-8">
                    14.5K seg • 0.9% engagement • Publicación irregular
                  </p>
                </div>
                <span className="text-xs font-bold text-amber-400 group-hover:translate-x-1 transition-transform">
                  {loadingDemoTier === 'D' ? '...' : 'Ver →'}
                </span>
              </button>

              <button
                onClick={() => handleLoadDemo('F')}
                disabled={Boolean(loadingDemoTier)}
                className="w-full p-3.5 rounded-2xl bg-slate-950/70 border border-rose-500/30 hover:border-rose-500/60 hover:bg-rose-950/20 text-left transition-all flex items-center justify-between group"
              >
                <div>
                  <div className="flex items-center gap-2">
                    <span className="w-6 h-6 rounded-lg bg-rose-500 text-white font-black text-xs flex items-center justify-center shadow-sm shadow-rose-500/30">
                      F
                    </span>
                    <span className="text-xs font-bold text-white group-hover:text-rose-300">
                      @crypto_signals_daily
                    </span>
                  </div>
                  <p className="text-[11px] text-slate-400 mt-1 pl-8">
                    32.4K seg • 0.08% engagement • Cuenta inactiva
                  </p>
                </div>
                <span className="text-xs font-bold text-rose-400 group-hover:translate-x-1 transition-transform">
                  {loadingDemoTier === 'F' ? '...' : 'Ver →'}
                </span>
              </button>
            </div>
          </div>

          <div className="text-[11px] text-slate-500 bg-slate-950/40 p-3 rounded-2xl border border-slate-800">
            💡 Haz clic en una demo para ver la estructura del informe completo y sus gráficos interactivos.
          </div>
        </div>
      </div>
    </div>
  );
};
