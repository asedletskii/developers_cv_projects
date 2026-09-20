import { authApi } from '@features/auth/api'
import { Session } from '@features/auth/types'
import { createAsyncThunk } from '@reduxjs/toolkit'
import { RootState } from '@shared/types'
import axios from 'axios'
import { selectSessionStatus } from './session.selectors'

type ThunkConfig = { state: RootState; rejectValue: 'unauthenticated' | 'error' }

export const fetchSession = createAsyncThunk<Session, void, ThunkConfig>(
  'session/fetchSession',
  async (_, { rejectWithValue, signal }) => {
    try {
      return await authApi.getSession(signal)
    } catch (e) {
      if (axios.isAxiosError(e) && e.response?.status === 401) {
        return rejectWithValue('unauthenticated')
      }
      return rejectWithValue('error')
    }
  },
  {
    condition: (_, { getState }) => selectSessionStatus(getState()) === 'idle',
  },
)

export const logout = createAsyncThunk('session/logout', async () => {
  await authApi.logout()
})
