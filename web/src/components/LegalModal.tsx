import React, { useState } from 'react';
import { X, Shield, Lock, FileText, Cookie, Trash2, CheckCircle } from 'lucide-react';

export type LegalDocType = 'aviso-legal' | 'privacidad' | 'cookies' | 'terminos' | 'eliminacion-datos';

interface Props {
  isOpen: boolean;
  initialDoc?: LegalDocType;
  onClose: () => void;
}

export const LegalModal: React.FC<Props> = ({ isOpen, initialDoc = 'aviso-legal', onClose }) => {
  const [activeTab, setActiveTab] = useState<LegalDocType>(initialDoc);

  React.useEffect(() => {
    setActiveTab(initialDoc);
  }, [initialDoc]);

  if (!isOpen) return null;

  const tabs: { id: LegalDocType; label: string; icon: React.ComponentType<{ className?: string }> }[] = [
    { id: 'aviso-legal', label: 'Aviso Legal', icon: FileText },
    { id: 'privacidad', label: 'Política de Privacidad', icon: Lock },
    { id: 'cookies', label: 'Política de Cookies', icon: Cookie },
    { id: 'terminos', label: 'Términos del Servicio', icon: Shield },
    { id: 'eliminacion-datos', label: 'Eliminación de Datos', icon: Trash2 },
  ];

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/85 backdrop-blur-md animate-in fade-in duration-200">
      <div className="bg-slate-900 border border-slate-800 rounded-3xl w-full max-w-4xl shadow-2xl overflow-hidden flex flex-col max-h-[90vh]">
        <div className="p-5 sm:p-6 border-b border-slate-800 flex items-center justify-between bg-slate-950/40">
          <div className="flex items-center gap-3">
            <div className="p-2.5 rounded-2xl bg-indigo-500/10 border border-indigo-500/20 text-indigo-400">
              <Shield className="w-5 h-5" />
            </div>
            <div>
              <h3 className="text-lg font-bold text-white">
                Información Legal y Cumplimiento Normativo
              </h3>
              <p className="text-xs text-slate-400">
                Conforme al RGPD (UE 2016/679), LOPDGDD y Políticas de la Plataforma de Meta
              </p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="p-2 rounded-xl text-slate-400 hover:text-white hover:bg-slate-800 transition-colors"
            aria-label="Cerrar modal"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        <div className="flex border-b border-slate-800 overflow-x-auto bg-slate-950/60 scrollbar-none px-4 pt-2 gap-1.5">
          {tabs.map((t) => {
            const Icon = t.icon;
            const isActive = activeTab === t.id;
            return (
              <button
                key={t.id}
                onClick={() => setActiveTab(t.id)}
                className={`flex items-center gap-2 px-4 py-3 text-xs font-semibold rounded-t-2xl whitespace-nowrap transition-all border-t border-x ${
                  isActive
                    ? 'bg-slate-900 text-pink-400 border-slate-800 border-b-slate-900 shadow-sm'
                    : 'text-slate-400 border-transparent hover:text-slate-200 hover:bg-slate-900/40'
                }`}
              >
                <Icon className={`w-4 h-4 ${isActive ? 'text-pink-400' : 'text-slate-500'}`} />
                <span>{t.label}</span>
              </button>
            );
          })}
        </div>

        <div className="p-6 sm:p-8 overflow-y-auto space-y-6 text-slate-300 text-sm leading-relaxed max-h-[calc(90vh-140px)]">
          {activeTab === 'aviso-legal' && (
            <div className="space-y-5 animate-in fade-in duration-300">
              <div>
                <h4 className="text-xl font-bold text-white mb-1">Aviso Legal</h4>
                <p className="text-xs text-slate-500">Última actualización: Agosto de 2026</p>
              </div>

              <section className="space-y-2">
                <h5 className="text-sm font-bold text-slate-100">1. Datos Identificativos del Responsable</h5>
                <p className="text-xs text-slate-400 leading-relaxed">
                  En cumplimiento con el artículo 10 de la Ley 34/2002, de 11 de julio, de Servicios de la Sociedad de la Información y de Comercio Electrónico (LSSI-CE), se informa que la plataforma <strong>ReviewMySocialNetworks</strong> es una herramienta digital de auditoría, diagnóstico y analítica web orientada a cuentas profesionales y de creadores de redes sociales.
                </p>
              </section>

              <section className="space-y-2">
                <h5 className="text-sm font-bold text-slate-100">2. Finalidad de la Plataforma</h5>
                <p className="text-xs text-slate-400 leading-relaxed">
                  La plataforma proporciona métricas cuantitativas, análisis de interacción (engagement rate), cadencia de publicación y recomendaciones de optimización de contenido a los usuarios titulares de cuentas de Instagram que autoricen el acceso mediante el protocolo estándar de OAuth de Meta / Instagram Graph API.
                </p>
              </section>

              <section className="space-y-2">
                <h5 className="text-sm font-bold text-slate-100">3. Propiedad Intelectual y Marcas Registradas</h5>
                <p className="text-xs text-slate-400 leading-relaxed">
                  Todos los derechos de propiedad industrial e intelectual del software, código fuente, algoritmos de cálculo, diseño gráfico e interfaces pertenecen a ReviewMySocialNetworks.
                </p>
                <div className="p-3.5 rounded-2xl bg-slate-950/70 border border-slate-800 text-xs text-slate-400 flex items-start gap-2">
                  <span className="text-pink-400 shrink-0 font-bold">ℹ️</span>
                  <span>
                    <strong>Instagram®</strong> y <strong>Meta®</strong> son marcas registradas de Meta Platforms, Inc. ReviewMySocialNetworks es una aplicación independiente y no está patrocinada, respaldada ni administrada directamente por Meta Platforms, Inc.
                  </span>
                </div>
              </section>

              <section className="space-y-2">
                <h5 className="text-sm font-bold text-slate-100">4. Limitación de Responsabilidad</h5>
                <p className="text-xs text-slate-400 leading-relaxed">
                  Las valoraciones, notas (A, B, D, F) y planes de acción ofrecidos por el algoritmo de la plataforma tienen carácter meramente orientativo y estadístico. No garantizan resultados comerciales, incremento específico de ventas ni posicionamiento algorítmico en Instagram.
                </p>
              </section>
            </div>
          )}

          {activeTab === 'privacidad' && (
            <div className="space-y-5 animate-in fade-in duration-300">
              <div>
                <h4 className="text-xl font-bold text-white mb-1">Política de Privacidad y Tratamiento de Datos (RGPD)</h4>
                <p className="text-xs text-slate-500">Conforme al Reglamento General de Protección de Datos (UE 2016/679)</p>
              </div>

              <section className="space-y-2">
                <h5 className="text-sm font-bold text-slate-100">1. Principio de Minimización y Privacidad por Diseño</h5>
                <p className="text-xs text-slate-400 leading-relaxed">
                  En ReviewMySocialNetworks aplicamos una política estricta de <strong>Cero Almacenamiento Permanente de Credenciales</strong>:
                </p>
                <ul className="list-disc list-inside space-y-1.5 text-xs text-slate-400 pl-2">
                  <li><strong>Contraseñas:</strong> Nunca solicitamos ni tenemos acceso a tu contraseña de Instagram. El login se realiza exclusivamente en los servidores seguros de Instagram (OAuth).</li>
                  <li><strong>Tokens de Acceso:</strong> El token de acceso temporal se utiliza en tiempo real exclusivamente durante la sesión activa para obtener las métricas de tu cuenta y generar tu informe.</li>
                  <li><strong>Datos Consultados:</strong> Únicamente se consultan los datos estrictamente necesarios concedidos por los permisos: nombre de usuario, foto de perfil, número de seguidores, número de seguidos, lista de publicaciones recientes (fotos, carruseles, reels) y sus contadores públicos de likes, comentarios y métricas de alcance.</li>
                </ul>
              </section>

              <section className="space-y-2">
                <h5 className="text-sm font-bold text-slate-100">2. Base Jurídica del Tratamiento</h5>
                <p className="text-xs text-slate-400 leading-relaxed">
                  La base legal para el tratamiento de tus datos es el <strong>consentimiento explícito</strong> del interesado otorgado al iniciar sesión y aceptar los permisos en la pantalla de autorización de Instagram (Artículo 6.1.a del RGPD).
                </p>
              </section>

              <section className="space-y-2">
                <h5 className="text-sm font-bold text-slate-100">3. Cesión y Transferencias Internacionales</h5>
                <p className="text-xs text-slate-400 leading-relaxed">
                  Tus datos nunca son vendidos, alquilados ni transferidos a terceros con fines publicitarios o comerciales. Las comunicaciones con la API de Instagram se realizan mediante conexiones cifradas SSL/TLS de 256 bits hacia la infraestructura oficial de Meta Platforms, Inc.
                </p>
              </section>

              <section className="space-y-2">
                <h5 className="text-sm font-bold text-slate-100">4. Tus Derechos (ARCO / RGPD)</h5>
                <p className="text-xs text-slate-400 leading-relaxed">
                  Puedes revocar el acceso concedido a la aplicación en cualquier momento directamente desde la configuración de tu cuenta de Instagram (<em>Configuración y actividad ➔ Permisos de sitios web ➔ Aplicaciones y sitios web</em>) o solicitar la eliminación total conforme se describe en la sección de Eliminación de Datos.
                </p>
              </section>
            </div>
          )}

          {activeTab === 'cookies' && (
            <div className="space-y-5 animate-in fade-in duration-300">
              <div>
                <h4 className="text-xl font-bold text-white mb-1">Política de Cookies</h4>
                <p className="text-xs text-slate-500">Transparencia sobre el uso de almacenamiento local y cookies técnicas</p>
              </div>

              <section className="space-y-2">
                <h5 className="text-sm font-bold text-slate-100">1. ¿Qué son las cookies y tecnologías similares?</h5>
                <p className="text-xs text-slate-400 leading-relaxed">
                  Una cookie es un pequeño archivo que se almacena en tu navegador web. Esta aplicación web utiliza almacenamiento de sesión estrictamente técnico y funcional para permitir el funcionamiento de la Single Page Application (SPA).
                </p>
              </section>

              <section className="space-y-2">
                <h5 className="text-sm font-bold text-slate-100">2. Cookies y Almacenamiento Utilizados</h5>
                <div className="overflow-x-auto">
                  <table className="w-full text-xs text-left border border-slate-800 rounded-xl overflow-hidden">
                    <thead className="bg-slate-950 text-slate-300 font-bold border-b border-slate-800">
                      <tr>
                        <th className="p-2.5">Tipo</th>
                        <th className="p-2.5">Nombre / Finalidad</th>
                        <th className="p-2.5">Duración</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-800 text-slate-400">
                      <tr>
                        <td className="p-2.5 font-semibold text-white">Técnica / Estado</td>
                        <td className="p-2.5">Mantiene el estado temporal del informe analizado durante la sesión de navegación.</td>
                        <td className="p-2.5">Sesión actual</td>
                      </tr>
                      <tr>
                        <td className="p-2.5 font-semibold text-white">Seguridad OAuth</td>
                        <td className="p-2.5">Parámetro `state` criptográfico para prevenir ataques CSRF en la autenticación.</td>
                        <td className="p-2.5">Transitoria (10 min)</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </section>

              <section className="space-y-2">
                <h5 className="text-sm font-bold text-slate-100">3. Sin Cookies de Rastreo de Terceros</h5>
                <p className="text-xs text-slate-400 leading-relaxed">
                  No empleamos cookies publicitarias, trackers entre sitios (cross-site trackers) ni vendemos perfiles de navegación a empresas de publicidad.
                </p>
              </section>
            </div>
          )}

          {activeTab === 'terminos' && (
            <div className="space-y-5 animate-in fade-in duration-300">
              <div>
                <h4 className="text-xl font-bold text-white mb-1">Términos y Condiciones del Servicio</h4>
                <p className="text-xs text-slate-500">Condiciones de uso de la herramienta de auditoría de redes sociales</p>
              </div>

              <section className="space-y-2">
                <h5 className="text-sm font-bold text-slate-100">1. Aceptación de los Términos</h5>
                <p className="text-xs text-slate-400 leading-relaxed">
                  Al conectar tu cuenta o utilizar las herramientas de análisis de ReviewMySocialNetworks, aceptas cumplir los presentes Términos de Servicio y todas las leyes y regulaciones aplicables.
                </p>
              </section>

              <section className="space-y-2">
                <h5 className="text-sm font-bold text-slate-100">2. Uso Autorizado</h5>
                <p className="text-xs text-slate-400 leading-relaxed">
                  El usuario garantiza que es el titular legítimo o cuenta con la debida autorización del propietario para conectar y auditar la cuenta de Instagram vinculada. Queda prohibido el uso de la plataforma con fines fraudulentos, de scraping no autorizado o que vulneren las políticas de la plataforma de Meta.
                </p>
              </section>

              <section className="space-y-2">
                <h5 className="text-sm font-bold text-slate-100">3. Disponibilidad y Límites de la API de Instagram</h5>
                <p className="text-xs text-slate-400 leading-relaxed">
                  El servicio depende de la disponibilidad, cuotas de llamada (Rate Limits) y políticas técnicas de Instagram Graph API. ReviewMySocialNetworks no se hace responsable de interrupciones del servicio derivadas de cambios, caídas o limitaciones técnicas impuestas por Meta Platforms, Inc.
                </p>
              </section>
            </div>
          )}

          {activeTab === 'eliminacion-datos' && (
            <div className="space-y-5 animate-in fade-in duration-300">
              <div>
                <h4 className="text-xl font-bold text-white mb-1">Instrucciones de Eliminación de Datos de Usuario</h4>
                <p className="text-xs text-slate-500">Conforme a las Políticas de la Plataforma de Meta para Desarrolladores</p>
              </div>

              <div className="p-4 rounded-2xl bg-emerald-500/10 border border-emerald-500/30 text-xs text-emerald-300 flex items-start gap-2.5">
                <CheckCircle className="w-4 h-4 text-emerald-400 shrink-0 mt-0.5" />
                <div>
                  <strong>Política de Almacenamiento Cero:</strong> ReviewMySocialNetworks procesa las métricas en memoria volátil de sesión para generar el informe y no almacena tus datos de Instagram en bases de datos permanentes.
                </div>
              </div>

              <section className="space-y-3">
                <h5 className="text-sm font-bold text-slate-100">¿Cómo revocar el acceso y eliminar cualquier dato asociado?</h5>
                <p className="text-xs text-slate-400 leading-relaxed">
                  Si deseas revocar el acceso concedido a la aplicación y asegurarte de que ningún dato de tu perfil continúe vinculado, sigue estos pasos:
                </p>
                <ol className="list-decimal list-inside space-y-2 text-xs text-slate-300 pl-2 bg-slate-950/60 p-4 rounded-2xl border border-slate-800">
                  <li>Abre la aplicación de <strong>Instagram</strong> en tu móvil o accede a <strong>instagram.com</strong> en tu navegador.</li>
                  <li>Ve a tu <strong>Perfil</strong> ➔ Menú de opciones (tres líneas) ➔ <strong>Configuración y privacidad</strong>.</li>
                  <li>En el apartado de <em>"Permisos de sitios web"</em>, selecciona <strong>Aplicaciones y sitios web</strong>.</li>
                  <li>Localiza <strong>ReviewMySocialNetworks</strong> en la pestaña de <em>"Activas"</em> y haz clic en <strong>Eliminar</strong> (Revocar acceso).</li>
                </ol>
                <p className="text-xs text-slate-400">
                  Al completar este proceso, Instagram invalidará de inmediato cualquier token de acceso y la aplicación no podrá realizar ninguna consulta futura sobre tu cuenta.
                </p>
              </section>
            </div>
          )}
        </div>

        <div className="p-4 border-t border-slate-800 bg-slate-950/80 flex items-center justify-between text-xs text-slate-500">
          <span>ReviewMySocialNetworks • Cumplimiento Normativo UE</span>
          <button
            onClick={onClose}
            className="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-white font-semibold transition-colors"
          >
            Entendido
          </button>
        </div>
      </div>
    </div>
  );
};
