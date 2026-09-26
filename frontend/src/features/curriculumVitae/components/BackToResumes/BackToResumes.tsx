import { ArrowBackOutlined, EditRounded } from '@mui/icons-material'
import { Button, Stack } from '@mui/material'
import { ROUTES } from '@shared/constants'
import { useTranslate } from '@shared/hooks'
import { generatePath, Link as RouterLink } from 'react-router'
import { PILL_SX } from './constants'

type Props = {
  userId: string
  cvId?: string
}

export const BackToResumes = ({ userId, cvId }: Props) => {
  const translate = useTranslate('Lk.Resume')

  return (
    <Stack
      data-print-hidden
      direction="row"
      spacing={1}
      sx={{ position: 'fixed', top: 12, left: 12, zIndex: 1201 }}
    >
      <Button
        component={RouterLink}
        to={generatePath(ROUTES.lkProfile, { userId })}
        startIcon={<ArrowBackOutlined />}
        sx={PILL_SX}
      >
        {translate('back')}
      </Button>

      {cvId && (
        <Button
          component={RouterLink}
          to={generatePath(ROUTES.lkResume, { userId, cvId })}
          startIcon={<EditRounded />}
          sx={PILL_SX}
        >
          {translate('edit')}
        </Button>
      )}
    </Stack>
  )
}
