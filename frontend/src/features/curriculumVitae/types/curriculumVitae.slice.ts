import { SliceCaseReducers } from '@shared/types'
import { Person } from './curriculumVitae'

export type CurriculumVitaeError = 'notFound' | 'cannotGetPersonData'

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
