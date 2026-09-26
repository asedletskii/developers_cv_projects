import { useEffect, useMemo } from 'react'
import { getVisibleSections, SECTION_DEFINITIONS } from '@features/curriculumVitae/sections'
import { selectPersonData } from '@features/curriculumVitae/store/curriculumVitae'
import { ScrollTrigger } from '@shared/animations'
import { useDocumentTitle, useLanguage, usePrintLightTheme, useScrollReveal } from '@shared/hooks'
import { useAppSelector } from '@shared/store'
import { getPersonFIO } from '@shared/utils'
import { useNavItems } from './useNavItems'

export const useCurriculumVitae = () => {
  const person = useAppSelector(selectPersonData)
  const { language } = useLanguage()
  const containerRef = useScrollReveal()

  usePrintLightTheme()
  useDocumentTitle(getPersonFIO(person))

  const sections = useMemo(() => getVisibleSections(person), [person])
  const navItems = useNavItems(sections)

  const initiallyCollapsed = sections
    .filter(({ id, collapsed }) => collapsed && SECTION_DEFINITIONS[id].isCollapsible)
    .map(({ id }) => id)

  //от языка зависит длина текстов, а значит и высота страницы — пересчитываем позиции триггеров
  useEffect(() => {
    const timer = window.setTimeout(() => ScrollTrigger.refresh(), 100)
    return () => window.clearTimeout(timer)
  }, [language])

  return { person, sections, navItems, initiallyCollapsed, containerRef }
}
