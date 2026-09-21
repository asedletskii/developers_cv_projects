import { SchoolOutlined } from '@mui/icons-material'
import { Stack, Typography } from '@mui/material'
import { IconTile } from '@features/curriculumVitae/components/IconTile'
import { glassCard } from '@features/curriculumVitae/styles'
import { Education } from '@features/curriculumVitae/types'
import { useTranslate } from '@shared/hooks'

type Props = {
  education: Education
}

export const Entry = ({ education }: Props) => {
  const translate = useTranslate('CurriculumVitae.Education')

  return (
    <Stack className="reveal-item" direction="row" spacing={2.5} sx={{ p: 3, ...glassCard }}>
      <IconTile>
        <SchoolOutlined />
      </IconTile>

      <Stack spacing={0.5} sx={{ minWidth: 0 }}>
        <Typography variant="h6" sx={{ fontWeight: 700 }}>
          {education.university}
        </Typography>

        <Typography>
          {[
            translate(`Levels.level${education.grade}`),
            education.programName,
            education.profileName,
          ]
            .filter(Boolean)
            .join(' · ')}
        </Typography>

        <Typography variant="caption" sx={{ color: 'gray.light' }}>
          {[
            education.faculty,
            education.location,
            education.endDate
              ? translate('graduated', { year: education.endDate })
              : education.course && translate('course', { course: education.course }),
          ]
            .filter(Boolean)
            .join(' · ')}
        </Typography>
      </Stack>
    </Stack>
  )
}
