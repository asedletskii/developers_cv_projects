import { useEffect, useMemo } from 'react'
import { AccentTheme } from '@features/curriculumVitae/components/AccentTheme'
import { ControlsBar } from '@features/curriculumVitae/components/ControlsBar'
import { CursorGlow } from '@features/curriculumVitae/components/CursorGlow'
import { FilterBar } from '@features/curriculumVitae/components/FilterBar'
import { HeroSection } from '@features/curriculumVitae/components/HeroSection'
import { ScrollProgress } from '@features/curriculumVitae/components/ScrollProgress'
import { SectionCollapseProvider } from '@features/curriculumVitae/components/SectionCollapseProvider'
import { SectionNav } from '@features/curriculumVitae/components/SectionNav'
import { TechnologyFilterProvider } from '@features/curriculumVitae/components/TechnologyFilterProvider'
import { usePrintLightTheme, useScrollReveal } from '@features/curriculumVitae/hooks'
import { getVisibleSections, SECTION_DEFINITIONS } from '@features/curriculumVitae/sections'
import { selectPersonData } from '@features/curriculumVitae/store/curriculumVitae'
import { Box } from '@mui/material'
import { ScrollTrigger } from '@shared/animations'
import { useDocumentTitle, useLanguage } from '@shared/hooks'
import { useAppSelector } from '@shared/store'
import { getPersonFIO } from '@shared/utils'
import { Sections } from './Sections'
import { useNavItems } from './useNavItems'

/** Готовое резюме: hero, плавающие элементы управления и секции в заданном порядке. */
export const CurriculumVitae = () => {
  const person = useAppSelector(selectPersonData)
  const { language } = useLanguage()
  const containerRef = useScrollReveal()

  usePrintLightTheme()
  useDocumentTitle(getPersonFIO(person))

  const sections = useMemo(() => getVisibleSections(person), [person])
  const navItems = useNavItems(sections)

  // От языка зависит длина текстов, а значит и высота страницы — пересчитываем позиции триггеров
  useEffect(() => {
    const timer = window.setTimeout(() => ScrollTrigger.refresh(), 100)
    return () => window.clearTimeout(timer)
  }, [language])

  return (
    <AccentTheme colors={person.theme}>
      <SectionCollapseProvider
        initiallyCollapsed={sections
          .filter(({ id, collapsed }) => collapsed && SECTION_DEFINITIONS[id].isCollapsible)
          .map(({ id }) => id)}
      >
        <TechnologyFilterProvider projects={person.projects}>
          <Box
            ref={containerRef}
            sx={{ position: 'relative', isolation: 'isolate', width: '100%', overflowX: 'clip' }}
          >
            <ScrollProgress />
            <ControlsBar resumeLink={person.resumeLink} />
            <SectionNav items={navItems} />
            <CursorGlow />

            <HeroSection
              greeting={person.greeting}
              resumeLink={person.resumeLink}
              availability={person.availability}
              updatedAt={person.updatedAt}
            />

            <Sections person={person} sections={sections} />

            <FilterBar experience={person.experience} projects={person.projects} />
          </Box>
        </TechnologyFilterProvider>
      </SectionCollapseProvider>
    </AccentTheme>
  )
}
