import { SliceCaseReducers } from '@shared/types'
import { CurriculumVitaeError } from './curriculumVitae.slice'
import { LkProfile, LkResumeSummary } from './lk'

// Профиль и список резюме — два независимых запроса (getLkProfile/getLkResumes), поэтому у
// каждого свои isLoading/error/requestId, по тому же принципу отмены, что и у резюме.
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
