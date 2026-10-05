import { ExternalLink } from '@features/curriculumVitae/components/ExternalLink'
import { TechnologyList } from '@features/curriculumVitae/components/TechnologyList'
import { TiltCard } from '@features/curriculumVitae/components/TiltCard'
import { useTechnologyFilter } from '@features/curriculumVitae/hooks'
import { dimmed, glassCard } from '@features/curriculumVitae/styles'
import { Resume } from '@features/curriculumVitae/types'
import { Box, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'

type Props = {
  project: Resume['projects'][number]
}

export const Project = ({ project }: Props) => {
  const translate = useTranslate('CurriculumVitae.Projects')
  const { isDimmed } = useTechnologyFilter()

  return (
    <Box className="reveal-item" sx={{ display: 'flex' }}>
      <TiltCard
        sx={{
          flex: 1,
          display: 'flex',
          flexDirection: 'column',
          gap: 2,
          p: 3,
          ...glassCard,
          ...dimmed(isDimmed(project.stack)),
          transition: 'box-shadow .25s, border-color .25s, filter .3s',
          '&:hover': { boxShadow: 2, borderColor: 'primary.light' },
        }}
      >
        <Stack spacing={1} flex={1}>
          <Typography variant="h5" sx={{ fontSize: '1.375rem' }}>
            {project.name}
          </Typography>

          {project.description && (
            <Typography sx={{ color: 'text.secondary' }}>{project.description}</Typography>
          )}
        </Stack>

        <TechnologyList technologies={project.stack} />

        {project.link && (
          <ExternalLink href={project.link} kind="project">
            {translate('open')}
          </ExternalLink>
        )}
      </TiltCard>
    </Box>
  )
}
