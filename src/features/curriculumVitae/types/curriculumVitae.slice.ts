import { SliceCaseReducers } from '@shared/types'
import { Person } from './curriculumVitae'

export type CurriculumVitaeError = 'notFound' | 'cannotGetPersonData'

export type CurriculumVitaeSliceState = {
  person: Person
  isLoading: boolean
  error: CurriculumVitaeError | null
}

export type CurriculumVitaeSliceReducers = SliceCaseReducers<CurriculumVitaeSliceState> & {
  clearCurriculumVitaeState(state: CurriculumVitaeSliceState): void
}
