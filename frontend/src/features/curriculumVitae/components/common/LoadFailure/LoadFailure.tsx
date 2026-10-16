import { CurriculumVitaeError } from '@features/curriculumVitae/types'
import { Button, Stack, Typography } from '@mui/material'
import { ROUTES } from '@shared/constants'
import { useTranslate } from '@shared/hooks'
import { Link as RouterLink } from 'react-router'

type Props = {
  error: CurriculumVitaeError
  subject?: 'resume' | 'profile'
}

export const LoadFailure = ({ error, subject = 'resume' }: Props) => {
  const translate = useTranslate('CurriculumVitae.Error')
  const translateProfile = useTranslate('Lk.Error')
  const translateMessage = subject === 'profile' ? translateProfile : translate

  const isNotFound = error.kind === 'notFound'

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
        {translateMessage(isNotFound ? 'notFoundTitle' : 'failedTitle')}
      </Typography>

      <Typography sx={{ color: 'text.secondary' }}>
        {translateMessage(isNotFound ? 'notFoundText' : 'failedText')}
      </Typography>

      {diagnostics && (
        <Typography variant="caption" sx={{ color: 'gray.light' }}>
          {diagnostics}
        </Typography>
      )}

      <Button
        component={RouterLink}
        to={ROUTES.home}
        variant="contained"
        sx={{ borderRadius: 999, px: 3, textTransform: 'none' }}
      >
        {translate('home')}
      </Button>
    </Stack>
  )
}
