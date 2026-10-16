import { AccentTheme } from '@features/curriculumVitae/components/cv/AccentTheme'
import { ControlsBar } from '@features/curriculumVitae/components/cv/ControlsBar'
import { CursorGlow } from '@features/curriculumVitae/components/cv/CursorGlow'
import { FilterBar } from '@features/curriculumVitae/components/cv/FilterBar'
import { HeroSection } from '@features/curriculumVitae/components/cv/HeroSection'
import { ScrollProgress } from '@features/curriculumVitae/components/cv/ScrollProgress'
import { SectionCollapseProvider } from '@features/curriculumVitae/components/cv/SectionCollapseProvider'
import { SectionNav } from '@features/curriculumVitae/components/cv/SectionNav'
import { TechnologyFilterProvider } from '@features/curriculumVitae/components/cv/TechnologyFilterProvider'
import { TurnstileProvider } from '@features/curriculumVitae/components/cv/TurnstileProvider'
import { Box } from '@mui/material'
import { useCurriculumVitae } from './hooks'
import { Sections } from './Sections'

export const CurriculumVitae = () => {
  const { resume, sections, navItems, initiallyCollapsed, containerRef } = useCurriculumVitae()

  return (
    <AccentTheme colors={resume.overview.theme}>
      <SectionCollapseProvider initiallyCollapsed={initiallyCollapsed}>
        <TechnologyFilterProvider projects={resume.projects}>
          <TurnstileProvider>
            <Box
              ref={containerRef}
              sx={{ position: 'relative', isolation: 'isolate', width: '100%', overflowX: 'clip' }}
            >
              <ScrollProgress />
              <ControlsBar />
              <SectionNav items={navItems} />
              <CursorGlow />

              <HeroSection
                greeting={resume.overview.greeting}
                availability={resume.overview.availability}
                updatedAt={resume.updatedAt}
              />

              <Sections resume={resume} sections={sections} />

              <FilterBar experiences={resume.experiences} projects={resume.projects} />
            </Box>
          </TurnstileProvider>
        </TechnologyFilterProvider>
      </SectionCollapseProvider>
    </AccentTheme>
  )
}
