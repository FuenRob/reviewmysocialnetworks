import React from 'react';
import { ArrowLeft, FileText, ShieldCheck } from 'lucide-react';

export const TermsOfServicePage: React.FC = () => {
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
            <div className="rounded-2xl bg-indigo-500/10 border border-indigo-500/20 p-3 text-indigo-400">
              <FileText className="w-6 h-6" />
            </div>
            <div>
              <h1 className="text-3xl sm:text-4xl font-black text-white">Términos del Servicio</h1>
              <p className="mt-2 text-sm text-slate-400">ReviewMySocialNetworks · Última actualización: 26 de septiembre de 2026</p>
            </div>
          </div>

          <div className="space-y-8 text-sm leading-7 text-slate-300">
            <section>
              <h2 className="text-lg font-bold text-white mb-2">1. Aceptación</h2>
              <p>
                Al acceder a ReviewMySocialNetworks o utilizar sus herramientas de análisis, aceptas estos términos y las leyes aplicables. Si no estás de acuerdo, no utilices el servicio.
              </p>
            </section>

            <section>
              <h2 className="text-lg font-bold text-white mb-2">2. Descripción del servicio</h2>
              <p>
                ReviewMySocialNetworks proporciona métricas, análisis de interacción, cadencia, alcance y recomendaciones orientativas para cuentas de Instagram y TikTok, utilizando las API oficiales y los permisos autorizados por el usuario.
              </p>
            </section>

            <section>
              <h2 className="text-lg font-bold text-white mb-2">3. Uso autorizado</h2>
              <p>
                Solo puedes conectar cuentas de Instagram o TikTok de las que seas titular o para las que tengas autorización suficiente. No está permitido utilizar el servicio para acceder a cuentas ajenas, realizar scraping no autorizado, interferir con las plataformas o vulnerar sus políticas.
              </p>
            </section>

            <section>
              <h2 className="text-lg font-bold text-white mb-2">4. Autenticación y datos</h2>
              <p>
                La autenticación se realiza mediante OAuth en los servicios oficiales de Meta/Instagram o TikTok. No solicitamos contraseñas. El tratamiento de los datos se rige por nuestra <a href="/privacy-policy" className="text-pink-400 hover:text-pink-300">Política de Privacidad</a>.
              </p>
            </section>

            <section>
              <h2 className="text-lg font-bold text-white mb-2">5. Naturaleza orientativa de los informes</h2>
              <p>
                Las calificaciones, métricas y recomendaciones tienen carácter informativo y estadístico. No garantizan resultados comerciales, crecimiento de audiencia, alcance ni posicionamiento algorítmico en Instagram o TikTok.
              </p>
            </section>

            <section>
              <h2 className="text-lg font-bold text-white mb-2">6. Disponibilidad y limitaciones</h2>
              <p>
                El servicio depende de la disponibilidad, cuotas, cambios técnicos y políticas de las API de Meta/Instagram y TikTok. Podemos modificar, suspender o limitar funciones cuando sea necesario para mantener la seguridad o el funcionamiento del servicio.
              </p>
            </section>

            <section>
              <h2 className="text-lg font-bold text-white mb-2">7. Propiedad intelectual y marcas</h2>
              <p>
                El software, diseño, código, algoritmos y contenidos de ReviewMySocialNetworks pertenecen a sus respectivos titulares. Instagram y Meta son marcas de Meta Platforms, Inc.; TikTok es una marca de ByteDance Ltd. ReviewMySocialNetworks es una aplicación independiente y no está patrocinada ni administrada por dichas compañías.
              </p>
            </section>

            <section>
              <h2 className="text-lg font-bold text-white mb-2">8. Revocación y contacto</h2>
              <p>
                Puedes dejar de utilizar el servicio y revocar el acceso concedido desde la configuración de aplicaciones conectadas de Instagram o TikTok. Para consultas sobre estos términos, utiliza el canal de contacto indicado por el titular de la aplicación.
              </p>
            </section>
          </div>

          <div className="mt-10 rounded-2xl border border-indigo-500/20 bg-indigo-500/10 p-4 flex gap-3 text-sm text-indigo-200">
            <ShieldCheck className="w-5 h-5 shrink-0 text-indigo-400 mt-0.5" />
            <p>Estos términos se aplican al uso de ReviewMySocialNetworks con cuentas de Instagram y TikTok.</p>
          </div>
        </article>
      </main>
    </div>
  );
};

