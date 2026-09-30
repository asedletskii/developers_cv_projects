import { createLazySliceStateSelector } from '@shared/utils/createLazySliceStateSelector'
import { getSessionInitialState } from './session.state'

const sliceStateSelector = createLazySliceStateSelector('session', getSessionInitialState())

export const selectSessionUserId = sliceStateSelector((state) => state.userId)
export const selectSessionStatus = sliceStateSelector((state) => state.status)
