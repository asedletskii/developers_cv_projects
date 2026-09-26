import { createContext, PropsWithChildren } from 'react'
import { SectionCollapseValue, SectionId } from '@features/curriculumVitae/types'
import { useSectionCollapseState } from './hooks'

const noop = () => {}

// Значение по умолчанию: без провайдера секции просто всегда развёрнуты
export const SectionCollapseContext = createContext<SectionCollapseValue>({
  isCollapsed: () => false,
  toggle: noop,
  expand: noop,
})

type Props = PropsWithChildren<{
  /** Секции, свёрнутые при первом показе (приходят из настроек API) */
  initiallyCollapsed: SectionId[]
}>

/** Аккордеон секций: значение состояния отдаётся вниз по дереву через контекст. */
export const SectionCollapseProvider = ({ initiallyCollapsed, children }: Props) => {
  const value = useSectionCollapseState(initiallyCollapsed)

  return <SectionCollapseContext.Provider value={value}>{children}</SectionCollapseContext.Provider>
}
