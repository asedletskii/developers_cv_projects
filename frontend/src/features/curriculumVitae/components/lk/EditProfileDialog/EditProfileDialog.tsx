import { dialogPaper } from '@features/curriculumVitae/styles'
import { LkProfile } from '@features/curriculumVitae/types'
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Stack,
  Tab,
  Tabs,
  TextField,
  Typography,
} from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { Languages } from '@shared/types'
import { useEditProfile } from './hooks'

type Props = {
  userId: string
  profile: LkProfile
  onClose: () => void
  intro?: string
}

export const EditProfileDialog = ({ userId, profile, onClose, intro }: Props) => {
  const translate = useTranslate('Lk.Profile')
  const edit = useEditProfile({ userId, profile, onClose })

  return (
    <Dialog
      open
      onClose={onClose}
      fullWidth
      maxWidth="xs"
      slotProps={{ paper: { sx: dialogPaper } }}
    >
      <DialogTitle>{translate('editProfile')}</DialogTitle>

      <DialogContent>
        <Stack spacing={2}>
          {intro && (
            <Typography variant="body2" sx={{ color: 'text.secondary' }}>
              {intro}
            </Typography>
          )}

          <TextField
            fullWidth
            label={translate('photoPath')}
            value={edit.photoPath}
            onChange={(e) => edit.setPhotoPath(e.target.value)}
            disabled={edit.isSaving}
          />

          <Tabs value={edit.language} onChange={(_, value) => edit.setLanguage(value)}>
            <Tab value={Languages.Ru} label="RU" />
            <Tab value={Languages.En} label="EN" />
          </Tabs>

          <TextField
            fullWidth
            label={translate('name')}
            value={edit.current.name ?? ''}
            onChange={(e) => edit.updateField('name', e.target.value)}
            disabled={edit.isSaving}
          />

          <TextField
            fullWidth
            label={translate('surname')}
            value={edit.current.surname ?? ''}
            onChange={(e) => edit.updateField('surname', e.target.value)}
            disabled={edit.isSaving}
          />

          <TextField
            fullWidth
            label={translate('patronymic')}
            value={edit.current.patronymic ?? ''}
            onChange={(e) => edit.updateField('patronymic', e.target.value)}
            disabled={edit.isSaving}
          />

          <TextField
            fullWidth
            label={translate('position')}
            value={edit.current.position ?? ''}
            onChange={(e) => edit.updateField('position', e.target.value)}
            disabled={edit.isSaving}
          />

          <TextField
            fullWidth
            label={translate('location')}
            value={edit.current.location ?? ''}
            onChange={(e) => edit.updateField('location', e.target.value)}
            disabled={edit.isSaving}
          />

          {edit.hasError && (
            <Typography variant="body2" sx={{ color: 'error.main' }}>
              {translate('saveError')}
            </Typography>
          )}
        </Stack>
      </DialogContent>

      <DialogActions sx={{ px: 3, pb: 3 }}>
        <Button
          onClick={onClose}
          disabled={edit.isSaving}
          sx={{ borderRadius: 999, px: 3, textTransform: 'none' }}
        >
          {translate('cancel')}
        </Button>

        <Button
          onClick={edit.save}
          disabled={edit.isSaving}
          variant="contained"
          sx={{ borderRadius: 999, px: 3, textTransform: 'none' }}
        >
          {translate('save')}
        </Button>
      </DialogActions>
    </Dialog>
  )
}
