import { createContext, PropsWithChildren, useCallback, useMemo, useState } from 'react'
import { SectionId } from '@features/curriculumVitae/types'

export type SectionCollapseValue = {
  isCollapsed(id: SectionId): boolean
  toggle(id: SectionId): void
  expand(id: SectionId): void
}

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

/** Состояние аккордеона. Начальное значение берётся из API, дальше управляет пользователь. */
export const SectionCollapseProvider = ({ initiallyCollapsed, children }: Props) => {
  const [collapsed, setCollapsed] = useState<Set<SectionId>>(() => new Set(initiallyCollapsed))

  const toggle = useCallback((id: SectionId) => {
    setCollapsed((current) => {
      const next = new Set(current)
      if (!next.delete(id)) next.add(id)
      return next
    })
  }, [])

  const expand = useCallback((id: SectionId) => {
    setCollapsed((current) => {
      if (!current.has(id)) return current

      const next = new Set(current)
      next.delete(id)
      return next
    })
  }, [])

  const value = useMemo(
    () => ({ isCollapsed: (id: SectionId) => collapsed.has(id), toggle, expand }),
    [collapsed, toggle, expand],
  )

  return <SectionCollapseContext.Provider value={value}>{children}</SectionCollapseContext.Provider>
}
