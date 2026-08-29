import React from 'react';
import type { AccountReport, Recommendation } from '../types/instagram';
import { CheckCircle2, AlertTriangle, ArrowUpRight, Target, ShieldCheck, Flame, Database } from 'lucide-react';

interface Props {
  report: AccountReport;
}

export const ReportActionPlan: React.FC<Props> = ({ report }) => {
  const { strengths, weaknesses, recommendations } = report;

  const getPriorityBadge = (p: Recommendation['priority']) => {
    switch (p) {
      case 'Alta':
        return 'bg-rose-500/10 text-rose-400 border-rose-500/20';
      case 'Media':
        return 'bg-amber-500/10 text-amber-400 border-amber-500/20';
      case 'Baja':
      default:
        return 'bg-blue-500/10 text-blue-400 border-blue-500/20';
    }
  };

  return (
    <div className="space-y-8">
      {report.platform === 'tiktok' && (
        <div className="bg-cyan-950/20 border border-cyan-500/20 rounded-3xl p-5 flex flex-col sm:flex-row gap-4">
          <Database className="w-5 h-5 text-cyan-400 shrink-0 mt-0.5" />
          <div className="space-y-2">
            <h3 className="text-sm font-bold text-white">Cobertura del análisis oficial</h3>
            <p className="text-xs text-slate-300">Se analizaron {report.data_coverage.analyzed_posts} de hasta {report.data_coverage.max_posts} vídeos recientes: {report.data_coverage.available.join(', ')}.</p>
            {report.data_coverage.unavailable?.length ? <p className="text-[11px] text-slate-500">TikTok Display API no facilita: {report.data_coverage.unavailable.join(', ')}. El informe no estima estos datos para evitar conclusiones engañosas.</p> : null}
          </div>
        </div>
      )}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="bg-slate-900/80 border border-slate-800 rounded-3xl p-6 backdrop-blur-xl">
          <div className="flex items-center gap-3 mb-5">
            <div className="p-2.5 rounded-2xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-400">
              <ShieldCheck className="w-5 h-5" />
            </div>
            <div>
              <h3 className="text-base sm:text-lg font-bold text-white">
                Puntos Fuertes Detectados
              </h3>
              <p className="text-xs text-slate-400">
                Aspectos sobresalientes de tu cuenta
              </p>
            </div>
          </div>

          <ul className="space-y-3">
            {strengths.map((s, idx) => (
              <li
                key={idx}
                className="flex items-start gap-3 bg-slate-950/50 border border-slate-800/80 rounded-2xl p-3.5 text-xs sm:text-sm text-slate-200"
              >
                <CheckCircle2 className="w-4 h-4 text-emerald-400 shrink-0 mt-0.5" />
                <span>{s}</span>
              </li>
            ))}
          </ul>
        </div>

        <div className="bg-slate-900/80 border border-slate-800 rounded-3xl p-6 backdrop-blur-xl">
          <div className="flex items-center gap-3 mb-5">
            <div className="p-2.5 rounded-2xl bg-rose-500/10 border border-rose-500/20 text-rose-400">
              <AlertTriangle className="w-5 h-5" />
            </div>
            <div>
              <h3 className="text-base sm:text-lg font-bold text-white">
                Áreas de Mejora Críticas
              </h3>
              <p className="text-xs text-slate-400">
                Factores que limitan tu crecimiento
              </p>
            </div>
          </div>

          <ul className="space-y-3">
            {weaknesses.map((w, idx) => (
              <li
                key={idx}
                className="flex items-start gap-3 bg-slate-950/50 border border-slate-800/80 rounded-2xl p-3.5 text-xs sm:text-sm text-slate-200"
              >
                <AlertTriangle className="w-4 h-4 text-rose-400 shrink-0 mt-0.5" />
                <span>{w}</span>
              </li>
            ))}
          </ul>
        </div>
      </div>

      <div className="bg-slate-900/80 border border-slate-800 rounded-3xl p-6 sm:p-8 backdrop-blur-xl">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
          <div className="flex items-center gap-3">
            <div className="p-2.5 rounded-2xl bg-indigo-500/10 border border-indigo-500/20 text-indigo-400">
              <Target className="w-5 h-5" />
            </div>
            <div>
              <h3 className="text-base sm:text-lg font-bold text-white">
                Plan de Acción y Recomendaciones
              </h3>
              <p className="text-xs text-slate-400">
                Pasos recomendados para subir a la siguiente escala de calificación
              </p>
            </div>
          </div>

          <div className="flex items-center gap-2 text-xs font-semibold px-3 py-1.5 rounded-full bg-slate-800 text-slate-300 border border-slate-700 w-fit">
            <Flame className="w-3.5 h-3.5 text-amber-400" />
            Objetivo: Grado A (Sobresaliente)
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {recommendations.map((rec, idx) => (
            <div
              key={idx}
              className="bg-slate-950/70 border border-slate-800 rounded-2xl p-5 flex flex-col justify-between hover:border-slate-700 transition-all group"
            >
              <div className="space-y-3">
                <div className="flex items-center justify-between gap-2">
                  <span className="text-[10px] font-bold uppercase tracking-wider text-slate-400">
                    {rec.category}
                  </span>
                  <span className={`text-[10px] font-bold uppercase px-2 py-0.5 rounded-full border ${getPriorityBadge(rec.priority)}`}>
                    Prioridad {rec.priority}
                  </span>
                </div>

                <h4 className="text-sm font-bold text-white group-hover:text-indigo-400 transition-colors">
                  {rec.title}
                </h4>

                <p className="text-xs text-slate-300 leading-relaxed">
                  {rec.action}
                </p>
              </div>

              {rec.impact && (
                <div className="mt-4 pt-3 border-t border-slate-800/80 flex items-center gap-1.5 text-xs font-semibold text-emerald-400">
                  <ArrowUpRight className="w-3.5 h-3.5 shrink-0" />
                  <span>Impacto: {rec.impact}</span>
                </div>
              )}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};
