import {
  CurriculumVitaeSliceReducers,
  CurriculumVitaeSliceState,
} from '@features/curriculumVitae/types'
import { createSlice, WithSlice } from '@reduxjs/toolkit'
import { rootReducer } from '@shared/store'
import { getCurriculumVitaeInitialState } from './curriculumVitae.state'
import { getPublicResume } from './curriculumVitae.thunks'

export const curriculumVitaeSlice = createSlice<
  CurriculumVitaeSliceState,
  CurriculumVitaeSliceReducers,
  'curriculumVitae',
  any
>({
  name: 'curriculumVitae',
  initialState: getCurriculumVitaeInitialState(),
  reducers: {
    clearCurriculumVitaeState() {
      return getCurriculumVitaeInitialState()
    },
  },

  extraReducers: (builder) => {
    builder
      .addCase(getPublicResume.fulfilled, (state, { payload, meta }) => {
        if (meta.requestId !== state.currentRequestId) return

        state.resume = payload.resume
        state.viewId = payload.viewId
        state.isLoading = false
        state.error = null
      })
      .addCase(getPublicResume.rejected, (state, { payload, meta }) => {
        if (meta.aborted || meta.requestId !== state.currentRequestId) return

        state.isLoading = false
        state.error = payload ?? { kind: 'cannotGetPersonData' }
      })
      .addCase(getPublicResume.pending, (state, { meta }) => {
        state.currentRequestId = meta.requestId
        state.isLoading = state.resume.id !== meta.arg.cvId
        state.error = null
      })
  },
})

declare module '@shared/store' {
  interface AppLazySlices extends WithSlice<typeof curriculumVitaeSlice> {}
}

curriculumVitaeSlice.injectInto(rootReducer)
export const { clearCurriculumVitaeState } = curriculumVitaeSlice.actions
