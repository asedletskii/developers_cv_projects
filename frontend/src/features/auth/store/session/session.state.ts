import { SessionSliceState } from '@features/auth/types'

export const getSessionInitialState = (): SessionSliceState => ({
  userId: null,
  status: 'idle',
})
