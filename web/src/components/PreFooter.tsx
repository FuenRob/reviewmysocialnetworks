import React from 'react';
import { Shield, Lock, FileText, Cookie, Trash2, CheckCircle2 } from 'lucide-react';
import { InstagramIcon } from './InstagramIcon';
import type { LegalDocType } from './LegalModal';

interface Props {
  onOpenLegal: (doc: LegalDocType) => void;
}

export const PreFooter: React.FC<Props> = ({ onOpenLegal }) => {
  return (
    <section className="mt-20 border-t border-slate-900 bg-slate-950/90 pt-12 pb-8 no-print">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 md:grid-cols-12 gap-8 lg:gap-12 pb-10 border-b border-slate-900">
          <div className="md:col-span-5 space-y-4">
            <div className="flex items-center gap-3">
              <div className="w-9 h-9 rounded-xl bg-gradient-to-tr from-yellow-500 via-pink-500 to-purple-600 p-0.5 shadow-md shadow-pink-500/10">
                <div className="w-full h-full bg-slate-950 rounded-[10px] flex items-center justify-center">
                  <InstagramIcon className="w-4 h-4 text-white" />
                </div>
              </div>
              <span className="text-base font-black tracking-tight text-white">
                ReviewMy<span className="bg-gradient-to-r from-pink-500 to-purple-400 bg-clip-text text-transparent">SocialNetworks</span>
              </span>
            </div>

            <p className="text-xs text-slate-400 leading-relaxed max-w-md">
              Plataforma de auditoría inteligente y diagnóstico de cuentas de Instagram. Analiza el rendimiento orgánico, engagement real y ritmo de publicación a través de la API oficial Graph de Meta.
            </p>

            <div className="flex flex-wrap items-center gap-3 pt-1">
              <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 text-[11px] font-semibold">
                <CheckCircle2 className="w-3.5 h-3.5" /> Conexión Oficial Meta OAuth
              </span>
              <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-xl bg-blue-500/10 border border-blue-500/20 text-blue-400 text-[11px] font-semibold">
                <Lock className="w-3.5 h-3.5" /> Sesión OAuth protegida
              </span>
            </div>
          </div>

          <div className="md:col-span-4 space-y-3">
            <h4 className="text-xs font-bold uppercase tracking-wider text-slate-300 flex items-center gap-2">
              <Shield className="w-4 h-4 text-pink-400" />
              Marco Legal y Privacidad
            </h4>
            <ul className="space-y-2 text-xs">
              <li>
                <button
                  onClick={() => onOpenLegal('aviso-legal')}
                  className="text-slate-400 hover:text-pink-400 transition-colors flex items-center gap-2 text-left"
                >
                  <FileText className="w-3.5 h-3.5 text-slate-500" />
                  Aviso Legal y LSSI-CE
                </button>
              </li>
              <li>
                <button
                  onClick={() => onOpenLegal('privacidad')}
                  className="text-slate-400 hover:text-pink-400 transition-colors flex items-center gap-2 text-left"
                >
                  <Lock className="w-3.5 h-3.5 text-slate-500" />
                  Política de Privacidad (RGPD / GDPR)
                </button>
              </li>
              <li>
                <button
                  onClick={() => onOpenLegal('cookies')}
                  className="text-slate-400 hover:text-pink-400 transition-colors flex items-center gap-2 text-left"
                >
                  <Cookie className="w-3.5 h-3.5 text-slate-500" />
                  Política de Cookies
                </button>
              </li>
              <li>
                <button
                  onClick={() => onOpenLegal('terminos')}
                  className="text-slate-400 hover:text-pink-400 transition-colors flex items-center gap-2 text-left"
                >
                  <Shield className="w-3.5 h-3.5 text-slate-500" />
                  Términos y Condiciones de Uso
                </button>
              </li>
              <li>
                <button
                  onClick={() => onOpenLegal('eliminacion-datos')}
                  className="text-slate-400 hover:text-pink-400 transition-colors flex items-center gap-2 text-left"
                >
                  <Trash2 className="w-3.5 h-3.5 text-slate-500" />
                  Eliminación de Datos (Meta Platform)
                </button>
              </li>
            </ul>
          </div>

          <div className="md:col-span-3 space-y-3">
            <h4 className="text-xs font-bold uppercase tracking-wider text-slate-300">
              Garantía de Privacidad
            </h4>
            <div className="bg-slate-900/60 border border-slate-800/80 rounded-2xl p-3.5 space-y-2 text-[11px] text-slate-400 leading-relaxed">
              <p>
                <strong>Almacenamiento Cero:</strong> No guardamos tus contraseñas ni tus publicaciones en bases de datos. Los datos se procesan en memoria exclusivamente para calcular tu informe.
              </p>
              <p className="text-[10px] text-slate-500 border-t border-slate-800/80 pt-2">
                Instagram® es una marca comercial de Meta Platforms, Inc. Este sitio web no está afiliado ni patrocinado por Meta.
              </p>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
};
