import { createSlice, WithSlice } from '@reduxjs/toolkit'
import { rootReducer } from '@shared/store'
import { getSessionInitialState } from './session.state'
import { fetchSession, logout } from './session.thunks'

export const sessionSlice = createSlice({
  name: 'session',
  initialState: getSessionInitialState(),
  reducers: {
    sessionExpired(state) {
      state.status = 'unauthenticated'
      state.userId = null
    },
  },

  extraReducers: (builder) => {
    builder
      .addCase(fetchSession.pending, (state) => {
        state.status = 'loading'
      })
      .addCase(fetchSession.fulfilled, (state, { payload }) => {
        state.status = 'authenticated'
        state.userId = payload.userId
      })
      .addCase(fetchSession.rejected, (state, { payload }) => {
        state.status = payload === 'unauthenticated' ? 'unauthenticated' : 'error'
        state.userId = null
      })

      .addCase(logout.fulfilled, () => getSessionInitialState())
  },
})

declare module '@shared/store' {
  interface AppLazySlices extends WithSlice<typeof sessionSlice> {}
}

sessionSlice.injectInto(rootReducer)
export const { sessionExpired } = sessionSlice.actions
