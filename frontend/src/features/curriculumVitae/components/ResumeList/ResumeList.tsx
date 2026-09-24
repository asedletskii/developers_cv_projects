import { useState } from 'react'
import { glassCard } from '@features/curriculumVitae/styles'
import { LkResumeSummary } from '@features/curriculumVitae/types'
import { AddRounded, ChevronRightRounded, EditRounded } from '@mui/icons-material'
import { Button, ButtonBase, IconButton, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { Link as RouterLink } from 'react-router'
import { EditResumeTitleDialog } from './EditResumeTitleDialog'

type Props = {
  userId: string
  resumes: LkResumeSummary[]
}

/** Правая колонка ЛК: список резюме пользователя. */
export const ResumeList = ({ userId, resumes }: Props) => {
  const translate = useTranslate('Lk.Resumes')

  const [editingResume, setEditingResume] = useState<LkResumeSummary | null>(null)

  return (
    <Stack spacing={2}>
      <Stack direction="row" alignItems="center" justifyContent="space-between">
        <Typography variant="h5">{translate('title')}</Typography>

        <Button
          component={RouterLink}
          to={`/cv/lk/${encodeURIComponent(userId)}/new`}
          startIcon={<AddRounded />}
          size="small"
          sx={{ textTransform: 'none' }}
        >
          {translate('addResume')}
        </Button>
      </Stack>

      {resumes.length === 0 ? (
        <Typography sx={{ color: 'text.secondary' }}>{translate('empty')}</Typography>
      ) : (
        <Stack spacing={1.5}>
          {resumes.map((resume) => (
            <Stack
              key={resume.id}
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
                to={`/cv/lk/${encodeURIComponent(userId)}/${encodeURIComponent(resume.id)}`}
                sx={{ flex: 1, justifyContent: 'space-between', px: 3, py: 2 }}
              >
                <Typography>{resume.title}</Typography>

                <ChevronRightRounded sx={{ color: 'gray.light' }} />
              </ButtonBase>

              <IconButton
                aria-label={translate('editTitle')}
                onClick={() => setEditingResume(resume)}
              >
                <EditRounded fontSize="small" />
              </IconButton>
            </Stack>
          ))}
        </Stack>
      )}

      {editingResume && (
        <EditResumeTitleDialog
          userId={userId}
          resume={editingResume}
          onClose={() => setEditingResume(null)}
        />
      )}
    </Stack>
  )
}
