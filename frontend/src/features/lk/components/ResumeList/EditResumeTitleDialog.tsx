import { useState } from 'react'
import { glassCard } from '@features/curriculumVitae/styles'
import { LkResumeSummary } from '@features/lk/types'
import { updateLkResumeTitle } from '@features/lk/store/lk'
import { Button, Dialog, DialogActions, DialogContent, DialogTitle, TextField } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useAppDispatch } from '@shared/store'

type Props = {
  userId: string
  resume: LkResumeSummary
  onClose: () => void
}

export const EditResumeTitleDialog = ({ userId, resume, onClose }: Props) => {
  const dispatch = useAppDispatch()
  const translate = useTranslate('Lk.Resumes')

  const [title, setTitle] = useState(resume.title)
  const [isSaving, setIsSaving] = useState(false)
  const [error, setError] = useState(false)

  const handleSave = async () => {
    setIsSaving(true)
    setError(false)

    try {
      await dispatch(updateLkResumeTitle({ userId, resumeId: resume.id, title })).unwrap()
      onClose()
    } catch {
      setError(true)
      setIsSaving(false)
    }
  }

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
          error={error}
          helperText={error ? translate('editTitleError') : undefined}
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
          onClick={handleSave}
          disabled={isSaving || !title.trim()}
          variant="contained"
          sx={{ borderRadius: 999, px: 3, textTransform: 'none' }}
        >
          {translate('save')}
        </Button>
      </DialogActions>
    </Dialog>
  )
}
