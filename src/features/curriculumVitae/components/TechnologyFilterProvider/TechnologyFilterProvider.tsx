import { createContext, PropsWithChildren, useCallback, useEffect, useMemo, useState } from 'react'
import { useSectionCollapse } from '@features/curriculumVitae/components/SectionCollapseProvider'
import { Person, SectionId, Technology } from '@features/curriculumVitae/types'
import { hasTechnology, isSameTechnology } from '@features/curriculumVitae/utils'

export type TechnologyFilterValue = {
  /** Выбранная технология или null, если фильтра нет */
  selected: string | null
  isSelected(name: string): boolean
  toggle(name: string): void
  clear(): void

  /** Нужно ли приглушить элемент со стеком: фильтр включён, а технологии в стеке нет */
  isDimmed(stack: Technology[]): boolean
}

// Значение по умолчанию: без провайдера фильтра нет, и всё выглядит как обычно
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

/**
 * Фильтр по технологии: клик выбирает её, повторный клик или Escape сбрасывает.
 * Если совпадения есть в свёрнутой секции проектов, она раскрывается.
 */
export const TechnologyFilterProvider = ({ projects, children }: Props) => {
  const [selected, setSelected] = useState<string | null>(null)
  const { expand } = useSectionCollapse()

  const clear = useCallback(() => setSelected(null), [])

  const toggle = useCallback(
    (name: string) => {
      const isSame = selected !== null && isSameTechnology(selected, name)
      setSelected(isSame ? null : name)

      if (!isSame && projects.some((project) => hasTechnology(project.stack, name))) {
        expand(SectionId.Projects)
      }
    },
    [selected, projects, expand],
  )

  useEffect(() => {
    if (!selected) return

    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') clear()
    }

    window.addEventListener('keydown', onKeyDown)
    return () => window.removeEventListener('keydown', onKeyDown)
  }, [selected, clear])

  const value = useMemo<TechnologyFilterValue>(
    () => ({
      selected,
      isSelected: (name) => selected !== null && isSameTechnology(selected, name),
      toggle,
      clear,
      isDimmed: (stack) => selected !== null && !hasTechnology(stack, selected),
    }),
    [selected, toggle, clear],
  )

  return (
    <TechnologyFilterContext.Provider value={value}>{children}</TechnologyFilterContext.Provider>
  )
}
