import React from 'react';
import { RefreshCw, Printer } from 'lucide-react';
import { InstagramIcon } from './InstagramIcon';
import type { AccountReport } from '../types/instagram';

interface Props {
  onReset: () => void;
  report: AccountReport | null;
  onLoadDemo: (tier: 'A' | 'B' | 'D' | 'F', platform?: 'instagram' | 'tiktok') => void;
}

export const Navbar: React.FC<Props> = ({
  onReset,
  report,
  onLoadDemo,
}) => {
  const handlePrint = () => {
    window.print();
  };

  return (
    <header className="sticky top-0 z-40 w-full border-b border-slate-800/80 bg-slate-950/80 backdrop-blur-xl no-print">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
        <div
          onClick={onReset}
          className="flex items-center gap-3 cursor-pointer group"
        >
          <div className="w-10 h-10 rounded-2xl bg-gradient-to-tr from-yellow-500 via-pink-500 to-purple-600 p-0.5 shadow-lg shadow-pink-500/20 group-hover:scale-105 transition-transform">
            <div className="w-full h-full bg-slate-950 rounded-[14px] flex items-center justify-center">
              <InstagramIcon className="w-5 h-5 text-white" />
            </div>
          </div>
          <div>
            <span className="text-base font-black tracking-tight text-white flex items-center gap-1.5">
              ReviewMy<span className="bg-gradient-to-r from-pink-500 to-purple-400 bg-clip-text text-transparent">SocialNetworks</span>
            </span>
            <span className="text-[10px] text-slate-400 block -mt-1 font-medium">
              Instagram + TikTok Audit Engine
            </span>
          </div>
        </div>

        <div className="flex items-center gap-2 sm:gap-3">
          {report && (
            <div className="hidden md:flex items-center gap-1 bg-slate-900 border border-slate-800 rounded-2xl p-1">
              <span className="text-[10px] text-slate-500 font-bold px-2 uppercase">Demo:</span>
              {(['A', 'B', 'D', 'F'] as const).map((tier) => (
                <button
                  key={tier}
                  onClick={() => onLoadDemo(tier, report.platform)}
                  className={`w-6 h-6 rounded-xl text-xs font-black transition-all ${
                    report.overall_grade === tier
                      ? 'bg-indigo-600 text-white shadow-sm'
                      : 'text-slate-400 hover:text-white hover:bg-slate-800'
                  }`}
                >
                  {tier}
                </button>
              ))}
            </div>
          )}

          {report && (
            <button
              onClick={handlePrint}
              className="inline-flex items-center gap-1.5 px-3.5 py-2 rounded-2xl bg-slate-900 hover:bg-slate-800 text-slate-200 border border-slate-800 text-xs font-semibold transition-colors"
            >
              <Printer className="w-3.5 h-3.5 text-indigo-400" />
              <span className="hidden sm:inline">Exportar / Imprimir</span>
            </button>
          )}

          {report && (
            <button
              onClick={onReset}
              className="inline-flex items-center gap-1.5 px-3.5 py-2 rounded-2xl bg-slate-900 hover:bg-slate-800 text-slate-200 border border-slate-800 text-xs font-semibold transition-colors"
            >
              <RefreshCw className="w-3.5 h-3.5 text-pink-400" />
              <span className="hidden sm:inline">Nuevo Análisis</span>
            </button>
          )}
        </div>
      </div>
    </header>
  );
};
