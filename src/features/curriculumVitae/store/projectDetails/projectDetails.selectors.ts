import { createLazySliceStateSelector } from '@shared/utils/createLazySliceStateSelector'
import { getProjectDetailsInitialState } from './projectDetails.state'

const sliceStateSelector = createLazySliceStateSelector(
  'projectDetails',
  getProjectDetailsInitialState(),
)

export const selectBudgetSettings = sliceStateSelector((state) => state.project)

export const selectIsBudgetSettingsLoading = sliceStateSelector((state) => state.isLoading)
