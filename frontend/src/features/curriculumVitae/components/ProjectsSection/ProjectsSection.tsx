import { Section } from '@features/curriculumVitae/components/Section'
import { Resume, SectionId } from '@features/curriculumVitae/types'
import { Box } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { Project } from './Project'

type Props = {
  projects: Resume['projects']
}

export const ProjectsSection = ({ projects }: Props) => {
  const translate = useTranslate('CurriculumVitae.Projects')

  if (!projects.length) return null

  return (
    <Section id={SectionId.Projects} title={translate('title')}>
      <Box
        className="reveal-group"
        sx={{
          display: 'grid',
          gap: 2.5,
          gridTemplateColumns: { xs: '1fr', md: 'repeat(2, 1fr)' },
        }}
      >
        {projects.map((project) => (
          <Project key={project.id} project={project} />
        ))}
      </Box>
    </Section>
  )
}
