import { DownloadOutlined } from '@mui/icons-material'
import { Button, Stack } from '@mui/material'
import { useTranslate } from '@shared/hooks'

type Props = {
  resumeLink?: string
}

/** Главная кнопка hero: скачать файл резюме. Без ссылки на файл не показывается. */
export const Cta = ({ resumeLink }: Props) => {
  const translate = useTranslate('CurriculumVitae.Controls')

  if (!resumeLink) return null

  return (
    <Stack alignItems="center" data-print-hidden sx={{ mt: { xs: 4, md: 6 } }}>
      <Button
        className="hero-action"
        component="a"
        href={resumeLink}
        download
        variant="contained"
        size="large"
        startIcon={<DownloadOutlined />}
        sx={{ borderRadius: 999, px: 4, py: 1.5, textTransform: 'none', fontSize: '1rem' }}
      >
        {translate('download')}
      </Button>
    </Stack>
  )
}
