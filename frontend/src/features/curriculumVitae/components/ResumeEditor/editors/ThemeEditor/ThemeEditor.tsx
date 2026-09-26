import { ColorField } from '@features/curriculumVitae/components/form/ColorField'
import { ThemePreviewDialog } from '@features/curriculumVitae/components/ThemePreviewDialog'
import { glassCard } from '@features/curriculumVitae/styles'
import { VisibilityOutlined } from '@mui/icons-material'
import { Button, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useThemeEditor } from './hooks'

export const ThemeEditor = () => {
  const translate = useTranslate('Lk.Edit.Theme')
  const { defaultPrimary, defaultSecondary, isPreviewOpen, openPreview, closePreview } =
    useThemeEditor()

  return (
    <Stack spacing={2} sx={{ p: { xs: 2, sm: 3 }, ...glassCard }}>
      <Stack direction="row" alignItems="center" justifyContent="space-between">
        <Typography variant="h6">{translate('title')}</Typography>

        <Button
          size="small"
          startIcon={<VisibilityOutlined />}
          onClick={openPreview}
          sx={{ textTransform: 'none' }}
        >
          {translate('preview')}
        </Button>
      </Stack>

      <Stack direction={{ xs: 'column', sm: 'row' }} spacing={3}>
        <ColorField
          name="theme.primary"
          label={translate('primary')}
          defaultColor={defaultPrimary}
        />
        <ColorField
          name="theme.secondary"
          label={translate('secondary')}
          defaultColor={defaultSecondary}
        />
      </Stack>

      {isPreviewOpen && <ThemePreviewDialog onClose={closePreview} />}
    </Stack>
  )
}
