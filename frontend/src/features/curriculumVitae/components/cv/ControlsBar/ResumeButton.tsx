import { useDownloadResume } from '@features/curriculumVitae/hooks'
import { DownloadOutlined } from '@mui/icons-material'
import { CircularProgress, IconButton } from '@mui/material'
import { useTranslate } from '@shared/hooks'

/** Резюме собирается в HTML на странице и скачивается — см. `useDownloadResume` */
export const ResumeButton = () => {
  const translate = useTranslate('CurriculumVitae.Controls')
  const { download, isDownloading } = useDownloadResume()

  return (
    <IconButton
      data-analytics-id="download-html"
      onClick={download}
      disabled={isDownloading}
      aria-label={translate('download')}
      title={translate('download')}
      sx={{ width: 40, height: 40 }}
    >
      {isDownloading ? <CircularProgress size={20} color="inherit" /> : <DownloadOutlined />}
    </IconButton>
  )
}
