import { LkError } from '@features/lk/types'
import axios from 'axios'

/** Форма тела ошибки, которую отдаёт бэкенд: `{ error: { code, message, requestId } }` */
type BackendErrorBody = {
  error?: {
    requestId?: string
  }
}

/**
 * Превращает ошибку axios в `LkError`. Если ответ реально пришёл от бэкенда, в нём остаются
 * HTTP-статус и requestId (для саппорта); сетевая ошибка/таймаут — без них, ответа-то не было.
 */
export const parseLkError = (error: unknown): LkError => {
  if (!axios.isAxiosError(error) || !error.response) {
    return { kind: 'cannotGetLk' }
  }

  const { status, data } = error.response
  const requestId = (data as BackendErrorBody | undefined)?.error?.requestId

  return {
    kind: status === 404 ? 'notFound' : 'cannotGetLk',
    status,
    requestId,
  }
}
