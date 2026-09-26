import React from 'react';
import { ArrowLeft, Lock, ShieldCheck } from 'lucide-react';

export const PrivacyPolicyPage: React.FC = () => {
  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 font-sans">
      <header className="border-b border-slate-800 bg-slate-950/95">
        <div className="max-w-4xl mx-auto px-6 py-5 flex items-center justify-between gap-4">
          <a href="/" className="text-sm font-black tracking-tight text-white hover:text-pink-400 transition-colors">
            ReviewMy<span className="text-pink-400">SocialNetworks</span>
          </a>
          <a href="/" className="inline-flex items-center gap-2 text-xs font-semibold text-slate-400 hover:text-white transition-colors">
            <ArrowLeft className="w-4 h-4" /> Volver a la aplicación
          </a>
        </div>
      </header>

      <main className="max-w-4xl mx-auto px-6 py-12 sm:py-16">
        <article className="rounded-3xl border border-slate-800 bg-slate-900/80 p-6 sm:p-10 shadow-2xl shadow-black/20">
          <div className="flex items-start gap-4 mb-10">
            <div className="rounded-2xl bg-pink-500/10 border border-pink-500/20 p-3 text-pink-400">
              <Lock className="w-6 h-6" />
            </div>
            <div>
              <h1 className="text-3xl sm:text-4xl font-black text-white">Política de Privacidad</h1>
              <p className="mt-2 text-sm text-slate-400">ReviewMySocialNetworks · Última actualización: 26 de septiembre de 2026</p>
            </div>
          </div>

          <div className="space-y-8 text-sm leading-7 text-slate-300">
            <section>
              <h2 className="text-lg font-bold text-white mb-2">1. Responsable y alcance</h2>
              <p>
                ReviewMySocialNetworks es una aplicación independiente para auditar y analizar cuentas de Instagram y TikTok. Esta política explica qué datos se tratan cuando autorizas el acceso mediante OAuth y utilizas el servicio.
              </p>
            </section>

            <section>
              <h2 className="text-lg font-bold text-white mb-2">2. Datos que tratamos</h2>
              <p>
                Solicitamos únicamente los datos necesarios para generar el informe: identificador y datos básicos del perfil, seguidores y seguidos, publicaciones o vídeos recientes y sus métricas disponibles, como likes, comentarios, compartidos, visualizaciones y alcance.
              </p>
              <p className="mt-3">
                Nunca solicitamos ni conocemos tu contraseña. La autenticación se realiza exclusivamente en los servidores oficiales de Meta/Instagram o TikTok.
              </p>
            </section>

            <section>
              <h2 className="text-lg font-bold text-white mb-2">3. Finalidad y base jurídica</h2>
              <p>
                Tratamos los datos para autenticar tu cuenta, calcular métricas y mostrarte un diagnóstico de rendimiento y recomendaciones. La base jurídica es tu consentimiento explícito al aceptar los permisos en la pantalla de autorización de Instagram o TikTok, conforme al artículo 6.1.a del RGPD.
              </p>
            </section>

            <section>
              <h2 className="text-lg font-bold text-white mb-2">4. Conservación y seguridad</h2>
              <p>
                Aplicamos una política de almacenamiento cero permanente: los tokens y las métricas se procesan en memoria durante la sesión y no se guardan en bases de datos permanentes. Usamos conexiones cifradas y controles técnicos para proteger las comunicaciones.
              </p>
            </section>

            <section>
              <h2 className="text-lg font-bold text-white mb-2">5. Cesiones y terceros</h2>
              <p>
                No vendemos tus datos ni los usamos para publicidad. El servicio se comunica con las API oficiales de Meta/Instagram y TikTok únicamente para obtener los datos que hayas autorizado. Instagram y Meta son marcas de Meta Platforms, Inc.; TikTok es una marca de ByteDance Ltd. ReviewMySocialNetworks no está patrocinada ni administrada por dichas compañías.
              </p>
            </section>

            <section>
              <h2 className="text-lg font-bold text-white mb-2">6. Tus derechos y revocación</h2>
              <p>
                Puedes retirar el consentimiento y revocar el acceso en cualquier momento desde la configuración de aplicaciones conectadas de Instagram o TikTok. También puedes solicitar acceso, rectificación, supresión, limitación u oposición al tratamiento contactando con el responsable del servicio.
              </p>
            </section>

            <section>
              <h2 className="text-lg font-bold text-white mb-2">7. Contacto</h2>
              <p>
                Para consultas sobre privacidad o eliminación de datos, utiliza el canal de contacto indicado por el titular de la aplicación. La revocación del acceso en Instagram o TikTok impide futuras consultas de datos.
              </p>
            </section>
          </div>

          <div className="mt-10 rounded-2xl border border-emerald-500/20 bg-emerald-500/10 p-4 flex gap-3 text-sm text-emerald-200">
            <ShieldCheck className="w-5 h-5 shrink-0 text-emerald-400 mt-0.5" />
            <p>Esta página es la URL pública de política de privacidad de ReviewMySocialNetworks para las aplicaciones de Instagram y TikTok.</p>
          </div>
        </article>
      </main>
    </div>
  );
};

