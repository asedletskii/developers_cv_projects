import { getCurriculumVitaeInitialState } from '@features/curriculumVitae/store/curriculumVitae/curriculumVitae.state'
import { RootState } from '@shared/types'
import { createLazySliceStateSelector } from './createLazySliceStateSelector'

const defaultState = getCurriculumVitaeInitialState()

const createSelector = () =>
  createLazySliceStateSelector('curriculumVitae', defaultState)((sliceState, rootState) => ({
    sliceState,
    rootState,
  }))

describe('createLazySliceStateSelector', () => {
  it('отдаёт состояние среза, если он уже подключён', () => {
    const sliceState = { ...defaultState, isLoading: false }
    const rootState = { curriculumVitae: sliceState } as unknown as RootState

    expect(createSelector()(rootState).sliceState).toBe(sliceState)
  })

  it('пока срез не подключён, отдаёт состояние по умолчанию', () => {
    const rootState = {} as RootState

    expect(createSelector()(rootState).sliceState).toBe(defaultState)
  })

  it('передаёт в селектор и корневое состояние', () => {
    const rootState = {} as RootState

    expect(createSelector()(rootState).rootState).toBe(rootState)
  })

  it('возвращает то, что вернул селектор', () => {
    const selector = createLazySliceStateSelector('curriculumVitae', defaultState)(
      (sliceState) => sliceState.isLoading,
    )

    expect(selector({} as RootState)).toBe(defaultState.isLoading)
  })
})
