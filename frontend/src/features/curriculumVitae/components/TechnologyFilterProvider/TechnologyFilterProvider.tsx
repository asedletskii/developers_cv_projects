import { createContext, PropsWithChildren } from 'react'
import { Person, TechnologyFilterValue } from '@features/curriculumVitae/types'
import { useTechnologyFilterState } from './hooks'

export const TechnologyFilterContext = createContext<TechnologyFilterValue>({
  selected: null,
  isSelected: () => false,
  toggle: () => {},
  clear: () => {},
  isDimmed: () => false,
})

type Props = PropsWithChildren<{
  projects: Person['projects']
}>

export const TechnologyFilterProvider = ({ projects, children }: Props) => {
  const value = useTechnologyFilterState(projects)

  return (
    <TechnologyFilterContext.Provider value={value}>{children}</TechnologyFilterContext.Provider>
  )
}
