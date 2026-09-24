import { useState } from 'react'
import { glassCard } from '@features/curriculumVitae/styles'
import { VisibilityOutlined } from '@mui/icons-material'
import { Button, Stack, Typography } from '@mui/material'
import { useTheme } from '@mui/material/styles'
import { useTranslate } from '@shared/hooks'
import { ColorField } from './ColorField'
import { ThemePreviewDialog } from './ThemePreviewDialog'

/** Цветовая гамма резюме — необязательная, без значений используются цвета темы сайта по умолчанию. */
export const ThemeEditor = () => {
  const translate = useTranslate('Lk.Edit.Theme')
  const [isPreviewOpen, setIsPreviewOpen] = useState(false)
  const { palette } = useTheme()

  return (
    <Stack spacing={2} sx={{ p: { xs: 2, sm: 3 }, ...glassCard }}>
      <Stack direction="row" alignItems="center" justifyContent="space-between">
        <Typography variant="h6">{translate('title')}</Typography>

        <Button
          size="small"
          startIcon={<VisibilityOutlined />}
          onClick={() => setIsPreviewOpen(true)}
          sx={{ textTransform: 'none' }}
        >
          {translate('preview')}
        </Button>
      </Stack>

      <Stack direction={{ xs: 'column', sm: 'row' }} spacing={3}>
        <ColorField
          name="theme.primary"
          label={translate('primary')}
          defaultColor={palette.primary.main}
        />
        <ColorField
          name="theme.secondary"
          label={translate('secondary')}
          defaultColor={palette.secondary.main}
        />
      </Stack>

      {isPreviewOpen && <ThemePreviewDialog onClose={() => setIsPreviewOpen(false)} />}
    </Stack>
  )
}
