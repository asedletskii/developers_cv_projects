import { SliceCaseReducers } from '@shared/types'
import { Resume } from './resume'

export type CurriculumVitaeErrorKind = 'notFound' | 'cannotGetPersonData'

export type CurriculumVitaeError = {
  kind: CurriculumVitaeErrorKind
  status?: number
  requestId?: string
}

export type CurriculumVitaeSliceState = {
  resume: Resume
  isLoading: boolean
  error: CurriculumVitaeError | null
  currentRequestId: string | null
}

export type CurriculumVitaeSliceReducers = SliceCaseReducers<CurriculumVitaeSliceState> & {
  clearCurriculumVitaeState(state: CurriculumVitaeSliceState): void
}
