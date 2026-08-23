# 📊 ReviewMySocialNetworks

Plataforma web integral para **auditoría, analítica y diagnóstico de cuentas de Instagram** utilizando la **Instagram Graph API de Meta**. Desarrollada con un backend en **Go puro (utilizando exclusivamente la librería estándar `net/http`)** y un frontend moderno y reactivo en **TypeScript, React, Vite, Tailwind CSS y Chart.js**.

---

## 🌟 Características Principales

- **Auditoría Integral y Algoritmo de Calificación A/B/D/F**:
  - 🟢 **Grado A (85 - 100 pts)**: Cuenta Excelente / Óptima. Alta tasa de interacción (>3.5-5%), cadencia constante de publicación, diversificación ejemplar (Reels + Carruseles).
  - 🔵 **Grado B (70 - 84 pts)**: Cuenta Buena con base sólida. Rendimiento positivo con margen de optimización en formatos o regularidad.
  - 🟠 **Grado D (50 - 69 pts)**: Cuenta Decente que requiere optimizaciones importantes. Baja interacción relativa o publicación irregular.
  - 🔴 **Grado F (0 - 49 pts)**: Cuenta Crítica o de nivel muy bajo. Inactividad prolongada o indicios de seguidores fantasma.
- **Métricas Evaluadas**:
  - **Seguidores y Autoridad**: Seguidores, Seguidos, Ratio Seguidores/Seguidos.
  - **Interacción (Engagement)**: Tasa de interacción por publicación y media global, likes promedio, comentarios promedio, guardados (saves) y comparación contra benchmarks del sector (2.0%).
  - **Cadencia y Frecuencia**: Días promedio entre publicaciones, ritmo semanal/mensual, días sin publicar y detección del **Mejor Día** y **Mejor Hora** para publicar según el histórico.
  - **Diversidad de Formatos**: Desglose y rendimiento comparativo de Reels vs Carruseles vs Fotos individuales.
- **Informe Ejecutivo Cualitativo**:
  - Diagnóstico ejecutivo en español.
  - Lista de **Puntos Fuertes Detectados**.
  - Lista de **Áreas de Mejora Críticas**.
  - **Plan de Acción Priorizado** (Prioridad Alta, Media, Baja) con estimación de impacto en alcance y engagement.
- **Gráficos Interactivos**:
  - Línea de tiempo de Engagement y Likes por publicación.
  - Distribución de formatos con métricas de rendimiento.
  - Distribución de actividad por día de la semana.
- **Múltiples Modos de Conexión**:
  1. **Inicio de Sesión Oficial OAuth 2.0** con Instagram / Meta Graph API.
  2. **Análisis con Access Token Directo** (útil para tokens de *Graph API Explorer* o tokens de larga duración).
  3. **Modo Demo Interactivo**: 4 perfiles reales precargados (Grados A, B, D y F) para pruebas inmediatas sin necesidad de tokens.
- **Exportación e Impresión**:
  - Vista limpia optimizada para impresión o guardado en PDF del informe de auditoría.

---

## 🛠️ Stack Tecnológico

- **Backend**: Go (Go 1.22+) con **100% librería estándar `net/http`** (sin frameworks externos como Gin o Echo).
- **Frontend**: TypeScript, React 19, Vite, Tailwind CSS v4, Lucide Icons, Chart.js, react-chartjs-2, canvas-confetti.
- **API**: Meta / Instagram Graph API v20.0 (`https://graph.facebook.com/v20.0` y `https://graph.instagram.com`).

---

## 🚀 Inicio Rápido

### Requisitos Previos
- **Go** (versión 1.22 o superior)
- **Node.js** y **pnpm**

### 1. Clonar e Instalar Dependencias
```bash
# Instalar dependencias del frontend con pnpm
pnpm --dir web install
```

### 2. Configurar Variables de Entorno (Opcional)
Copia `.env.example` a `.env` y rellena tus credenciales de Meta for Developers:
```env
PORT=8080
INSTAGRAM_APP_ID=tu_app_id_aqui
INSTAGRAM_APP_SECRET=tu_app_secret_aqui
INSTAGRAM_REDIRECT_URI=http://localhost:8080/api/auth/callback
FRONTEND_URL=http://localhost:8080
```
> *Nota*: También puedes configurar o cambiar el `App ID` y `App Secret` directamente desde el modal de ajustes en la interfaz web sin necesidad de editar archivos.

### 3. Compilar Frontend
```bash
pnpm --dir web build
```

### 4. Iniciar el Servidor Go
```bash
go run ./cmd/server
```

Abre tu navegador en **`http://localhost:8080`**.

---

## 💻 Desarrollo Local con Hot-Reload

Si deseas desarrollar con recarga en vivo:
1. Inicia el backend Go en una terminal:
   ```bash
   go run ./cmd/server
   ```
2. Inicia el servidor de desarrollo de Vite en otra terminal:
   ```bash
   pnpm --dir web dev
   ```
   (El servidor Vite en `http://localhost:5173` redirige automáticamente todas las peticiones `/api/*` al backend Go en `http://localhost:8080`).

---

## 🧪 Ejecución de Tests Automatizados

```bash
# Ejecutar toda la suite de pruebas unitarias en Go (Scoring, Analyzers, Handlers)
go test -v ./...
```

---

## 🔐 Configuración de la App en Meta for Developers

1. Accede a [Meta for Developers](https://developers.facebook.com/apps).
2. Crea una aplicación y selecciona el caso de uso **Instagram Graph API** / **Inicio de sesión de Facebook para empresas**.
3. En la sección *Configuración básica*, copia tu **Identificador de la app (App ID)** y **Clave secreta (App Secret)**.
4. En *Inicio de sesión con Facebook > Configuración*, añade en **URI de redireccionamiento de OAuth válidos**:
   `http://localhost:8080/api/auth/callback`
5. Introduce el App ID y App Secret en el panel de ReviewMySocialNetworks para iniciar sesión con 1 clic.
