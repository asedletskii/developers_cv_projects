import { DownloadOutlined } from '@mui/icons-material'
import { Button, CircularProgress, Stack } from '@mui/material'
import { useTranslate } from '@shared/hooks'

type Props = {
  /** Скачать резюме. Без обработчика кнопка только выглядит так же — для предпросмотра в редакторе */
  onClick?: () => void
  isLoading?: boolean
}

export const Cta = ({ onClick, isLoading = false }: Props) => {
  const translate = useTranslate('CurriculumVitae.Controls')

  return (
    <Stack alignItems="center" data-print-hidden sx={{ mt: { xs: 4, md: 6 } }}>
      <Button
        className="hero-action"
        data-analytics-id="download-html"
        onClick={onClick}
        disabled={isLoading}
        variant="contained"
        size="large"
        startIcon={
          isLoading ? <CircularProgress size={18} color="inherit" /> : <DownloadOutlined />
        }
        sx={{ borderRadius: 999, px: 4, py: 1.5, textTransform: 'none', fontSize: '1rem' }}
      >
        {translate('download')}
      </Button>
    </Stack>
  )
}
