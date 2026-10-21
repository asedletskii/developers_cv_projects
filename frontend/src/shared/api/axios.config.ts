import { sessionExpired } from '@features/auth/store/session'
import { dispatch } from '@shared/store'
import axios, { AxiosError } from 'axios'

const baseURL = import.meta.env.VITE_API_BASE_URL ?? '/api/v1'

export const api = axios.create({
  baseURL,
  timeout: 3 * 60 * 1000,
  withCredentials: true,
  headers: {
    Accept: 'application/json',
    'Content-Type': 'application/json',
    'X-Requested-With': 'XMLHttpRequest',
  },
})

if (import.meta.env.DEV && import.meta.env.VITE_MOCK_USER_ID) {
  api.defaults.headers.common['X-Mock-User-ID'] = import.meta.env.VITE_MOCK_USER_ID
}

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
