import { CurriculumVitaeError } from '@features/curriculumVitae/types'
import axios from 'axios'

/** Форма тела ошибки, которую отдаёт бэкенд: `{ error: { code, message, requestId } }` */
type BackendErrorBody = {
  error?: {
    requestId?: string
  }
}

/**
 * Превращает ошибку axios в `CurriculumVitaeError`. Если ответ реально пришёл от бэкенда,
 * в нём остаются HTTP-статус и requestId (для саппорта); сетевая ошибка/таймаут — без них,
 * ответа-то не было.
 */
export const parseApiError = (error: unknown): CurriculumVitaeError => {
  if (!axios.isAxiosError(error) || !error.response) {
    return { kind: 'cannotGetPersonData' }
  }

  const { status, data } = error.response
  const requestId = (data as BackendErrorBody | undefined)?.error?.requestId

  return {
    kind: status === 404 ? 'notFound' : 'cannotGetPersonData',
    status,
    requestId,
  }
}
