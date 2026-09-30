export type Session = {
  userId: string
  authenticated: boolean
}

export type SessionStatus = 'idle' | 'loading' | 'authenticated' | 'unauthenticated' | 'error'

export type SessionSliceState = {
  userId: string | null
  status: SessionStatus
}
