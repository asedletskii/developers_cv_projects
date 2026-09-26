import { AccentTheme } from '@features/curriculumVitae/components/AccentTheme'
import { ControlsBar } from '@features/curriculumVitae/components/ControlsBar'
import { CursorGlow } from '@features/curriculumVitae/components/CursorGlow'
import { FilterBar } from '@features/curriculumVitae/components/FilterBar'
import { HeroSection } from '@features/curriculumVitae/components/HeroSection'
import { ScrollProgress } from '@features/curriculumVitae/components/ScrollProgress'
import { SectionCollapseProvider } from '@features/curriculumVitae/components/SectionCollapseProvider'
import { SectionNav } from '@features/curriculumVitae/components/SectionNav'
import { TechnologyFilterProvider } from '@features/curriculumVitae/components/TechnologyFilterProvider'
import { Box } from '@mui/material'
import { useCurriculumVitae } from './hooks'
import { Sections } from './Sections'

export const CurriculumVitae = () => {
  const { person, sections, navItems, initiallyCollapsed, containerRef } = useCurriculumVitae()

  return (
    <AccentTheme colors={person.theme}>
      <SectionCollapseProvider initiallyCollapsed={initiallyCollapsed}>
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
