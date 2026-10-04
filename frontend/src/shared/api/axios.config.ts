import { sessionExpired } from '@features/auth/store/session'
import { installAnalyticsMocks } from '@features/curriculumVitae/mocks/analyticsMocks'
import { dispatch } from '@shared/store'
import axios, { AxiosError } from 'axios'

// В браузере запрос по умолчанию идёт на тот же origin: в проде это Caddy, который проксирует
// /api/v1 на бэкенд внутри docker-сети — без CORS и без хардкода адреса.
const baseURL = import.meta.env.VITE_API_BASE_URL ?? '/api/v1'

export const api = axios.create({
  baseURL,
  timeout: 3 * 60 * 1000,
  withCredentials: true,
  headers: {
    Accept: 'application/json',
    'Content-Type': 'application/json',
  },
})

// Локальная разработка без Telegram: бэк с AUTH_MODE=mock пускает по этому заголовку вместо
// сессионной cookie. Только для dev-сборки — в проде переменной нет, и ветка не компилируется в бандл.
if (import.meta.env.DEV && import.meta.env.VITE_MOCK_USER_ID) {
  api.defaults.headers.common['X-Mock-User-ID'] = import.meta.env.VITE_MOCK_USER_ID
}

// Моки ручек аналитики, пока бэкенд их не сделал (см. docs/analytics). Только dev.
if (import.meta.env.DEV && import.meta.env.VITE_MOCK_ANALYTICS === 'true') {
  installAnalyticsMocks(api)
}

// Сессия — httpOnly cookie, обновить её самим нельзя (нет ручки /auth/refresh — только полный
// повторный вход через Telegram). Если владельческая ручка ответила 401, cookie протухла или её
// нет: сразу помечаем сессию как «не залогинен», чтобы экран показал форму входа.
api.interceptors.response.use(
  (res) => res,
  (error: AxiosError) => {
    const url = error.config?.url ?? ''
    const isAuthEndpoint = url.startsWith('/auth/')

    if (error.response?.status === 401 && !isAuthEndpoint) {
      dispatch(sessionExpired())
    }

    return Promise.reject(error)
  },
)
