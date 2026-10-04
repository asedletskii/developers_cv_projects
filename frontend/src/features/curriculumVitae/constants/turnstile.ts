/** Публичный site key Cloudflare Turnstile. Без него защита контактов выключена */
export const TURNSTILE_SITE_KEY = import.meta.env.VITE_TURNSTILE_SITE_KEY

export const TURNSTILE_SCRIPT_SRC =
  'https://challenges.cloudflare.com/turnstile/v0/api.js?render=explicit'

/** Сколько ждём токен, если задача так и не решена */
export const TURNSTILE_TOKEN_TIMEOUT_MS = 60_000
