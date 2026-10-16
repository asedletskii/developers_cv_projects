import { glassCard } from '@features/curriculumVitae/styles'
import { RESUME_LANGUAGES, ResumeLanguage } from '@features/curriculumVitae/types'
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  MenuItem,
  Stack,
  TextField,
} from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useCreateResume } from './hooks'

type Props = {
  onClose: () => void
}

export const CreateResumeDialog = ({ onClose }: Props) => {
  const translate = useTranslate('Lk.Resumes')
  const translateAbout = useTranslate('CurriculumVitae.About')
  const { title, setTitle, language, setLanguage, isSaving, hasError, canCreate, create } =
    useCreateResume()

  return (
    <Dialog
      open
      onClose={onClose}
      fullWidth
      maxWidth="xs"
      slotProps={{ paper: { sx: { ...glassCard, backgroundImage: 'none' } } }}
    >
      <DialogTitle>{translate('addResume')}</DialogTitle>

      <DialogContent>
        <Stack spacing={2} sx={{ pt: 1 }}>
          <TextField
            autoFocus
            fullWidth
            label={translate('titleLabel')}
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            disabled={isSaving}
          />

          <TextField
            select
            fullWidth
            label={translate('languageLabel')}
            value={language}
            onChange={(e) => setLanguage(e.target.value as ResumeLanguage)}
            error={hasError}
            helperText={hasError ? translate('createError') : undefined}
            disabled={isSaving}
          >
            {RESUME_LANGUAGES.map((code) => (
              <MenuItem key={code} value={code}>
                {translateAbout(`Languages.${code}`)}
              </MenuItem>
            ))}
          </TextField>
        </Stack>
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
          onClick={create}
          disabled={!canCreate}
          variant="contained"
          sx={{ borderRadius: 999, px: 3, textTransform: 'none' }}
        >
          {translate('create')}
        </Button>
      </DialogActions>
    </Dialog>
  )
}
