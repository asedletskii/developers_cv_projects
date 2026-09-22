import { SliceCaseReducers } from '@shared/types'
import { Person } from './curriculumVitae'

export type CurriculumVitaeErrorKind = 'notFound' | 'cannotGetPersonData'

export type CurriculumVitaeError = {
  kind: CurriculumVitaeErrorKind

  /** HTTP-статус ответа бэкенда — есть только если ошибка реально пришла оттуда, а не из сети/таймаута */
  status?: number

  /** requestId из тела ошибки бэкенда — по нему ищут проблему в его логах */
  requestId?: string
}

export type CurriculumVitaeSliceState = {
  person: Person
  isLoading: boolean
  error: CurriculumVitaeError | null

  /** requestId последнего запущенного getPersonData — чтобы игнорировать ответ уже отменённого запроса */
  currentRequestId: string | null
}

export type CurriculumVitaeSliceReducers = SliceCaseReducers<CurriculumVitaeSliceState> & {
  clearCurriculumVitaeState(state: CurriculumVitaeSliceState): void
}
