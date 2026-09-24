import { LkSliceState } from '@features/lk/types'

export function getLkInitialState(): LkSliceState {
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
