import { CreateResumeDialog } from '@features/curriculumVitae/components/CreateResumeDialog'
import { EditResumeTitleDialog } from '@features/curriculumVitae/components/EditResumeTitleDialog'
import { LkResumeSummary } from '@features/curriculumVitae/types'
import { AddRounded } from '@mui/icons-material'
import { Button, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useResumeList } from './hooks'
import { ResumeActionsMenu } from './ResumeActionsMenu'
import { ResumeListItem } from './ResumeListItem'

type Props = {
  userId: string
  resumes: LkResumeSummary[]
}

export const ResumeList = ({ userId, resumes }: Props) => {
  const translate = useTranslate('Lk.Resumes')
  const list = useResumeList(userId)

  return (
    <Stack spacing={2}>
      <Stack direction="row" alignItems="center" justifyContent="space-between">
        <Typography variant="h5">{translate('title')}</Typography>

        <Button
          onClick={list.openCreating}
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
            <ResumeListItem
              key={resume.id}
              resume={resume}
              onOpenMenu={(e) => list.openMenu(e.currentTarget, resume)}
            />
          ))}
        </Stack>
      )}

      <ResumeActionsMenu
        anchor={list.menu?.anchor ?? null}
        status={list.menu?.resume.status ?? 'draft'}
        onClose={list.closeMenu}
        onRename={list.renameFromMenu}
        onPublish={list.publishFromMenu}
        onUnpublish={list.unpublishFromMenu}
        onRemove={list.removeFromMenu}
      />

      {list.isCreating && <CreateResumeDialog userId={userId} onClose={list.closeCreating} />}

      {list.editingResume && (
        <EditResumeTitleDialog
          userId={userId}
          resume={list.editingResume}
          onClose={list.closeEditing}
        />
      )}
    </Stack>
  )
}
