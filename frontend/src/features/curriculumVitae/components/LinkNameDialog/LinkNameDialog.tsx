import { glassCard } from '@features/curriculumVitae/styles'
import { Button, Dialog, DialogActions, DialogContent, DialogTitle, TextField } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useLinkNameDialog } from './hooks'

type Props = {
  title: string
  initialName?: string
  onSave: (name: string) => void
  onClose: () => void
}

export const LinkNameDialog = ({ title, initialName, onSave, onClose }: Props) => {
  const translate = useTranslate('Lk.Links')
  const { name, setName, canSave, save, handleKeyDown } = useLinkNameDialog({
    initialName,
    onSave,
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
      <DialogTitle>{title}</DialogTitle>

      <DialogContent>
        <TextField
          autoFocus
          fullWidth
          margin="dense"
          label={translate('name')}
          placeholder={translate('namePlaceholder')}
          value={name}
          onChange={(e) => setName(e.target.value)}
          onKeyDown={handleKeyDown}
        />
      </DialogContent>

      <DialogActions sx={{ px: 3, pb: 3 }}>
        <Button onClick={onClose} sx={{ borderRadius: 999, px: 3, textTransform: 'none' }}>
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
