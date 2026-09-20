import { useEffect, useMemo } from 'react'
import { EMPTY_PROFILE } from '@features/curriculumVitae/constants'
import { getVisibleSections, SECTION_DEFINITIONS } from '@features/curriculumVitae/sections'
import { selectResume } from '@features/curriculumVitae/store/curriculumVitae'
import { ScrollTrigger } from '@shared/animations'
import { useDocumentTitle, useLanguage, usePrintLightTheme, useScrollReveal } from '@shared/hooks'
import { useAppSelector } from '@shared/store'
import { getPersonFIO } from '@shared/utils'
import { useNavItems } from './useNavItems'

export const useCurriculumVitae = () => {
  const resume = useAppSelector(selectResume)
  const { language } = useLanguage()
  const containerRef = useScrollReveal()

  usePrintLightTheme()
  useDocumentTitle(getPersonFIO(resume.profile ?? EMPTY_PROFILE))

  const sections = useMemo(() => getVisibleSections(resume), [resume])
  const navItems = useNavItems(sections)

  const initiallyCollapsed = sections
    .filter(({ id, collapsed }) => collapsed && SECTION_DEFINITIONS[id].isCollapsible)
    .map(({ id }) => id)

  //от языка зависит длина текстов, а значит и высота страницы — пересчитываем позиции триггеров
  useEffect(() => {
    const timer = window.setTimeout(() => ScrollTrigger.refresh(), 100)
    return () => window.clearTimeout(timer)
  }, [language])

  return { resume, sections, navItems, initiallyCollapsed, containerRef }
}
