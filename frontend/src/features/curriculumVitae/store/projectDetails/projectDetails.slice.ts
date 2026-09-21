import {
  ProjectDetailsSliceReducers,
  ProjectDetailsSliceState,
} from '@features/curriculumVitae/types'
import { createSlice, WithSlice } from '@reduxjs/toolkit'
import { rootReducer } from '@shared/store'
import { getProjectDetailsInitialState } from './projectDetails.state'
import { getProject } from './projectDetails.thunks'

export const ProjectDetailsSlice = createSlice<
  ProjectDetailsSliceState,
  ProjectDetailsSliceReducers,
  'projectDetails',
  any
>({
  name: 'projectDetails',
  initialState: getProjectDetailsInitialState(),
  reducers: {
    clearProjectDetailsState() {
      getProjectDetailsInitialState()
    },
  },

  extraReducers: (builder) => {
    builder
      .addCase(getProject.fulfilled, (state, { payload }) => {
        state.project = {
          ...payload,
        }
        state.isLoading = false
      })
      .addCase(getProject.rejected, (state) => {
        state.isLoading = false
      })
      .addCase(getProject.pending, (state) => {
        state.isLoading = true
      })
  },
})

declare module '@shared/store' {
  interface AppLazySlices extends WithSlice<typeof ProjectDetailsSlice> {}
}

ProjectDetailsSlice.injectInto(rootReducer)
export const { clearProjectDetailsState } = ProjectDetailsSlice.actions
