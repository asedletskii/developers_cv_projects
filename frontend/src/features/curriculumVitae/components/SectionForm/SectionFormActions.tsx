import { glassCard } from '@features/curriculumVitae/styles'
import { Button, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { SaveStatus } from './hooks'

type Props = {
  status: SaveStatus
  isDirty: boolean
  onReset: () => void
}

/** Плавающая панель формы: статус сохранения и кнопки «Сбросить»/«Сохранить». Прилипает к низу
 *  окна, пока форма длиннее экрана. */
export const SectionFormActions = ({ status, isDirty, onReset }: Props) => {
  const translate = useTranslate('Lk.Edit')

  return (
    <Stack
      direction="row"
      alignItems="center"
      spacing={1.5}
      sx={{
        position: 'sticky',
        bottom: 16,
        zIndex: 10,
        px: 2,
        py: 1,
        ...glassCard,
        borderRadius: 999,
      }}
    >
      <Typography
        variant="body2"
        role="status"
        sx={{ flex: 1, color: status === 'error' ? 'error.main' : 'text.secondary' }}
      >
        {status === 'error' && translate('saveError')}
        {status === 'saved' && translate('saved')}
        {status === 'idle' && isDirty && translate('unsaved')}
      </Typography>

      <Button
        onClick={onReset}
        disabled={!isDirty || status === 'saving'}
        sx={{ borderRadius: 999, px: 3, textTransform: 'none' }}
      >
        {translate('reset')}
      </Button>

      <Button
        type="submit"
        variant="contained"
        disabled={status === 'saving'}
        sx={{ borderRadius: 999, px: 3, textTransform: 'none' }}
      >
        {status === 'saving' ? translate('saving') : translate('save')}
      </Button>
    </Stack>
  )
}
