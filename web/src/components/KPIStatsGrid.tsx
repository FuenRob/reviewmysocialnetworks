import React from 'react';
import type { AccountReport } from '../types/instagram';
import { Users, Heart, Calendar, Zap, Award, Flame } from 'lucide-react';

interface Props {
  report: AccountReport;
}

export const KPIStatsGrid: React.FC<Props> = ({ report }) => {
  const { profile, engagement_metrics, cadence_metrics, growth_metrics } = report;

  const formatNumber = (num: number) => {
    if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M';
    if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
    return num.toLocaleString();
  };

  const formatDayName = (day: string) => {
    if (!day) return 'Sin datos';
    const trimmed = day.trim();
    if (trimmed.toLowerCase().endsWith('s')) {
      return trimmed;
    }
    return `${trimmed}s`;
  };

  const kpis = [
    {
      label: 'Seguidores',
      value: formatNumber(profile.followers_count),
      subtitle: `${formatNumber(profile.follows_count)} seguidos`,
      icon: Users,
      color: 'from-blue-500/20 to-indigo-500/20 text-blue-400 border-blue-500/30',
    },
    {
      label: 'Tasa de Interacción Media',
      value: `${engagement_metrics.average_engagement_rate}%`,
      subtitle: engagement_metrics.benchmark_comparison,
      icon: Zap,
      color: 'from-emerald-500/20 to-teal-500/20 text-emerald-400 border-emerald-500/30',
      highlight: true,
    },
    {
      label: 'Interacciones Promedio',
      value: `${formatNumber(engagement_metrics.average_likes)} likes`,
      subtitle: `${formatNumber(engagement_metrics.average_comments)} comentarios / post`,
      icon: Heart,
      color: 'from-rose-500/20 to-pink-500/20 text-rose-400 border-rose-500/30',
    },
    {
      label: 'Frecuencia de Publicación',
      value: `~${cadence_metrics.estimated_posts_per_week} posts/sem`,
      subtitle: cadence_metrics.cadence_status,
      icon: Calendar,
      color: 'from-purple-500/20 to-violet-500/20 text-purple-400 border-purple-500/30',
    },
    {
      label: 'Mejor Momento para Publicar',
      value: formatDayName(cadence_metrics.best_posting_day),
      subtitle: `A las ${String(cadence_metrics.best_posting_hour).padStart(2, '0')}:00h (Mayor interacción)`,
      icon: Flame,
      color: 'from-amber-500/20 to-orange-500/20 text-amber-400 border-amber-500/30',
    },
    {
      label: 'Ratio de Audiencia',
      value: `${growth_metrics.follower_to_following_ratio}:1`,
      subtitle: growth_metrics.audience_health_status,
      icon: Award,
      color: 'from-cyan-500/20 to-sky-500/20 text-cyan-400 border-cyan-500/30',
    },
  ];

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6">
      {kpis.map((kpi, idx) => {
        const Icon = kpi.icon;
        return (
          <div
            key={idx}
            className={`relative p-5 sm:p-6 rounded-3xl border bg-slate-900/80 backdrop-blur-xl transition-all duration-300 hover:scale-[1.02] hover:shadow-xl ${kpi.color.split(' ').slice(2).join(' ')}`}
          >
            <div className="flex items-start justify-between">
              <div className="space-y-1">
                <span className="text-xs font-semibold uppercase tracking-wider text-slate-400">
                  {kpi.label}
                </span>
                <div className="text-2xl sm:text-3xl font-black text-white tracking-tight">
                  {kpi.value}
                </div>
              </div>
              <div className={`p-3 rounded-2xl bg-gradient-to-br ${kpi.color.split(' ').slice(0, 2).join(' ')} border border-white/10`}>
                <Icon className={`w-6 h-6 ${kpi.color.split(' ')[2]}`} />
              </div>
            </div>
            <div className="mt-3 text-xs text-slate-400 truncate">
              {kpi.subtitle}
            </div>
          </div>
        );
      })}
    </div>
  );
};
