import { CreateResumeDialog } from '@features/curriculumVitae/components/lk/CreateResumeDialog'
import { EditResumeTitleDialog } from '@features/curriculumVitae/components/lk/EditResumeTitleDialog'
import { LkResumeSummary } from '@features/curriculumVitae/types'
import { AddRounded } from '@mui/icons-material'
import { Button, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { EmptyResumes } from './EmptyResumes'
import { useResumeList } from './hooks'
import { ResumeActionsMenu } from './ResumeActionsMenu'
import { ResumeListItem } from './ResumeListItem'

type Props = {
  resumes: LkResumeSummary[]
}

export const ResumeList = ({ resumes }: Props) => {
  const translate = useTranslate('Lk.Resumes')
  const list = useResumeList()

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
        <EmptyResumes onCreate={list.openCreating} />
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

      {list.isCreating && <CreateResumeDialog onClose={list.closeCreating} />}

      {list.editingResume && (
        <EditResumeTitleDialog resume={list.editingResume} onClose={list.closeEditing} />
      )}
    </Stack>
  )
}
