import React from 'react';
import type { AccountReport } from '../types/instagram';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  Filler,
} from 'chart.js';
import { Line, Doughnut, Bar } from 'react-chartjs-2';
import { BarChart3, PieChart, CalendarDays } from 'lucide-react';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  Filler
);

interface Props {
  report: AccountReport;
}

export const EngagementCharts: React.FC<Props> = ({ report }) => {
  const { media_analysis, content_metrics, cadence_metrics } = report;

  const chronologicalMedia = [...media_analysis].reverse();

  const timelineLabels = chronologicalMedia.map((m) => {
    const d = new Date(m.timestamp);
    return `${d.getDate()}/${d.getMonth() + 1}`;
  });

  const timelineData = {
    labels: timelineLabels,
    datasets: [
      {
        type: 'line' as const,
        label: 'Tasa Engagement (%)',
        data: chronologicalMedia.map((m) => m.engagement_rate),
        borderColor: '#10B981',
        backgroundColor: 'rgba(16, 185, 129, 0.1)',
        fill: true,
        tension: 0.4,
        yAxisID: 'y1',
        pointBackgroundColor: '#10B981',
        pointBorderColor: '#fff',
        pointHoverRadius: 6,
      },
      {
        type: 'bar' as const,
        label: 'Likes',
        data: chronologicalMedia.map((m) => m.like_count),
        backgroundColor: 'rgba(99, 102, 241, 0.4)',
        borderColor: '#6366F1',
        borderWidth: 1,
        borderRadius: 8,
        yAxisID: 'y',
      },
    ],
  };

  const timelineOptions = {
    responsive: true,
    maintainAspectRatio: false,
    interaction: {
      mode: 'index' as const,
      intersect: false,
    },
    plugins: {
      legend: {
        position: 'top' as const,
        labels: {
          color: '#94A3B8',
          font: { family: 'sans-serif', size: 12 },
        },
      },
      tooltip: {
        backgroundColor: '#0F172A',
        borderColor: '#334155',
        borderWidth: 1,
        padding: 12,
        titleColor: '#F8FAFC',
        bodyColor: '#CBD5E1',
      },
    },
    scales: {
      x: {
        grid: { color: 'rgba(51, 65, 85, 0.3)' },
        ticks: { color: '#94A3B8' },
      },
      y: {
        type: 'linear' as const,
        display: true,
        position: 'left' as const,
        title: { display: true, text: 'Likes', color: '#94A3B8' },
        grid: { color: 'rgba(51, 65, 85, 0.3)' },
        ticks: { color: '#94A3B8' },
      },
      y1: {
        type: 'linear' as const,
        display: true,
        position: 'right' as const,
        title: { display: true, text: 'Engagement (%)', color: '#10B981' },
        grid: { drawOnChartArea: false },
        ticks: {
          color: '#10B981',
          callback: (value: any) => `${value}%`,
        },
      },
    },
  };

  const formatData = {
    labels: ['Carruseles', 'Reels / Vídeos', 'Fotos Individuales'],
    datasets: [
      {
        data: [
          content_metrics.carousel_count,
          content_metrics.video_count,
          content_metrics.image_count,
        ],
        backgroundColor: [
          'rgba(168, 85, 247, 0.8)',
          'rgba(236, 72, 153, 0.8)',
          'rgba(59, 130, 246, 0.8)',
        ],
        borderColor: '#0F172A',
        borderWidth: 3,
        hoverOffset: 6,
      },
    ],
  };

  const formatOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'bottom' as const,
        labels: { color: '#94A3B8', boxWidth: 12, padding: 16 },
      },
    },
    cutout: '70%',
  };

  const dayLabels = ['Lunes', 'Martes', 'Miércoles', 'Jueves', 'Viernes', 'Sábado', 'Domingo'];
  const dayCounts = dayLabels.map((d) => cadence_metrics.day_distribution[d] || 0);

  const dayData = {
    labels: dayLabels,
    datasets: [
      {
        label: 'Publicaciones',
        data: dayCounts,
        backgroundColor: dayLabels.map((d) =>
          d === cadence_metrics.best_posting_day
            ? 'rgba(245, 158, 11, 0.85)'
            : 'rgba(99, 102, 241, 0.4)'
        ),
        borderColor: dayLabels.map((d) =>
          d === cadence_metrics.best_posting_day ? '#F59E0B' : '#6366F1'
        ),
        borderWidth: 1,
        borderRadius: 8,
      },
    ],
  };

  const dayOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: { display: false },
      tooltip: {
        callbacks: {
          afterLabel: (ctx: any) =>
            ctx.label === cadence_metrics.best_posting_day
              ? '⭐ Mejor día según interacción histórica'
              : '',
        },
      },
    },
    scales: {
      x: {
        grid: { display: false },
        ticks: { color: '#94A3B8' },
      },
      y: {
        beginAtZero: true,
        ticks: { stepSize: 1, color: '#94A3B8' },
        grid: { color: 'rgba(51, 65, 85, 0.3)' },
      },
    },
  };

  return (
    <div className="space-y-6">
      <div className="bg-slate-900/80 border border-slate-800 rounded-3xl p-6 backdrop-blur-xl">
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center gap-3">
            <div className="p-2.5 rounded-2xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-400">
              <BarChart3 className="w-5 h-5" />
            </div>
            <div>
              <h3 className="text-base sm:text-lg font-bold text-white">
                Rendimiento de Interacción por Publicación
              </h3>
              <p className="text-xs text-slate-400">
                Evolución de Likes y porcentaje de engagement a lo largo del tiempo
              </p>
            </div>
          </div>
        </div>

        <div className="h-72 w-full">
          {chronologicalMedia.length > 0 ? (
            <Line data={timelineData as any} options={timelineOptions} />
          ) : (
            <div className="h-full flex items-center justify-center text-slate-500 text-sm">
              Sin datos suficientes para graficar
            </div>
          )}
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-slate-900/80 border border-slate-800 rounded-3xl p-6 backdrop-blur-xl flex flex-col justify-between">
          <div>
            <div className="flex items-center gap-3 mb-4">
              <div className="p-2.5 rounded-2xl bg-purple-500/10 border border-purple-500/20 text-purple-400">
                <PieChart className="w-5 h-5" />
              </div>
              <div>
                <h3 className="text-base sm:text-lg font-bold text-white">
                  Distribución de Formatos
                </h3>
                <p className="text-xs text-slate-400">
                  Mejor formato:{' '}
                  <span className="text-purple-400 font-semibold">
                    {content_metrics.best_performing_type === 'CAROUSEL_ALBUM'
                      ? 'Carruseles'
                      : content_metrics.best_performing_type === 'VIDEO'
                      ? 'Reels'
                      : 'Fotos'}
                  </span>
                </p>
              </div>
            </div>

            <div className="h-56 w-full relative flex items-center justify-center">
              <Doughnut data={formatData} options={formatOptions} />
            </div>
          </div>

          <div className="mt-4 pt-4 border-t border-slate-800 grid grid-cols-3 gap-2 text-center text-xs">
            <div className="bg-slate-950/60 p-2.5 rounded-2xl border border-slate-800/80">
              <span className="text-purple-400 font-bold block">Carruseles</span>
              <span className="text-slate-200 font-semibold">
                {content_metrics.carousel_percentage}%
              </span>
              <span className="text-[10px] text-slate-500 block">
                {content_metrics.average_by_format?.CAROUSEL_ALBUM?.average_engagement_rate || 0}% eng
              </span>
            </div>
            <div className="bg-slate-950/60 p-2.5 rounded-2xl border border-slate-800/80">
              <span className="text-pink-400 font-bold block">Reels</span>
              <span className="text-slate-200 font-semibold">
                {content_metrics.video_percentage}%
              </span>
              <span className="text-[10px] text-slate-500 block">
                {content_metrics.average_by_format?.VIDEO?.average_engagement_rate || 0}% eng
              </span>
            </div>
            <div className="bg-slate-950/60 p-2.5 rounded-2xl border border-slate-800/80">
              <span className="text-blue-400 font-bold block">Fotos</span>
              <span className="text-slate-200 font-semibold">
                {content_metrics.image_percentage}%
              </span>
              <span className="text-[10px] text-slate-500 block">
                {content_metrics.average_by_format?.IMAGE?.average_engagement_rate || 0}% eng
              </span>
            </div>
          </div>
        </div>

        <div className="bg-slate-900/80 border border-slate-800 rounded-3xl p-6 backdrop-blur-xl flex flex-col justify-between">
          <div>
            <div className="flex items-center gap-3 mb-4">
              <div className="p-2.5 rounded-2xl bg-amber-500/10 border border-amber-500/20 text-amber-400">
                <CalendarDays className="w-5 h-5" />
              </div>
              <div>
                <h3 className="text-base sm:text-lg font-bold text-white">
                  Distribución por Día de Publicación
                </h3>
                <p className="text-xs text-slate-400">
                  Día estelar:{' '}
                  <span className="text-amber-400 font-semibold">
                    {cadence_metrics.best_posting_day}
                  </span>
                </p>
              </div>
            </div>

            <div className="h-56 w-full">
              <Bar data={dayData} options={dayOptions} />
            </div>
          </div>

          <div className="mt-4 pt-4 border-t border-slate-800 flex items-center justify-between text-xs text-slate-400 bg-slate-950/40 p-3 rounded-2xl">
            <span>Hora dorada estimada:</span>
            <span className="font-bold text-amber-400 text-sm">
              {String(cadence_metrics.best_posting_hour).padStart(2, '0')}:00 hrs
            </span>
          </div>
        </div>
      </div>
    </div>
  );
};
