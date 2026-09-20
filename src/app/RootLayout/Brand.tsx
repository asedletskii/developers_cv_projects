import { Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'

type Props = {
  appName: string
}

export const Brand = ({ appName }: Props) => {
  const translate = useTranslate('Footer')

  return (
    <Stack spacing={1} maxWidth={'250px'}>
      <Typography variant="h3" color="gray.dark">
        {appName}
      </Typography>

      <Typography variant="caption" sx={{ whiteSpace: 'pre-line' }}>
        {translate('info')}
      </Typography>
    </Stack>
  )
}
