import { createLazySliceStateSelector } from '@shared/utils/createLazySliceStateSelector'
import { getCurriculumVitaeInitialState } from './curriculumVitae.state'

const sliceStateSelector = createLazySliceStateSelector(
  'curriculumVitae',
  getCurriculumVitaeInitialState(),
)

export const selectPersonData = sliceStateSelector((state) => state.person)

export const selectPersonDataError = sliceStateSelector((state) => state.error)

export const selectIsPersonDataLoading = sliceStateSelector((state) => state.isLoading)
