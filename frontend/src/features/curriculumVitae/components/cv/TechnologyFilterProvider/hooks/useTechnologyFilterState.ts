import { useCallback, useEffect, useMemo, useState } from 'react'
import { useAnalytics } from '@features/curriculumVitae/hooks/useAnalytics'
import { useSectionCollapse } from '@features/curriculumVitae/hooks/useSectionCollapse'
import {
  AnalyticsEventType,
  Resume,
  SectionId,
  TechnologyFilterValue,
} from '@features/curriculumVitae/types'
import { hasTechnology, isSameTechnology } from '@features/curriculumVitae/utils'

export const useTechnologyFilterState = (projects: Resume['projects']): TechnologyFilterValue => {
  const [selected, setSelected] = useState<string | null>(null)
  const { expand } = useSectionCollapse()
  const { track } = useAnalytics()

  const clear = useCallback(() => {
    if (selected) {
      track({ type: AnalyticsEventType.TechFilter, technology: selected, action: 'clear' })
    }
    setSelected(null)
  }, [selected, track])

  const toggle = useCallback(
    (name: string) => {
      const isSame = selected !== null && isSameTechnology(selected, name)
      setSelected(isSame ? null : name)
      track({
        type: AnalyticsEventType.TechFilter,
        technology: name,
        action: isSame ? 'clear' : 'select',
      })

      if (!isSame && projects.some((project) => hasTechnology(project.stack, name))) {
        expand(SectionId.Projects)
      }
    },
    [selected, projects, expand, track],
  )

  useEffect(() => {
    if (!selected) return

    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') clear()
    }

    window.addEventListener('keydown', onKeyDown)
    return () => window.removeEventListener('keydown', onKeyDown)
  }, [selected, clear])

  return useMemo(
    () => ({
      selected,
      isSelected: (name) => selected !== null && isSameTechnology(selected, name),
      toggle,
      clear,
      isDimmed: (stack) => selected !== null && !hasTechnology(stack, selected),
    }),
    [selected, toggle, clear],
  )
}
