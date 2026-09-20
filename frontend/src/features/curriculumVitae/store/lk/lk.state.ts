import { LkSliceState } from '@features/curriculumVitae/types'

export const getLkInitialState = (): LkSliceState => {
  return {
    profile: null,
    resumes: null,
    isProfileLoading: true,
    isResumesLoading: true,
    profileError: null,
    resumesError: null,
    currentProfileRequestId: null,
    currentResumesRequestId: null,
  }
}
