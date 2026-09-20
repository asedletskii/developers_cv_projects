import { SliceCaseReducers } from '@shared/types'
import { CurriculumVitaeError } from './curriculumVitae.slice'
import { LkProfile, LkResumeSummary } from './lk'

export type LkSliceState = {
  profile: LkProfile | null
  resumes: LkResumeSummary[] | null

  isProfileLoading: boolean
  isResumesLoading: boolean

  profileError: CurriculumVitaeError | null
  resumesError: CurriculumVitaeError | null

  currentProfileRequestId: string | null
  currentResumesRequestId: string | null
}

export type LkSliceReducers = SliceCaseReducers<LkSliceState> & {
  clearLkState(state: LkSliceState): void
}
