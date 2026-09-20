import { combineSlices } from '@reduxjs/toolkit'
import { toastSlice } from './toast'

export interface AppLazySlices extends Record<string, object> {}

export const rootReducer = combineSlices(toastSlice).withLazyLoadedSlices<AppLazySlices>()
