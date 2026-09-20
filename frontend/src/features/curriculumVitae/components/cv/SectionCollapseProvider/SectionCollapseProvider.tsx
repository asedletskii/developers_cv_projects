import { createContext, PropsWithChildren } from 'react'
import { SectionCollapseValue, SectionId } from '@features/curriculumVitae/types'
import { useSectionCollapseState } from './hooks'

const noop = () => {}

export const SectionCollapseContext = createContext<SectionCollapseValue>({
  isCollapsed: () => false,
  toggle: noop,
  expand: noop,
})

type Props = PropsWithChildren<{
  initiallyCollapsed: SectionId[]
}>

export const SectionCollapseProvider = ({ initiallyCollapsed, children }: Props) => {
  const value = useSectionCollapseState(initiallyCollapsed)

  return <SectionCollapseContext.Provider value={value}>{children}</SectionCollapseContext.Provider>
}
