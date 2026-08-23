import React from 'react';
import type { Grade } from '../types/instagram';

interface GradeBadgeProps {
  grade: Grade;
  score: number;
  title: string;
  size?: 'sm' | 'md' | 'lg' | 'xl';
}

export const GradeBadge: React.FC<GradeBadgeProps> = ({
  grade,
  score,
  title,
  size = 'lg',
}) => {
  const getGradeTheme = (g: Grade) => {
    switch (g) {
      case 'A':
        return {
          bg: 'bg-emerald-500/10 border-emerald-500/30 text-emerald-400',
          badgeBg: 'bg-emerald-500 text-slate-950 shadow-emerald-500/30 shadow-lg',
          ring: 'text-emerald-500',
          gradient: 'from-emerald-400 to-teal-500',
          glow: 'shadow-[0_0_50px_-12px_rgba(16,185,129,0.4)]',
          label: 'Excelente / Perfecta',
        };
      case 'B':
        return {
          bg: 'bg-blue-500/10 border-blue-500/30 text-blue-400',
          badgeBg: 'bg-blue-500 text-white shadow-blue-500/30 shadow-lg',
          ring: 'text-blue-500',
          gradient: 'from-blue-400 to-indigo-500',
          glow: 'shadow-[0_0_50px_-12px_rgba(59,130,246,0.4)]',
          label: 'Buena / Sólida',
        };
      case 'D':
        return {
          bg: 'bg-amber-500/10 border-amber-500/30 text-amber-400',
          badgeBg: 'bg-amber-500 text-slate-950 shadow-amber-500/30 shadow-lg',
          ring: 'text-amber-500',
          gradient: 'from-amber-400 to-orange-500',
          glow: 'shadow-[0_0_50px_-12px_rgba(245,158,11,0.4)]',
          label: 'Decente / A Mejorar',
        };
      case 'F':
      default:
        return {
          bg: 'bg-rose-500/10 border-rose-500/30 text-rose-400',
          badgeBg: 'bg-rose-500 text-white shadow-rose-500/30 shadow-lg',
          ring: 'text-rose-500',
          gradient: 'from-rose-500 to-red-600',
          glow: 'shadow-[0_0_50px_-12px_rgba(239,68,68,0.4)]',
          label: 'Nivel Muy Bajo',
        };
    }
  };

  const theme = getGradeTheme(grade);
  const radius = 42;
  const circumference = 2 * Math.PI * radius;
  const strokeDashoffset = circumference - (score / 100) * circumference;

  if (size === 'sm') {
    return (
      <span className={`inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-bold border ${theme.bg}`}>
        <span className="font-black">{grade}</span>
        <span>• {score} pts</span>
      </span>
    );
  }

  return (
    <div className={`relative flex flex-col items-center justify-center p-6 rounded-3xl border bg-slate-900/80 backdrop-blur-xl ${theme.bg} ${theme.glow}`}>
      <div className="relative w-36 h-36 flex items-center justify-center">
        <svg className="w-full h-full transform -rotate-90" viewBox="0 0 100 100">
          <circle
            cx="50"
            cy="50"
            r={radius}
            className="stroke-slate-800"
            strokeWidth="8"
            fill="transparent"
          />
          <circle
            cx="50"
            cy="50"
            r={radius}
            className={`transition-all duration-1000 ease-out ${theme.ring}`}
            strokeWidth="8"
            strokeDasharray={circumference}
            strokeDashoffset={strokeDashoffset}
            strokeLinecap="round"
            fill="transparent"
          />
        </svg>

        <div className="absolute inset-0 flex flex-col items-center justify-center">
          <span className={`text-5xl font-black tracking-tight bg-gradient-to-br ${theme.gradient} bg-clip-text text-transparent`}>
            {grade}
          </span>
          <span className="text-xs font-semibold text-slate-400 mt-0.5">
            {score}/100
          </span>
        </div>
      </div>

      <div className="mt-3 text-center">
        <span className={`text-xs font-bold tracking-wider uppercase px-2.5 py-1 rounded-full border ${theme.bg}`}>
          {title || theme.label}
        </span>
      </div>
    </div>
  );
};
