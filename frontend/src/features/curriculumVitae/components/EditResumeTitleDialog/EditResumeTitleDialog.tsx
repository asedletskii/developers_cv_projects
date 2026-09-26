import { glassCard } from '@features/curriculumVitae/styles'
import { LkResumeSummary } from '@features/curriculumVitae/types'
import { Button, Dialog, DialogActions, DialogContent, DialogTitle, TextField } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useEditResumeTitle } from './hooks'

type Props = {
  userId: string
  resume: LkResumeSummary
  onClose: () => void
}

export const EditResumeTitleDialog = ({ userId, resume, onClose }: Props) => {
  const translate = useTranslate('Lk.Resumes')
  const { title, setTitle, isSaving, hasError, canSave, save } = useEditResumeTitle({
    userId,
    resume,
    onClose,
  })

  return (
    <Dialog
      open
      onClose={onClose}
      fullWidth
      maxWidth="xs"
      slotProps={{ paper: { sx: { ...glassCard, backgroundImage: 'none' } } }}
    >
      <DialogTitle>{translate('editTitle')}</DialogTitle>

      <DialogContent>
        <TextField
          autoFocus
          fullWidth
          margin="dense"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          error={hasError}
          helperText={hasError ? translate('editTitleError') : undefined}
          disabled={isSaving}
        />
      </DialogContent>

      <DialogActions sx={{ px: 3, pb: 3 }}>
        <Button
          onClick={onClose}
          disabled={isSaving}
          sx={{ borderRadius: 999, px: 3, textTransform: 'none' }}
        >
          {translate('cancel')}
        </Button>

        <Button
          onClick={save}
          disabled={!canSave}
          variant="contained"
          sx={{ borderRadius: 999, px: 3, textTransform: 'none' }}
        >
          {translate('save')}
        </Button>
      </DialogActions>
    </Dialog>
  )
}
