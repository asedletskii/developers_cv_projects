import { createLazySliceStateSelector } from '@shared/utils/createLazySliceStateSelector'
import { getLkInitialState } from './lk.state'

const sliceStateSelector = createLazySliceStateSelector('lk', getLkInitialState())

export const selectLkProfile = sliceStateSelector((state) => state.profile)
export const selectLkResumes = sliceStateSelector((state) => state.resumes)

export const selectIsLkProfileLoading = sliceStateSelector((state) => state.isProfileLoading)
export const selectIsLkResumesLoading = sliceStateSelector((state) => state.isResumesLoading)

export const selectLkProfileError = sliceStateSelector((state) => state.profileError)
export const selectLkResumesError = sliceStateSelector((state) => state.resumesError)
