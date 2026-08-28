import React from 'react';
import type { AccountReport, MediaAnalysisItem } from '../types/instagram';
import { BarChart3, PieChart, CalendarDays } from 'lucide-react';

interface Props { report: AccountReport; }

const width = 800;
const height = 280;

function TimelineChart({ media }: { media: MediaAnalysisItem[] }) {
  if (media.length === 0) return <div className="h-full flex items-center justify-center text-slate-500 text-sm">Sin datos suficientes para graficar</div>;
  const ordered = [...media].reverse();
  const left = 56, right = 56, top = 20, bottom = 48;
  const plotWidth = width - left - right, plotHeight = height - top - bottom;
  const maxLikes = Math.max(1, ...ordered.map((item) => item.like_count));
  const maxEngagement = Math.max(1, ...ordered.map((item) => item.engagement_rate));
  const step = ordered.length > 1 ? plotWidth / (ordered.length - 1) : plotWidth;
  const barWidth = Math.min(38, plotWidth / ordered.length * 0.5);
  const x = (index: number) => ordered.length === 1 ? left + plotWidth / 2 : left + index * step;
  const likesY = (value: number) => top + plotHeight - value / maxLikes * plotHeight;
  const engagementY = (value: number) => top + plotHeight - value / maxEngagement * plotHeight;
  const points = ordered.map((item, index) => `${x(index)},${engagementY(item.engagement_rate)}`).join(' ');

  return (
    <svg viewBox={`0 0 ${width} ${height}`} className="w-full h-full" role="img" aria-label="Evolución de likes y tasa de interacción">
      {[0, 0.25, 0.5, 0.75, 1].map((ratio) => <line key={ratio} x1={left} x2={width - right} y1={top + plotHeight * ratio} y2={top + plotHeight * ratio} stroke="#334155" strokeOpacity="0.45" />)}
      {ordered.map((item, index) => {
        const date = new Date(item.timestamp);
        const barY = likesY(item.like_count);
        return <g key={item.id || index}>
          <rect x={x(index) - barWidth / 2} y={barY} width={barWidth} height={top + plotHeight - barY} rx="7" fill="#6366f1" fillOpacity="0.48"><title>{`${item.like_count} likes · ${item.engagement_rate}% engagement`}</title></rect>
          <text x={x(index)} y={height - 19} textAnchor="middle" fill="#94a3b8" fontSize="12">{date.getDate()}/{date.getMonth() + 1}</text>
        </g>;
      })}
      <polyline points={points} fill="none" stroke="#10b981" strokeWidth="4" strokeLinejoin="round" strokeLinecap="round" />
      {ordered.map((item, index) => <circle key={`point-${item.id || index}`} cx={x(index)} cy={engagementY(item.engagement_rate)} r="5" fill="#10b981" stroke="#f8fafc" strokeWidth="2"><title>{`${item.engagement_rate}% engagement`}</title></circle>)}
      <text x="12" y="18" fill="#818cf8" fontSize="12">Likes</text>
      <text x={width - 8} y="18" textAnchor="end" fill="#10b981" fontSize="12">Engagement %</text>
    </svg>
  );
}

function FormatDonut({ report }: { report: AccountReport }) {
  const values = [report.content_metrics.carousel_count, report.content_metrics.video_count, report.content_metrics.image_count];
  const colors = ['#a855f7', '#ec4899', '#3b82f6'];
  const labels = ['Carruseles', 'Reels', 'Fotos'];
  const total = values.reduce((sum, value) => sum + value, 0);
  const radius = 78, circumference = 2 * Math.PI * radius;
  let offset = 0;
  return <svg viewBox="0 0 240 240" className="w-full h-full" role="img" aria-label="Distribución de formatos de publicación">
    <circle cx="120" cy="120" r={radius} fill="none" stroke="#1e293b" strokeWidth="28" />
    {total > 0 && values.map((value, index) => {
      const fraction = value / total, dashOffset = -offset * circumference;
      offset += fraction;
      return <circle key={colors[index]} cx="120" cy="120" r={radius} fill="none" stroke={colors[index]} strokeWidth="28" strokeDasharray={`${fraction * circumference} ${circumference}`} strokeDashoffset={dashOffset} transform="rotate(-90 120 120)"><title>{`${labels[index]}: ${value}`}</title></circle>;
    })}
    <text x="120" y="116" textAnchor="middle" fill="#f8fafc" fontSize="30" fontWeight="700">{total}</text>
    <text x="120" y="140" textAnchor="middle" fill="#94a3b8" fontSize="13">publicaciones</text>
  </svg>;
}

function DayChart({ report }: { report: AccountReport }) {
  const labels = ['Lun', 'Mar', 'Mié', 'Jue', 'Vie', 'Sáb', 'Dom'];
  const fullLabels = ['Lunes', 'Martes', 'Miércoles', 'Jueves', 'Viernes', 'Sábado', 'Domingo'];
  const values = fullLabels.map((day) => report.cadence_metrics.day_distribution[day] || 0);
  const max = Math.max(1, ...values), left = 36, top = 18, bottom = 42;
  const plotHeight = height - top - bottom, columnWidth = (width - left * 2) / values.length;
  return <svg viewBox={`0 0 ${width} ${height}`} className="w-full h-full" role="img" aria-label="Publicaciones distribuidas por día de la semana">
    {[0, 0.5, 1].map((ratio) => <line key={ratio} x1={left} x2={width - left} y1={top + plotHeight * ratio} y2={top + plotHeight * ratio} stroke="#334155" strokeOpacity="0.45" />)}
    {values.map((value, index) => {
      const barHeight = value / max * plotHeight;
      const best = fullLabels[index] === report.cadence_metrics.best_posting_day;
      return <g key={fullLabels[index]}>
        <rect x={left + index * columnWidth + columnWidth * 0.2} y={top + plotHeight - barHeight} width={columnWidth * 0.6} height={barHeight} rx="8" fill={best ? '#f59e0b' : '#6366f1'} fillOpacity={best ? 0.9 : 0.55}><title>{`${fullLabels[index]}: ${value} publicaciones${best ? ' · mejor día' : ''}`}</title></rect>
        <text x={left + index * columnWidth + columnWidth / 2} y={height - 17} textAnchor="middle" fill={best ? '#fbbf24' : '#94a3b8'} fontSize="13">{labels[index]}</text>
      </g>;
    })}
  </svg>;
}

export const EngagementCharts: React.FC<Props> = ({ report }) => {
  const { content_metrics, cadence_metrics } = report;
  const bestFormat = content_metrics.best_performing_type === 'CAROUSEL_ALBUM' ? 'Carruseles' : content_metrics.best_performing_type === 'VIDEO' ? 'Reels' : 'Fotos';
  const summaries: Array<[string, number, number, string]> = [
    ['Carruseles', content_metrics.carousel_percentage, content_metrics.average_by_format?.CAROUSEL_ALBUM?.average_engagement_rate || 0, 'text-purple-400'],
    ['Reels', content_metrics.video_percentage, content_metrics.average_by_format?.VIDEO?.average_engagement_rate || 0, 'text-pink-400'],
    ['Fotos', content_metrics.image_percentage, content_metrics.average_by_format?.IMAGE?.average_engagement_rate || 0, 'text-blue-400'],
  ];

  return <div className="space-y-6">
    <section className="bg-slate-900/80 border border-slate-800 rounded-3xl p-6 backdrop-blur-xl">
      <div className="flex items-center gap-3 mb-6"><div className="p-2.5 rounded-2xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-400"><BarChart3 className="w-5 h-5" /></div><div><h3 className="text-base sm:text-lg font-bold text-white">Rendimiento de Interacción por Publicación</h3><p className="text-xs text-slate-400">Evolución de Likes y porcentaje de engagement a lo largo del tiempo</p></div></div>
      <div className="h-72 w-full"><TimelineChart media={report.media_analysis} /></div>
    </section>
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <section className="bg-slate-900/80 border border-slate-800 rounded-3xl p-6 backdrop-blur-xl">
        <div className="flex items-center gap-3 mb-4"><div className="p-2.5 rounded-2xl bg-purple-500/10 border border-purple-500/20 text-purple-400"><PieChart className="w-5 h-5" /></div><div><h3 className="text-base sm:text-lg font-bold text-white">Distribución de Formatos</h3><p className="text-xs text-slate-400">Mejor formato: <span className="text-purple-400 font-semibold">{bestFormat}</span></p></div></div>
        <div className="h-56 w-full"><FormatDonut report={report} /></div>
        <div className="mt-4 pt-4 border-t border-slate-800 grid grid-cols-3 gap-2 text-center text-xs">{summaries.map(([label, percentage, engagement, color]) => <div key={label} className="bg-slate-950/60 p-2.5 rounded-2xl border border-slate-800/80"><span className={`${color} font-bold block`}>{label}</span><span className="text-slate-200 font-semibold">{percentage}%</span><span className="text-[10px] text-slate-500 block">{engagement}% eng</span></div>)}</div>
      </section>
      <section className="bg-slate-900/80 border border-slate-800 rounded-3xl p-6 backdrop-blur-xl">
        <div className="flex items-center gap-3 mb-4"><div className="p-2.5 rounded-2xl bg-amber-500/10 border border-amber-500/20 text-amber-400"><CalendarDays className="w-5 h-5" /></div><div><h3 className="text-base sm:text-lg font-bold text-white">Distribución por Día de Publicación</h3><p className="text-xs text-slate-400">Día estelar: <span className="text-amber-400 font-semibold">{cadence_metrics.best_posting_day}</span></p></div></div>
        <div className="h-56 w-full"><DayChart report={report} /></div>
        <div className="mt-4 pt-4 border-t border-slate-800 flex items-center justify-between text-xs text-slate-400 bg-slate-950/40 p-3 rounded-2xl"><span>Hora dorada estimada:</span><span className="font-bold text-amber-400 text-sm">{String(cadence_metrics.best_posting_hour).padStart(2, '0')}:00 hrs</span></div>
      </section>
    </div>
  </div>;
};
