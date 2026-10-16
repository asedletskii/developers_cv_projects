import { useCallback, useMemo, useState } from 'react'
import { SectionCollapseValue, SectionId } from '@features/curriculumVitae/types'

export const useSectionCollapseState = (initiallyCollapsed: SectionId[]): SectionCollapseValue => {
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

  return useMemo(
    () => ({ isCollapsed: (id: SectionId) => collapsed.has(id), toggle, expand }),
    [collapsed, toggle, expand],
  )
}
