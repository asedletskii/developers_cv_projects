import { AUTH_API_ROUTES } from '@features/auth/constants'
import { Session } from '@features/auth/types'
import { api } from '@shared/api'

class AuthApi {
  /** Полная навигация браузера (не XHR) — бэк ставит cookie и редиректит на Telegram. */
  getTelegramLoginUrl(): string {
    return `${api.defaults.baseURL}${AUTH_API_ROUTES.telegramLogin}`
  }

  async getSession(signal?: AbortSignal): Promise<Session> {
    const response = await api.get<Session>(AUTH_API_ROUTES.me, { signal })
    return response.data
  }

  async logout(): Promise<void> {
    await api.post(AUTH_API_ROUTES.logout)
  }
}

export const authApi = new AuthApi()
