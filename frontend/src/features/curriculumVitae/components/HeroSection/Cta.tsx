import { useAnalytics } from '@features/curriculumVitae/hooks'
import { AnalyticsEventType } from '@features/curriculumVitae/types'
import { DownloadOutlined } from '@mui/icons-material'
import { Button, Stack } from '@mui/material'
import { useTranslate } from '@shared/hooks'

type Props = {
  resumeLink?: string
}

export const Cta = ({ resumeLink }: Props) => {
  const translate = useTranslate('CurriculumVitae.Controls')
  const { track } = useAnalytics()

  if (!resumeLink) return null

  return (
    <Stack alignItems="center" data-print-hidden sx={{ mt: { xs: 4, md: 6 } }}>
      <Button
        className="hero-action"
        component="a"
        href={resumeLink}
        download
        data-analytics-id="download"
        onClick={() => track({ type: AnalyticsEventType.Download })}
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
