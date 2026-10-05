export type RuntimeConfig = {
  respectDoNotTrack: boolean
}

declare global {
  interface Window {
    __CONFIG__?: Partial<Record<keyof RuntimeConfig, unknown>>
  }
}

const DEFAULTS: RuntimeConfig = { respectDoNotTrack: true }

/** Настройки, которые меняются без пересборки фронтенда: контейнер при каждом старте пишет `/config.js`
 *  из переменных окружения (см. `docker-entrypoint.sh`). Нет файла или значение неверного вида — берётся
 *  дефолт; поэтому в dev и в тестах конфиг не нужен. */
export const getRuntimeConfig = (
  source: Window['__CONFIG__'] = typeof window === 'undefined' ? undefined : window.__CONFIG__,
): RuntimeConfig => ({
  respectDoNotTrack:
    typeof source?.respectDoNotTrack === 'boolean'
      ? source.respectDoNotTrack
      : DEFAULTS.respectDoNotTrack,
})
