import { CurriculumVitaeError } from '@features/curriculumVitae/types'
import { Button, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { Link as RouterLink } from 'react-router'

type Props = {
  error: CurriculumVitaeError
}

export const LoadFailure = ({ error }: Props) => {
  const translate = useTranslate('CurriculumVitae.Error')

  const isNotFound = error.kind === 'notFound'

  // Показываем, только если ошибка реально пришла от бэкенда: пригодится, чтобы сообщить в
  // поддержку — по коду и requestId проблему найдут в логах бэкенда
  const diagnostics = [
    error.status && translate('errorCode', { status: error.status }),
    error.requestId && translate('requestId', { requestId: error.requestId }),
  ]
    .filter(Boolean)
    .join(' · ')

  return (
    <Stack
      spacing={2}
      alignItems="center"
      justifyContent="center"
      textAlign="center"
      sx={{ flex: 1, minHeight: '70vh', px: 2 }}
    >
      <Typography variant="h3">
        {translate(isNotFound ? 'notFoundTitle' : 'failedTitle')}
      </Typography>

      <Typography sx={{ color: 'text.secondary' }}>
        {translate(isNotFound ? 'notFoundText' : 'failedText')}
      </Typography>

      {diagnostics && (
        <Typography variant="caption" sx={{ color: 'gray.light' }}>
          {diagnostics}
        </Typography>
      )}

      <Button
        component={RouterLink}
        to="/"
        variant="contained"
        sx={{ borderRadius: 999, px: 3, textTransform: 'none' }}
      >
        {translate('home')}
      </Button>
    </Stack>
  )
}
