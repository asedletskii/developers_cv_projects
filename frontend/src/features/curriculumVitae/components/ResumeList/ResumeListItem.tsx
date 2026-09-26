import { glassCard } from '@features/curriculumVitae/styles'
import { LkResumeSummary } from '@features/curriculumVitae/types'
import { ChevronRightRounded, MoreVertRounded } from '@mui/icons-material'
import { ButtonBase, Chip, IconButton, Stack, Typography } from '@mui/material'
import { ROUTES } from '@shared/constants'
import { useTranslate } from '@shared/hooks'
import { MouseEvent } from 'react'
import { generatePath, Link as RouterLink } from 'react-router'

type Props = {
  userId: string
  resume: LkResumeSummary
  onOpenMenu: (e: MouseEvent<HTMLElement>) => void
}

/** Строка списка: название и язык резюме (ссылка в редактор) и меню действий. */
export const ResumeListItem = ({ userId, resume, onOpenMenu }: Props) => {
  const translate = useTranslate('Lk.Resumes')

  return (
    <Stack
      direction="row"
      alignItems="center"
      sx={{
        pr: 1,
        ...glassCard,
        transition: 'box-shadow .2s, border-color .2s',
        '&:hover': { boxShadow: 2, borderColor: 'primary.light' },
      }}
    >
      <ButtonBase
        component={RouterLink}
        to={generatePath(ROUTES.lkResume, { userId, cvId: resume.id })}
        sx={{ flex: 1, justifyContent: 'space-between', px: 3, py: 2 }}
      >
        <Stack direction="row" alignItems="center" spacing={1.5}>
          <Typography>{resume.title}</Typography>

          <Chip label={resume.language.toUpperCase()} size="small" />
        </Stack>

        <ChevronRightRounded sx={{ color: 'gray.light' }} />
      </ButtonBase>

      <IconButton aria-label={translate('actions')} onClick={onOpenMenu}>
        <MoreVertRounded fontSize="small" />
      </IconButton>
    </Stack>
  )
}
