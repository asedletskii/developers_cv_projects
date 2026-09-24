import { SliceCaseReducers } from '@shared/types'
import { LkProfile, LkResumeSummary } from './lk'

export type LkErrorKind = 'notFound' | 'cannotGetLk'

export type LkError = {
  kind: LkErrorKind

  /** HTTP-статус ответа бэкенда — есть только если ошибка реально пришла оттуда, а не из сети/таймаута */
  status?: number

  /** requestId из тела ошибки бэкенда — по нему ищут проблему в его логах */
  requestId?: string
}

// Профиль и список резюме — два независимых запроса (getLkProfile/getLkResumes), поэтому у
// каждого свои isLoading/error/requestId, по тому же принципу отмены, что и остальные фичи.
export type LkSliceState = {
  profile: LkProfile | null
  resumes: LkResumeSummary[] | null

  isProfileLoading: boolean
  isResumesLoading: boolean

  profileError: LkError | null
  resumesError: LkError | null

  currentProfileRequestId: string | null
  currentResumesRequestId: string | null
}

export type LkSliceReducers = SliceCaseReducers<LkSliceState> & {
  clearLkState(state: LkSliceState): void
}
