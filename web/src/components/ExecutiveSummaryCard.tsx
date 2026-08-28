import React from 'react';
import type { AccountReport } from '../types/instagram';
import { GradeBadge } from './GradeBadge';
import { Sparkles, Activity, Clock, Layers, Users, ExternalLink } from 'lucide-react';

interface Props {
  report: AccountReport;
}

export const ExecutiveSummaryCard: React.FC<Props> = ({ report }) => {
  const { profile, sub_scores, overall_grade, overall_score, grade_title, executive_summary } = report;

  const subScoreItems = [
    {
      label: 'Tasa de Interacción (Engagement)',
      weight: '35%',
      score: sub_scores.engagement_score,
      icon: Activity,
      color: 'progress-emerald',
      textColor: 'text-emerald-400',
    },
    {
      label: 'Cadencia y Regularidad de Publicación',
      weight: '25%',
      score: sub_scores.consistency_score,
      icon: Clock,
      color: 'progress-blue',
      textColor: 'text-blue-400',
    },
    {
      label: 'Estrategia y Diversidad de Formatos',
      weight: '20%',
      score: sub_scores.content_mix_score,
      icon: Layers,
      color: 'progress-purple',
      textColor: 'text-purple-400',
    },
    {
      label: 'Salud de Audiencia y Autoridad',
      weight: '20%',
      score: sub_scores.audience_health_score,
      icon: Users,
      color: 'progress-amber',
      textColor: 'text-amber-400',
    },
  ];

  return (
    <div className="bg-slate-900/90 border border-slate-800 rounded-3xl p-6 sm:p-8 backdrop-blur-xl shadow-2xl relative overflow-hidden">
      <div className="absolute -right-20 -top-20 w-80 h-80 bg-indigo-500/10 rounded-full blur-3xl pointer-events-none" />

      <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-center">
        <div className="lg:col-span-8 space-y-6">
          <div className="flex flex-col sm:flex-row items-start sm:items-center gap-4">
            <div className="relative">
              <img
                src={profile.profile_picture_url || 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=150&auto=format&fit=crop&q=80'}
                alt={profile.username}
                className="w-20 h-20 rounded-2xl object-cover ring-2 ring-indigo-500/30 p-0.5 shadow-lg bg-slate-800"
              />
              <div className="absolute -bottom-1 -right-1 bg-gradient-to-tr from-yellow-500 via-pink-500 to-purple-600 rounded-full p-1 shadow-md">
                <div className="w-2 h-2 bg-white rounded-full" />
              </div>
            </div>

            <div className="space-y-1">
              <div className="flex flex-wrap items-center gap-2">
                <h1 className="text-2xl sm:text-3xl font-black text-white tracking-tight">
                  @{profile.username}
                </h1>
                {profile.account_type && (
                  <span className="px-2.5 py-0.5 rounded-full text-xs font-semibold bg-slate-800 text-slate-300 border border-slate-700">
                    {profile.account_type}
                  </span>
                )}
              </div>
              {profile.name && (
                <p className="text-sm font-medium text-slate-400">{profile.name}</p>
              )}
              {profile.website && (
                <a
                  href={profile.website}
                  target="_blank"
                  rel="noreferrer"
                  className="inline-flex items-center gap-1 text-xs text-indigo-400 hover:text-indigo-300 transition-colors"
                >
                  <ExternalLink className="w-3 h-3" />
                  {profile.website.replace(/^https?:\/\//, '')}
                </a>
              )}
            </div>
          </div>

          {profile.biography && (
            <div className="bg-slate-950/60 border border-slate-800/80 rounded-2xl p-3.5 text-xs sm:text-sm text-slate-300 whitespace-pre-line leading-relaxed">
              {profile.biography}
            </div>
          )}

          <div className="bg-indigo-950/30 border border-indigo-500/20 rounded-2xl p-4 sm:p-5 relative">
            <div className="flex items-center gap-2 text-indigo-400 font-bold text-xs uppercase tracking-wider mb-2">
              <Sparkles className="w-4 h-4" />
              Diagnóstico Ejecutivo del Algoritmo
            </div>
            <p className="text-sm sm:text-base text-slate-200 leading-relaxed">
              {executive_summary}
            </p>
          </div>

          <div className="space-y-3 pt-2">
            <h4 className="text-xs font-bold uppercase tracking-wider text-slate-400">
              Desglose de Factores Evaluados
            </h4>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3.5">
              {subScoreItems.map((item, idx) => {
                const Icon = item.icon;
                return (
                  <div
                    key={idx}
                    className="bg-slate-950/50 border border-slate-800/60 rounded-2xl p-3.5 space-y-2 hover:border-slate-700/80 transition-all"
                  >
                    <div className="flex items-center justify-between text-xs">
                      <div className="flex items-center gap-2 text-slate-300 font-medium truncate">
                        <Icon className={`w-4 h-4 ${item.textColor}`} />
                        <span className="truncate">{item.label}</span>
                      </div>
                      <span className="font-bold text-white shrink-0">
                        {item.score}<span className="text-slate-500 font-normal">/100</span>
                      </span>
                    </div>

                    <progress
                      className={`score-progress ${item.color}`}
                      value={Math.max(5, item.score)}
                      max={100}
                      aria-label={`${item.label}: ${item.score} de 100`}
                    />
                    <div className="text-[10px] text-slate-500 text-right">
                      Peso: {item.weight}
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        </div>

        <div className="lg:col-span-4 flex flex-col items-center justify-center">
          <GradeBadge
            grade={overall_grade}
            score={overall_score}
            title={grade_title}
            size="lg"
          />
        </div>
      </div>
    </div>
  );
};
