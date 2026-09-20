import { LkSliceReducers, LkSliceState } from '@features/curriculumVitae/types'
import { createSlice, WithSlice } from '@reduxjs/toolkit'
import { rootReducer } from '@shared/store'
import { getLkInitialState } from './lk.state'
import {
  changeLkResumeStatus,
  deleteLkResume,
  getLkProfile,
  getLkResumes,
  updateLkProfile,
  updateLkResumeTitle,
} from './lk.thunks'

export const lkSlice = createSlice<LkSliceState, LkSliceReducers, 'lk', any>({
  name: 'lk',
  initialState: getLkInitialState(),
  reducers: {
    clearLkState() {
      return getLkInitialState()
    },
  },

  extraReducers: (builder) => {
    builder
      .addCase(getLkProfile.fulfilled, (state, { payload, meta }) => {
        if (meta.requestId !== state.currentProfileRequestId) return

        state.profile = payload
        state.isProfileLoading = false
        state.profileError = null
      })
      .addCase(getLkProfile.rejected, (state, { payload, meta }) => {
        if (meta.aborted || meta.requestId !== state.currentProfileRequestId) return

        state.isProfileLoading = false
        state.profileError = payload ?? { kind: 'cannotGetPersonData' }
      })
      .addCase(getLkProfile.pending, (state, { meta }) => {
        state.currentProfileRequestId = meta.requestId
        state.isProfileLoading = state.profile?.userId !== meta.arg
        state.profileError = null
      })

      .addCase(getLkResumes.fulfilled, (state, { payload, meta }) => {
        if (meta.requestId !== state.currentResumesRequestId) return

        state.resumes = payload
        state.isResumesLoading = false
        state.resumesError = null
      })
      .addCase(getLkResumes.rejected, (state, { payload, meta }) => {
        if (meta.aborted || meta.requestId !== state.currentResumesRequestId) return

        state.isResumesLoading = false
        state.resumesError = payload ?? { kind: 'cannotGetPersonData' }
      })
      .addCase(getLkResumes.pending, (state, { meta }) => {
        state.currentResumesRequestId = meta.requestId
        state.isResumesLoading = state.profile?.userId !== meta.arg
        state.resumesError = null
      })

      .addCase(updateLkProfile.fulfilled, (state, { payload }) => {
        state.profile = payload
        state.profileError = null
      })

      .addCase(deleteLkResume.fulfilled, (state, { payload }) => {
        if (state.resumes) state.resumes = state.resumes.filter((r) => r.id !== payload)
      })
      .addCase(updateLkResumeTitle.fulfilled, (state, { payload }) => {
        const resume = state.resumes?.find((r) => r.id === payload.id)
        if (resume) {
          resume.title = payload.title
          resume.version = payload.version
        }
      })
      .addCase(changeLkResumeStatus.fulfilled, (state, { payload }) => {
        const resume = state.resumes?.find((r) => r.id === payload.id)
        if (resume) {
          resume.status = payload.status
          resume.version = payload.version
        }
      })
  },
})

declare module '@shared/store' {
  interface AppLazySlices extends WithSlice<typeof lkSlice> {}
}

lkSlice.injectInto(rootReducer)
export const { clearLkState } = lkSlice.actions
