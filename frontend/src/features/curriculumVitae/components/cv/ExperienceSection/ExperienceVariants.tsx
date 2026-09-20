import { Section } from '@features/curriculumVitae/components/cv/Section'
import { TechnologyList } from '@features/curriculumVitae/components/cv/TechnologyList'
import { useTechnologyFilter } from '@features/curriculumVitae/hooks'
import { dimmed, iconTile, richCard } from '@features/curriculumVitae/styles'
import { Experience, SectionId } from '@features/curriculumVitae/types'
import { Box, Chip, Divider, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { formatMonth } from '@shared/utils'
import { Achievements } from './Achievements'
import { Dates } from './Dates'

type Props = {
  experience: Experience[]
}

export const ExperienceRich = ({ experience }: Props) => {
  const translate = useTranslate('CurriculumVitae.Experience')
  const { isDimmed } = useTechnologyFilter()

  if (!experience.length) return null

  return (
    <Section id={SectionId.Experience} title={translate('title')} collapsible={false}>
      <Stack className="reveal-group" spacing={2.5}>
        {experience.map((job) => (
          <Stack
            key={job.id}
            className="reveal-item"
            spacing={2}
            sx={{ p: { xs: 2.5, md: 3.5 }, ...richCard, ...dimmed(isDimmed(job.stack)) }}
          >
            <Stack direction="row" spacing={2} alignItems="center">
              <Box
                aria-hidden
                sx={{
                  display: 'grid',
                  placeItems: 'center',
                  width: 52,
                  height: 52,
                  borderRadius: 2,
                  flexShrink: 0,
                  fontSize: '1.5rem',
                  fontWeight: 800,
                  ...iconTile,
                }}
              >
                {job.company[0]}
              </Box>

              <Stack flex={1} minWidth={0}>
                <Typography variant="h5" sx={{ fontSize: { xs: '1.125rem', md: '1.375rem' } }}>
                  {job.position}
                </Typography>

                <Typography sx={{ color: 'primary.main', fontWeight: 600 }}>
                  {[job.company, job.location].filter(Boolean).join(' · ')}
                </Typography>
              </Stack>
            </Stack>

            <Box sx={{ alignSelf: 'flex-start' }}>
              <Chip size="small" variant="outlined" label={<Dates {...job} />} />
            </Box>

            {job.description && <Typography>{job.description}</Typography>}

            {!!job.tasks.length && <Achievements items={job.tasks} />}

            {!!job.achievements.length && <Achievements items={job.achievements} />}

            <TechnologyList technologies={job.stack} />
          </Stack>
        ))}
      </Stack>
    </Section>
  )
}

export const ExperienceCompact = ({ experience }: Props) => {
  const translate = useTranslate('CurriculumVitae.Experience')
  const { isDimmed } = useTechnologyFilter()

  if (!experience.length) return null

  return (
    <Section id={SectionId.Experience} title={translate('title')} collapsible={false}>
      <Stack className="reveal-group" divider={<Divider flexItem />} spacing={2}>
        {experience.map((job) => (
          <Box key={job.id} className="reveal-item" sx={dimmed(isDimmed(job.stack))}>
            <Stack direction="row" spacing={2} justifyContent="space-between" alignItems="baseline">
              <Typography sx={{ fontWeight: 700, minWidth: 0, overflowWrap: 'anywhere' }}>
                {job.position}
                <Typography component="span" sx={{ color: 'text.secondary', fontWeight: 400 }}>
                  {' · '}
                  {job.company}
                </Typography>
              </Typography>

              <Typography
                variant="caption"
                sx={{ color: 'gray.light', flexShrink: 0, whiteSpace: 'nowrap' }}
              >
                {formatMonth(job.startDate)} —{' '}
                {job.endDate ? formatMonth(job.endDate) : translate('present')}
              </Typography>
            </Stack>

            {!!job.stack.length && (
              <Typography variant="body2" sx={{ color: 'text.secondary', mt: 0.5 }}>
                {job.stack.map(({ name }) => name).join(', ')}
              </Typography>
            )}
          </Box>
        ))}
      </Stack>
    </Section>
  )
}
