import { CurriculumVitaeError } from '@features/curriculumVitae/types'
import { Button, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { Link as RouterLink } from 'react-router'

type Props = {
  error: CurriculumVitaeError
}

export const LoadFailure = ({ error }: Props) => {
  const translate = useTranslate('CurriculumVitae.Error')

  const isNotFound = error === 'notFound'

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
