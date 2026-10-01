import { TechnologyList } from '@features/curriculumVitae/components/TechnologyList'
import { useTechnologyFilter } from '@features/curriculumVitae/hooks'
import { dimmed, glassCard } from '@features/curriculumVitae/styles'
import { Experience } from '@features/curriculumVitae/types'
import { Box, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { Achievements } from './Achievements'
import { Dates } from './Dates'

type Props = {
  job: Experience
  isExtra: boolean
}

export const Job = ({ job, isExtra }: Props) => {
  const translate = useTranslate('CurriculumVitae.Experience')
  const { isDimmed } = useTechnologyFilter()

  return (
    <Box className={isExtra ? 'timeline-extra' : 'timeline-item'} sx={{ position: 'relative' }}>
      <Box
        className="timeline-dot"
        aria-hidden
        sx={{
          position: 'absolute',
          top: 6,
          left: { xs: -32, sm: -44 },
          width: 20,
          height: 20,
          borderRadius: '50%',
          bgcolor: 'primary.main',
          border: 4,
          borderColor: 'surface.main',
          boxShadow: ({ palette }) => `0 0 0 3px ${palette.primary.light}`,
        }}
      />

      <Stack
        className="timeline-card"
        spacing={1.5}
        sx={{ p: { xs: 2.5, sm: 3 }, ...glassCard, ...dimmed(isDimmed(job.stack)) }}
      >
        <Stack spacing={0.5}>
          <Typography variant="h5" sx={{ fontSize: '1.25rem' }}>
            {job.position}
          </Typography>

          <Typography sx={{ color: 'primary.main', fontWeight: 600 }}>
            {[job.company, job.location].filter(Boolean).join(' · ')}
          </Typography>

          <Dates startDate={job.startDate} endDate={job.endDate} />
        </Stack>

        {job.description && (
          <Typography sx={{ color: 'text.secondary' }}>{job.description}</Typography>
        )}

        {!!job.tasks.length && (
          <Stack spacing={0.5}>
            <Typography variant="subtitle2" sx={{ color: 'text.secondary' }}>
              {translate('tasks')}
            </Typography>
            <Achievements items={job.tasks} />
          </Stack>
        )}

        {!!job.achievements.length && (
          <Stack spacing={0.5}>
            <Typography variant="subtitle2" sx={{ color: 'text.secondary' }}>
              {translate('achievements')}
            </Typography>
            <Achievements items={job.achievements} />
          </Stack>
        )}

        <TechnologyList technologies={job.stack} />
      </Stack>
    </Box>
  )
}
