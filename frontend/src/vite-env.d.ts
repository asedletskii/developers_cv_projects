/// <reference types="vite/client" />

interface ImportMetaEnv {
  /** Базовый URL API. По умолчанию относительный `/api/v1` — расчёт на reverse-proxy перед фронтом. */
  readonly VITE_API_BASE_URL?: string

  /** uuid существующего пользователя — для бэка с AUTH_MODE=mock вместо входа через Telegram. */
  readonly VITE_MOCK_USER_ID?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
