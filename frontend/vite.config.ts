import react from '@vitejs/plugin-react'
import { defineConfig, loadEnv, Plugin } from 'vite'
import tsconfigPaths from 'vite-tsconfig-paths'

const FALSE_VALUES = new Set(['false', '0', 'no'])

/** В Docker `/config.js` пишет `docker-entrypoint.sh` из переменных окружения контейнера. В dev контейнера
 *  нет, поэтому отдаём тот же файл из `.env` / окружения — иначе локально всегда был бы дефолт из `public/` */
const runtimeConfigDev = (env: Record<string, string>): Plugin => ({
  name: 'runtime-config-dev',
  configureServer(server) {
    server.middlewares.use('/config.js', (_req, res) => {
      const respectDoNotTrack = !FALSE_VALUES.has((env.RESPECT_DO_NOT_TRACK ?? 'true').toLowerCase())

      res.setHeader('Content-Type', 'application/javascript')
      res.setHeader('Cache-Control', 'no-cache')
      res.end(`window.__CONFIG__ = { respectDoNotTrack: ${respectDoNotTrack} }\n`)
    })
  },
})

export default defineConfig(({ mode }) => {
  // Vite не кладёт `.env` в `process.env` конфига, поэтому читаем его сами: переменные работают и из `.env`,
  // и из окружения (`VITE_PROXY_TARGET=… yarn start`)
  const env = loadEnv(mode, process.cwd(), ['VITE_', 'RESPECT_'])

  return {
    plugins: [react(), tsconfigPaths(), runtimeConfigDev(env)],
    server: {
      proxy: {
        '/api': {
          target: env.VITE_PROXY_TARGET || process.env.VITE_PROXY_TARGET || 'http://localhost:8080',
          changeOrigin: true,
          configure: (proxy) => proxy.on('proxyReq', (req) => req.removeHeader('origin')),
        },
      },
    },
    build: {
      outDir: 'build',
    },
  }
})
