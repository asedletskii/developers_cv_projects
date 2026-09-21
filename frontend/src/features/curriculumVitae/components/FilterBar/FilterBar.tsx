import { useEffect, useState } from 'react'
import { CloseRounded } from '@mui/icons-material'
import { Chip, IconButton, Stack, Typography } from '@mui/material'
import { useTechnologyFilter } from '@features/curriculumVitae/components/TechnologyFilterProvider'
import { hasTechnology } from '@features/curriculumVitae/utils'
import { glassCard } from '@features/curriculumVitae/styles'
import { Person } from '@features/curriculumVitae/types'
import { useTranslate } from '@shared/hooks'

type Props = Pick<Person, 'experience' | 'projects'>

/** Плавающая плашка активного фильтра: технология, сколько совпадений и кнопка сброса. */
export const FilterBar = ({ experience, projects }: Props) => {
  const translate = useTranslate('CurriculumVitae.Filter')
  const { selected, clear } = useTechnologyFilter()

  // Помним последнюю технологию, чтобы плашка не «схлопывалась» текстом, пока исчезает
  const [shown, setShown] = useState(selected)

  useEffect(() => {
    if (selected) setShown(selected)
  }, [selected])

  const isVisible = !!selected

  const jobs = shown ? experience.filter((job) => hasTechnology(job.stack, shown)).length : 0
  const matchedProjects = shown
    ? projects.filter((project) => hasTechnology(project.stack, shown)).length
    : 0

  return (
    <Stack
      data-print-hidden
      role="status"
      direction="row"
      alignItems="center"
      spacing={1.5}
      sx={{
        ...glassCard,
        position: 'fixed',
        bottom: 16,
        left: '50%',
        zIndex: 1100,
        maxWidth: 'calc(100% - 24px)',
        pl: 1.5,
        pr: 0.5,
        py: 0.5,
        borderRadius: 999,
        opacity: isVisible ? 1 : 0,
        visibility: isVisible ? 'visible' : 'hidden',
        transform: `translate(-50%, ${isVisible ? 0 : 16}px)`,
        transition: isVisible
          ? 'opacity .3s, transform .4s cubic-bezier(0.32, 0.72, 0, 1)'
          : 'opacity .2s, transform .3s, visibility 0s linear .3s',
        '@media (prefers-reduced-motion: reduce)': { transition: 'none' },
      }}
    >
      <Typography variant="caption" sx={{ color: 'text.secondary' }}>
        {translate('technology')}
      </Typography>

      <Chip size="small" color="primary" label={shown} />

      <Typography variant="caption" sx={{ color: 'text.secondary', whiteSpace: 'nowrap' }}>
        {translate('experience')}: {jobs} · {translate('projects')}: {matchedProjects}
      </Typography>

      <IconButton size="small" aria-label={translate('reset')} onClick={clear}>
        <CloseRounded fontSize="small" />
      </IconButton>
    </Stack>
  )
}
