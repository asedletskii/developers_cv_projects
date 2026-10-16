import { ExternalLink } from '@features/curriculumVitae/components/cv/ExternalLink'
import { IconTile } from '@features/curriculumVitae/components/cv/IconTile'
import { glassCard } from '@features/curriculumVitae/styles'
import { Course as CourseData } from '@features/curriculumVitae/types'
import { WorkspacePremiumOutlined } from '@mui/icons-material'
import { Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { formatMonth } from '@shared/utils'

type Props = {
  course: CourseData
}

export const Course = ({ course }: Props) => {
  const translate = useTranslate('CurriculumVitae.Courses')

  return (
    <Stack className="reveal-item" direction="row" spacing={2} sx={{ p: 3, ...glassCard }}>
      <IconTile color="secondary.main">
        <WorkspacePremiumOutlined />
      </IconTile>

      <Stack spacing={0.5} sx={{ minWidth: 0 }}>
        <Typography variant="h6" sx={{ fontWeight: 700 }}>
          {course.name}
        </Typography>

        <Typography variant="caption" sx={{ color: 'gray.light' }}>
          {[
            course.organization,
            formatMonth(course.endDate),
            course.hours && translate('hours', { hours: course.hours }),
          ]
            .filter(Boolean)
            .join(' · ')}
        </Typography>

        {course.certificateLink && (
          <ExternalLink href={course.certificateLink}>{translate('certificate')}</ExternalLink>
        )}
      </Stack>
    </Stack>
  )
}
