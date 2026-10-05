import { useAnalytics } from '@features/curriculumVitae/hooks'
import { AnalyticsEventType } from '@features/curriculumVitae/types'
import { DownloadOutlined } from '@mui/icons-material'
import { IconButton } from '@mui/material'
import { useTranslate } from '@shared/hooks'

type Props = {
  resumeLink?: string
}

export const ResumeButton = ({ resumeLink }: Props) => {
  const translate = useTranslate('CurriculumVitae.Controls')
  const { track } = useAnalytics()

  if (!resumeLink) return null

  return (
    <IconButton
      component="a"
      href={resumeLink}
      download
      data-analytics-id="download"
      onClick={() => track({ type: AnalyticsEventType.Download })}
      aria-label={translate('download')}
      title={translate('download')}
      sx={{ width: 40, height: 40 }}
    >
      <DownloadOutlined />
    </IconButton>
  )
}
