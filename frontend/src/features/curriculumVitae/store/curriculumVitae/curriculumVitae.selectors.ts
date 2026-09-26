import { createLazySliceStateSelector } from '@shared/utils/createLazySliceStateSelector'
import { getCurriculumVitaeInitialState } from './curriculumVitae.state'

const sliceStateSelector = createLazySliceStateSelector(
  'curriculumVitae',
  getCurriculumVitaeInitialState(),
)

export const selectResume = sliceStateSelector((state) => state.resume)

export const selectResumeError = sliceStateSelector((state) => state.error)

export const selectIsResumeLoading = sliceStateSelector((state) => state.isLoading)
