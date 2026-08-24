import React, { lazy, Suspense, useState, useEffect, useRef } from 'react';
import { getAuthResult, analyzeDemo } from './api/client';
import type { AccountReport } from './types/instagram';
import { Navbar } from './components/Navbar';
import { HeroConnect } from './components/HeroConnect';
import { ExecutiveSummaryCard } from './components/ExecutiveSummaryCard';
import { KPIStatsGrid } from './components/KPIStatsGrid';
import { ReportActionPlan } from './components/ReportActionPlan';
import { MediaGrid } from './components/MediaGrid';
import { PreFooter } from './components/PreFooter';
import { LegalModal, type LegalDocType } from './components/LegalModal';
import { ArrowUp } from 'lucide-react';

const EngagementCharts = lazy(() => import('./components/EngagementCharts').then((module) => ({ default: module.EngagementCharts })));

export const App: React.FC = () => {
  const authResultRequested = useRef(false);
  const [report, setReport] = useState<AccountReport | null>(null);
  const [loading, setLoading] = useState(() => window.location.hash === '#auth_success');
  const [globalError, setGlobalError] = useState<string | null>(() => {
    const error = new URLSearchParams(window.location.search).get('error');
    return error ? `Error en autenticación de Instagram: ${error}` : null;
  });

  const [isLegalOpen, setIsLegalOpen] = useState(false);
  const [activeLegalDoc, setActiveLegalDoc] = useState<LegalDocType>('aviso-legal');

  const handleOpenLegal = (doc: LegalDocType) => {
    setActiveLegalDoc(doc);
    setIsLegalOpen(true);
  };

  const handleSetReport = (newReport: AccountReport) => {
    setReport(newReport);
    setGlobalError(null);
    window.scrollTo({ top: 0, behavior: 'smooth' });

    if (newReport.overall_grade === 'A') {
      void import('canvas-confetti').then(({ default: confetti }) => confetti({
        particleCount: 80,
        spread: 70,
        origin: { y: 0.6 },
        colors: ['#10B981', '#34D399', '#6EE7B7', '#FBBF24'],
      }));
    }
  };

  useEffect(() => {
    const searchParams = new URLSearchParams(window.location.search);
    const error = searchParams.get('error');
    const authSuccess = window.location.hash === '#auth_success';

    if (error) {
      window.history.replaceState({}, document.title, window.location.pathname);
    } else if (authSuccess) {
      if (authResultRequested.current) return;
      authResultRequested.current = true;
      setLoading(true);
      window.history.replaceState({}, document.title, window.location.pathname);
      getAuthResult()
        .then((res) => {
          handleSetReport(res);
        })
        .catch((err) => {
          setGlobalError(err.message || 'Error al analizar la cuenta autenticada');
        })
        .finally(() => setLoading(false));
    }
  }, []);

  const handleLoadDemo = async (tier: 'A' | 'B' | 'D' | 'F') => {
    setLoading(true);
    setGlobalError(null);
    try {
      const demoReport = await analyzeDemo(tier);
      handleSetReport(demoReport);
    } catch (err: any) {
      setGlobalError(err.message || 'Error al cargar la cuenta demo');
    } finally {
      setLoading(false);
    }
  };

  const handleReset = () => {
    setReport(null);
    setGlobalError(null);
  };

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-sans selection:bg-pink-500/30 selection:text-pink-200">
      <Navbar
        onReset={handleReset}
        report={report}
        onLoadDemo={handleLoadDemo}
      />

      <main className="flex-1 max-w-7xl w-full mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {globalError && (
          <div className="mb-8 p-4 rounded-2xl bg-rose-500/10 border border-rose-500/30 text-rose-400 text-sm font-medium flex items-center justify-between">
            <span>⚠️ {globalError}</span>
            <button
              onClick={() => setGlobalError(null)}
              className="text-xs text-slate-400 hover:text-white"
            >
              Cerrar
            </button>
          </div>
        )}

        {loading ? (
          <div className="flex flex-col items-center justify-center py-28 space-y-4">
            <div className="w-12 h-12 border-4 border-pink-500 border-t-transparent rounded-full animate-spin" />
            <p className="text-sm font-semibold text-slate-300">
              Extrayendo métricas de Instagram Graph API y calculando calificación...
            </p>
          </div>
        ) : report ? (
          <div className="space-y-8 animate-in fade-in duration-500">
            <ExecutiveSummaryCard report={report} />
            <KPIStatsGrid report={report} />
            <Suspense fallback={<div className="h-80 rounded-3xl bg-slate-900/80 animate-pulse" />}>
              <EngagementCharts report={report} />
            </Suspense>
            <ReportActionPlan report={report} />
            <MediaGrid media={report.media_analysis} />

            <div className="flex justify-center pt-6 no-print">
              <button
                onClick={() => window.scrollTo({ top: 0, behavior: 'smooth' })}
                className="inline-flex items-center gap-2 px-5 py-2.5 rounded-2xl bg-slate-900 hover:bg-slate-800 border border-slate-800 text-xs font-semibold text-slate-300 transition-colors"
              >
                <ArrowUp className="w-4 h-4 text-pink-400" />
                Volver Arriba
              </button>
            </div>
          </div>
        ) : (
          <HeroConnect
            onReportLoaded={handleSetReport}
          />
        )}
      </main>

      <PreFooter onOpenLegal={handleOpenLegal} />

      <footer className="border-t border-slate-900/80 bg-slate-950 py-6 text-xs text-slate-500 no-print">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 flex flex-col sm:flex-row items-center justify-between gap-4">
          <p>
            © {new Date().getFullYear()} ReviewMySocialNetworks. Todos los derechos reservados.
          </p>

          <div className="flex flex-wrap items-center gap-4 text-xs">
            <button
              onClick={() => handleOpenLegal('aviso-legal')}
              className="text-slate-400 hover:text-slate-200 transition-colors"
            >
              Aviso Legal
            </button>
            <span className="text-slate-700">•</span>
            <button
              onClick={() => handleOpenLegal('privacidad')}
              className="text-slate-400 hover:text-slate-200 transition-colors"
            >
              Privacidad
            </button>
            <span className="text-slate-700">•</span>
            <button
              onClick={() => handleOpenLegal('cookies')}
              className="text-slate-400 hover:text-slate-200 transition-colors"
            >
              Cookies
            </button>
            <span className="text-slate-700">•</span>
            <button
              onClick={() => handleOpenLegal('terminos')}
              className="text-slate-400 hover:text-slate-200 transition-colors"
            >
              Términos
            </button>
            <span className="text-slate-700">•</span>
            <button
              onClick={() => handleOpenLegal('eliminacion-datos')}
              className="text-slate-400 hover:text-slate-200 transition-colors"
            >
              Eliminar Datos
            </button>
          </div>
        </div>
      </footer>

      <LegalModal
        key={activeLegalDoc}
        isOpen={isLegalOpen}
        initialDoc={activeLegalDoc}
        onClose={() => setIsLegalOpen(false)}
      />
    </div>
  );
};

export default App;
