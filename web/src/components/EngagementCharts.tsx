import React from 'react';
import type { AccountReport, MediaAnalysisItem } from '../types/instagram';
import { BarChart3, PieChart, CalendarDays, Eye, Share2, Gauge, Heart, Clock3, TrendingUp, MessageCircle } from 'lucide-react';

interface Props { report: AccountReport; }

const width = 800;
const height = 280;

function TimelineChart({ media, isTikTok }: { media: MediaAnalysisItem[]; isTikTok: boolean }) {
  if (media.length === 0) return <div className="h-full flex items-center justify-center text-slate-500 text-sm">Sin datos suficientes para graficar</div>;
  const ordered = [...media].reverse();
  const chartHeight = 360;
  const left = 76, right = 76, top = 48, bottom = 68;
  const chartWidth = Math.max(900, left + right + ordered.length * 88);
  const plotWidth = chartWidth - left - right, plotHeight = chartHeight - top - bottom;
  const maxPrimary = Math.max(1, ...ordered.map((item) => isTikTok ? item.view_count || 0 : item.like_count));
  const engagementFor = (item: MediaAnalysisItem) => isTikTok ? item.view_engagement_rate || 0 : item.engagement_rate;
  const maxEngagement = Math.max(1, ...ordered.map(engagementFor));
  const step = ordered.length > 1 ? plotWidth / (ordered.length - 1) : plotWidth;
  const barWidth = Math.min(46, plotWidth / ordered.length * 0.52);
  const x = (index: number) => ordered.length === 1 ? left + plotWidth / 2 : left + index * step;
  const primaryY = (value: number) => top + plotHeight - value / maxPrimary * plotHeight;
  const engagementY = (value: number) => top + plotHeight - value / maxEngagement * plotHeight;
  const points = ordered.map((item, index) => `${x(index)},${engagementY(engagementFor(item))}`).join(' ');
  const compactNumber = (value: number) => value >= 1_000_000 ? `${(value / 1_000_000).toFixed(1)}M` : value >= 1_000 ? `${(value / 1_000).toFixed(1)}K` : Math.round(value).toLocaleString('es-ES');

  return (
    <svg
      viewBox={`0 0 ${chartWidth} ${chartHeight}`}
      className="block h-auto min-w-full max-w-none"
      style={{ width: `${chartWidth}px` }}
      role="img"
      aria-label={`Evolución de ${isTikTok ? 'visualizaciones' : 'likes'} y tasa de interacción`}
    >
      {[0, 0.25, 0.5, 0.75, 1].map((ratio) => {
        const y = top + plotHeight * ratio;
        return <g key={ratio}>
          <line x1={left} x2={chartWidth - right} y1={y} y2={y} stroke="#334155" strokeOpacity="0.6" strokeDasharray={ratio === 1 ? undefined : '4 6'} />
          <text x={left - 12} y={y + 4} textAnchor="end" fill="#94a3b8" fontSize="11">{compactNumber(maxPrimary * (1 - ratio))}</text>
          <text x={chartWidth - right + 12} y={y + 4} textAnchor="start" fill="#6ee7b7" fontSize="11">{(maxEngagement * (1 - ratio)).toFixed(1)}%</text>
        </g>;
      })}
      {ordered.map((item, index) => {
        const date = new Date(item.timestamp);
        const primaryValue = isTikTok ? item.view_count || 0 : item.like_count;
        const barY = primaryY(primaryValue);
        return <g key={item.id || index}>
          <rect x={x(index) - barWidth / 2} y={barY} width={barWidth} height={top + plotHeight - barY} rx="8" fill="#6366f1" fillOpacity="0.58"><title>{`${primaryValue.toLocaleString('es-ES')} ${isTikTok ? 'visualizaciones' : 'likes'} · ${engagementFor(item)}% engagement`}</title></rect>
          {ordered.length <= 12 && <text x={x(index)} y={Math.max(top + 12, barY - 9)} textAnchor="middle" fill="#a5b4fc" fontSize="11" fontWeight="600">{compactNumber(primaryValue)}</text>}
          <text x={x(index)} y={chartHeight - 29} textAnchor="middle" fill="#cbd5e1" fontSize="12" fontWeight="600">{date.getDate().toString().padStart(2, '0')}</text>
          <text x={x(index)} y={chartHeight - 13} textAnchor="middle" fill="#64748b" fontSize="10">{date.toLocaleDateString('es-ES', { month: 'short' }).replace('.', '')}</text>
        </g>;
      })}
      <polyline points={points} fill="none" stroke="#10b981" strokeWidth="4" strokeLinejoin="round" strokeLinecap="round" />
      {ordered.map((item, index) => <circle key={`point-${item.id || index}`} cx={x(index)} cy={engagementY(engagementFor(item))} r="6" fill="#10b981" stroke="#ecfdf5" strokeWidth="2.5"><title>{`${engagementFor(item)}% engagement`}</title></circle>)}
      <text x={left} y="24" fill="#a5b4fc" fontSize="12" fontWeight="600">{isTikTok ? 'Visualizaciones' : 'Likes'}</text>
      <text x={chartWidth - right} y="24" textAnchor="end" fill="#6ee7b7" fontSize="12" fontWeight="600">Engagement %</text>
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
  const isTikTok = report.platform === 'tiktok';
  const bestFormat = content_metrics.best_performing_type === 'CAROUSEL_ALBUM' ? 'Carruseles' : content_metrics.best_performing_type === 'VIDEO' ? 'Reels' : 'Fotos';
  const summaries: Array<[string, number, number, string]> = [
    ['Carruseles', content_metrics.carousel_percentage, content_metrics.average_by_format?.CAROUSEL_ALBUM?.average_engagement_rate || 0, 'text-purple-400'],
    ['Reels', content_metrics.video_percentage, content_metrics.average_by_format?.VIDEO?.average_engagement_rate || 0, 'text-pink-400'],
    ['Fotos', content_metrics.image_percentage, content_metrics.average_by_format?.IMAGE?.average_engagement_rate || 0, 'text-blue-400'],
  ];
  const compactNumber = (value: number) => value >= 1_000_000 ? `${(value / 1_000_000).toFixed(1)}M` : value >= 1_000 ? `${(value / 1_000).toFixed(1)}K` : value.toLocaleString('es-ES');

  return <div className="space-y-6">
    <section className="overflow-hidden bg-slate-900/80 border border-slate-800 rounded-3xl p-4 sm:p-6 backdrop-blur-xl">
      <div className="flex flex-col gap-4 mb-5 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-3"><div className="shrink-0 p-2.5 rounded-2xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-400"><BarChart3 className="w-5 h-5" /></div><div><h3 className="text-base sm:text-lg font-bold text-white">Rendimiento de Interacción por Publicación</h3><p className="text-xs text-slate-400">Evolución de {isTikTok ? 'visualizaciones' : 'likes'} y porcentaje de engagement a lo largo del tiempo</p></div></div>
        <div className="flex flex-wrap items-center gap-2 pl-12 sm:pl-0 text-[11px] font-medium">
          <span className="inline-flex items-center gap-1.5 rounded-full border border-indigo-500/20 bg-indigo-500/10 px-2.5 py-1 text-indigo-300"><span className="h-2.5 w-2.5 rounded-sm bg-indigo-500/70" />{isTikTok ? 'Visualizaciones' : 'Likes'}</span>
          <span className="inline-flex items-center gap-1.5 rounded-full border border-emerald-500/20 bg-emerald-500/10 px-2.5 py-1 text-emerald-300"><span className="h-0.5 w-3 rounded-full bg-emerald-400" />Engagement</span>
        </div>
      </div>
      <div className="rounded-2xl border border-slate-800/80 bg-slate-950/40 px-2 py-3 sm:px-4 sm:py-4">
        <div className="w-full overflow-x-auto pb-2"><TimelineChart media={report.media_analysis} isTikTok={isTikTok} /></div>
        <p className="mt-1 text-center text-[10px] text-slate-600 sm:hidden">Desliza horizontalmente para ver todas las publicaciones</p>
      </div>
    </section>
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
      {isTikTok ? <section className="bg-slate-900/80 border border-slate-800 rounded-3xl p-6 backdrop-blur-xl">
        <div className="flex items-center gap-3 mb-5"><div className="p-2.5 rounded-2xl bg-cyan-500/10 border border-cyan-500/20 text-cyan-400"><Eye className="w-5 h-5" /></div><div><h3 className="text-base sm:text-lg font-bold text-white">Alcance y Viralidad</h3><p className="text-xs text-slate-400">Distribución y conversión de las visualizaciones</p></div></div>
        <div className="grid grid-cols-2 gap-3 text-center">
          <div className="p-4 rounded-2xl bg-slate-950/60 border border-slate-800"><Eye className="w-4 h-4 text-cyan-400 mx-auto mb-2" /><span className="text-2xl font-black text-white block">{(report.tiktok_metrics?.median_views || 0).toLocaleString()}</span><span className="text-[10px] text-slate-500">mediana de vistas</span></div>
          <div className="p-4 rounded-2xl bg-slate-950/60 border border-slate-800"><Share2 className="w-4 h-4 text-purple-400 mx-auto mb-2" /><span className="text-2xl font-black text-white block">{report.tiktok_metrics?.share_rate || 0}%</span><span className="text-[10px] text-slate-500">tasa de compartidos</span></div>
          <div className="p-4 rounded-2xl bg-slate-950/60 border border-slate-800"><Gauge className="w-4 h-4 text-emerald-400 mx-auto mb-2" /><span className="text-2xl font-black text-white block">{report.tiktok_metrics?.views_per_follower || 0}x</span><span className="text-[10px] text-slate-500">vistas por seguidor</span></div>
          <div className="p-4 rounded-2xl bg-slate-950/60 border border-slate-800"><BarChart3 className="w-4 h-4 text-amber-400 mx-auto mb-2" /><span className="text-2xl font-black text-white block">{report.tiktok_metrics?.viral_videos_count || 0}</span><span className="text-[10px] text-slate-500">vídeos con alcance viral</span></div>
        </div>
      </section> : <section className="bg-slate-900/80 border border-slate-800 rounded-3xl p-6 backdrop-blur-xl">
        <div className="flex items-center gap-3 mb-4"><div className="p-2.5 rounded-2xl bg-purple-500/10 border border-purple-500/20 text-purple-400"><PieChart className="w-5 h-5" /></div><div><h3 className="text-base sm:text-lg font-bold text-white">Distribución de Formatos</h3><p className="text-xs text-slate-400">Mejor formato: <span className="text-purple-400 font-semibold">{bestFormat}</span></p></div></div>
        <div className="h-56 w-full"><FormatDonut report={report} /></div>
        <div className="mt-4 pt-4 border-t border-slate-800 grid grid-cols-3 gap-2 text-center text-xs">{summaries.map(([label, percentage, engagement, color]) => <div key={label} className="bg-slate-950/60 p-2.5 rounded-2xl border border-slate-800/80"><span className={`${color} font-bold block`}>{label}</span><span className="text-slate-200 font-semibold">{percentage}%</span><span className="text-[10px] text-slate-500 block">{engagement}% eng</span></div>)}</div>
      </section>}
      <section className="bg-slate-900/80 border border-slate-800 rounded-3xl p-6 backdrop-blur-xl">
        <div className="flex items-center gap-3 mb-4"><div className="p-2.5 rounded-2xl bg-amber-500/10 border border-amber-500/20 text-amber-400"><CalendarDays className="w-5 h-5" /></div><div><h3 className="text-base sm:text-lg font-bold text-white">Distribución por Día de Publicación</h3><p className="text-xs text-slate-400">Día estelar: <span className="text-amber-400 font-semibold">{cadence_metrics.best_posting_day}</span></p></div></div>
        <div className="h-56 w-full"><DayChart report={report} /></div>
        <div className="mt-4 pt-4 border-t border-slate-800 flex items-center justify-between text-xs text-slate-400 bg-slate-950/40 p-3 rounded-2xl"><span>Hora dorada estimada:</span><span className="font-bold text-amber-400 text-sm">{String(cadence_metrics.best_posting_hour).padStart(2, '0')}:00 hrs</span></div>
      </section>
    </div>
    {isTikTok && report.tiktok_metrics && <section className="bg-slate-900/80 border border-slate-800 rounded-3xl p-6 backdrop-blur-xl">
      <div className="flex items-center gap-3 mb-5"><div className="p-2.5 rounded-2xl bg-indigo-500/10 border border-indigo-500/20 text-indigo-400"><Gauge className="w-5 h-5" /></div><div><h3 className="text-base sm:text-lg font-bold text-white">Resumen Técnico de TikTok</h3><p className="text-xs text-slate-400">Métricas acumuladas del perfil y de los {report.data_coverage.analyzed_posts} vídeos analizados</p></div></div>
      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-3 text-center">
        {[
          ['Likes del perfil', compactNumber(report.profile.likes_count || 0), Heart, 'text-rose-400'],
          ['Vistas analizadas', compactNumber(report.tiktok_metrics.total_views), Eye, 'text-cyan-400'],
          ['Compartidos', compactNumber(report.tiktok_metrics.total_shares), Share2, 'text-purple-400'],
          ['Duración media', `${report.tiktok_metrics.average_duration_seconds}s`, Clock3, 'text-amber-400'],
          ['Tendencia reciente', `${report.growth_metrics.recent_trend_percentage > 0 ? '+' : ''}${report.growth_metrics.recent_trend_percentage}%`, TrendingUp, report.growth_metrics.recent_trend_percentage >= 0 ? 'text-emerald-400' : 'text-rose-400'],
          ['Comentarios / like', `${report.engagement_metrics.comment_to_like_ratio}%`, MessageCircle, 'text-blue-400'],
        ].map(([label, value, Icon, color]: any) => <div key={label} className="p-3.5 rounded-2xl bg-slate-950/60 border border-slate-800"><Icon className={`w-4 h-4 ${color} mx-auto mb-2`} /><span className="text-lg font-black text-white block">{value}</span><span className="text-[10px] text-slate-500">{label}</span></div>)}
      </div>
    </section>}
  </div>;
};
