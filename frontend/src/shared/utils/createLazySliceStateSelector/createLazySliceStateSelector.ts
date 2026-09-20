import { AppLazySlices } from '@shared/store'
import { RootState } from '@shared/types'

export const createLazySliceStateSelector = <SliceName extends keyof AppLazySlices>(
  sliceName: SliceName,
  defaultState: AppLazySlices[SliceName],
) => {
  const lazySliceStateSelector = <T>(
    selector: (sliceState: AppLazySlices[SliceName], rootState: RootState) => T,
  ) => {
    return (rootState: RootState) => selector(rootState[sliceName] ?? defaultState, rootState)
  }

  return lazySliceStateSelector
}
