import {
  CurriculumVitaeSliceReducers,
  CurriculumVitaeSliceState,
} from '@features/curriculumVitae/types'
import { createSlice, WithSlice } from '@reduxjs/toolkit'
import { rootReducer } from '@shared/store'
import { getCurriculumVitaeInitialState } from './curriculumVitae.state'
import { getPersonData } from './curriculumVitae.thunks'

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
      .addCase(getPersonData.fulfilled, (state, { payload, meta }) => {
        // Игнорируем ответ отменённого/устаревшего запроса: на экране должен остаться только
        // результат последнего из them, даже если старый ответ по сети пришёл позже нового
        if (meta.requestId !== state.currentRequestId) return

        state.person = {
          ...payload,
        }
        state.isLoading = false
        state.error = null
      })
      .addCase(getPersonData.rejected, (state, { payload, meta }) => {
        // Отменённый запрос — это не ошибка загрузки, его сменил более новый
        if (meta.aborted || meta.requestId !== state.currentRequestId) return

        state.isLoading = false
        state.error = payload ?? 'cannotGetPersonData'
      })
      .addCase(getPersonData.pending, (state, { meta }) => {
        state.currentRequestId = meta.requestId
        // Скелетон только когда на экране нет этого же резюме: смена языка не должна мигать
        state.isLoading = state.person.id !== meta.arg.cvId
        state.error = null
      })
  },
})

declare module '@shared/store' {
  interface AppLazySlices extends WithSlice<typeof curriculumVitaeSlice> {}
}

curriculumVitaeSlice.injectInto(rootReducer)
export const { clearCurriculumVitaeState } = curriculumVitaeSlice.actions
