import { DownloadOutlined } from '@mui/icons-material'
import { IconButton } from '@mui/material'
import { useTranslate } from '@shared/hooks'

type Props = {
  resumeLink?: string
}

export const ResumeButton = ({ resumeLink }: Props) => {
  const translate = useTranslate('CurriculumVitae.Controls')

  if (!resumeLink) return null

  return (
    <IconButton
      component="a"
      href={resumeLink}
      download
      aria-label={translate('download')}
      title={translate('download')}
      sx={{ width: 40, height: 40 }}
    >
      <DownloadOutlined />
    </IconButton>
  )
}
